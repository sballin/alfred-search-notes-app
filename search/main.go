package main

import (
	"os"
	"os/user"
	"strings"
	"fmt"
	"path/filepath"
	"database/sql"
	"golang.org/x/text/unicode/norm"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sballin/alfred-search-notes-app/alfred"
)

const (
	DbPath = "~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"

	TitleKey  = "noteTitle"
	FolderKey = "folderName"
	NoteIDKey = "'x-coredata://' || z_uuid || '/ICNote/p' || xcoreDataID"

	NOTES_BY_QUERY = `
SELECT 
    'x-coredata://' || z_uuid || '/ICNote/p' || xcoreDataID,
    noteTitle,
    folderName
FROM
    (SELECT c.ztitle1 AS noteTitle,
            c.zfolder AS noteFolderID,
            c.zmodificationdate1 AS modDate,
            c.z_pk AS xcoredataID,
            n.zdata AS noteBodyZipped
     FROM 
         ziccloudsyncingobject AS c
         INNER JOIN zicnotedata AS n ON c.znotedata = n.z_pk -- note id (int) distinct from xcoredataID
     WHERE noteTitle IS NOT NULL AND 
           noteFolderID > 1 AND -- 1 is the Recently Deleted folder
           modDate IS NOT NULL AND
           xcoredataID IS NOT NULL AND
           noteBodyZipped IS NOT NULL AND
           c.zmarkedfordeletion != 1) AS notes
    LEFT JOIN 
        (SELECT z_pk AS folderID,
                ztitle2 AS folderName
         FROM ziccloudsyncingobject
         WHERE ztitle2 IS NOT NULL AND 
         zmarkedfordeletion != 1) AS folders ON noteFolderID = folderID
    LEFT JOIN
        (SELECT z_uuid FROM z_metadata)
WHERE (lower(noteTitle) LIKE lower('%%%s%%') OR %s)
ORDER BY %s
LIMIT 100
`)

func Expanduser(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path[:2] == "~/" {
		path = filepath.Join(dir, path[2:])
	}
	return path
}

type LiteDB struct {
	db *sql.DB
}

func NewLiteDB(path string) (LiteDB, error) {
	db, err := sql.Open("sqlite3", path)
	litedb := LiteDB{db}
	return litedb, err
}

func NewNotesDB() (LiteDB, error) {
	path := Expanduser(DbPath)
	litedb, err := NewLiteDB(path)
	return litedb, err
}

func (lite LiteDB) Query(q string) ([]map[string]string, error) {
	results := []map[string]string{}
	rows, err := lite.db.Query(q)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return results, err
	}

	for rows.Next() {
		m := map[string]string{}
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return results, err
		}
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			s, ok := (*val).(string)
			if ok {
				m[colName] = s
			} else {
				m[colName] = ""
			}
		}
		results = append(results, m)
	}
	return results, err
}

func RowToItem(row map[string]string, query Query) alfred.Item {
	return alfred.Item{
		Title:    row[TitleKey],
		Subtitle: row[FolderKey],
		Arg:      row[NoteIDKey] + "?" + escape(query.WordString),
		QuicklookURL: nil,
	}
}

type Query struct {
	Tokens     []string
	Tags       []string
	LastToken  string
	WordString string
}

func (query Query) String() string {
	return strings.Join(query.Tokens, " ")
}

func ParseQuery(arg string) Query {
	query := Query{}
	query.Tokens = strings.Split(norm.NFC.String(arg), " ")
	query.Tags = []string{}
	words := []string{}
	for _, e := range query.Tokens {
		switch {
		case e == "":
		case strings.HasPrefix(e, "#"):
			query.Tags = append(query.Tags, e)
		default:
			words = append(words, e)
		}
	}
	query.LastToken = query.Tokens[len(query.Tokens)-1]
	query.WordString = strings.Join(words, " ")
	return query
}

func escape(s string) string {
	return strings.Replace(s, "'", "''", -1)
}

func GetSearchRows(litedb LiteDB, query Query) ([]map[string]string, error) {
	escapedWordString := escape(query.WordString)
	var sortByDate = os.Getenv("sortByDate")
	var orderBy = "modDate DESC"
	if (sortByDate == "0") {
		orderBy = "lower(noteTitle) ASC"
	}
	var searchFolders = os.Getenv("searchFolders")
	var foldersLike = fmt.Sprintf("lower(folderName) LIKE lower('%%%s%%')", escapedWordString)
	if (searchFolders == "0") {
		 foldersLike = "FALSE"
	}
	rows, err := litedb.Query(fmt.Sprintf(NOTES_BY_QUERY, escapedWordString, foldersLike, orderBy))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetCreateItem(query Query) (*alfred.Item, error) {
	title := query.WordString
	item := alfred.Item{
		Title: title,
		Arg:   title,
		Subtitle: "Create new note",
	}
	if len(query.Tags) != 0 {
		item.Subtitle = strings.Join(query.Tags, " ")
	}
	return &item, nil
}

func main() {
	if len(os.Args) > 1 {
		query := ParseQuery(os.Args[2])

		litedb, err := NewNotesDB()
		if err != nil {
			panic(err)
		}

		searchRows, err := GetSearchRows(litedb, query)
		if err != nil {
			panic(err)
		}

		createItem, err := GetCreateItem(query)
		if err != nil {
			panic(err)
		}

		if len(searchRows) > 0 {
			for _, row := range searchRows {
				alfred.Add(RowToItem(row, query))
			}
		} else {
			alfred.Add(*createItem)
		}

		alfred.Run()
	}
}
