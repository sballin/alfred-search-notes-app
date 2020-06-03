package main

import (
	"os"
	"os/user"
	"strings"
	"fmt"
	"path/filepath"
	"database/sql"
	"encoding/hex"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"golang.org/x/text/unicode/norm"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sballin/alfred-search-notes-app/alfred"
)

const (
	DbPath = "~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"

	TitleKey  = "title"
	SubtitleKey = "subtitle"
	ArgKey = "URL"
	BodyKey = "noteBodyHex"

	NOTES = `
SELECT 
    noteTitle as title,
    folderTitle as subtitle,
    'x-coredata://' || z_uuid || '/ICNote/p' || xcoreDataID as URL,
    HEX(noteBodyZipped) as noteBodyHex
FROM (
	SELECT
		c.ztitle1 AS noteTitle,
        c.zfolder AS noteFolderID,
        c.zmodificationdate1 AS modDate,
        c.z_pk AS xcoredataID,
        n.zdata AS noteBodyZipped
    FROM 
        ziccloudsyncingobject AS c
        INNER JOIN zicnotedata AS n ON c.znotedata = n.z_pk -- note id (int) distinct from xcoredataID
    WHERE 
        noteTitle IS NOT NULL AND 
        noteFolderID > 1 AND -- 1 is the Recently Deleted folder
        modDate IS NOT NULL AND
        xcoredataID IS NOT NULL AND
        noteBodyZipped IS NOT NULL AND
        c.zmarkedfordeletion != 1
) AS notes
LEFT JOIN (
    SELECT
        z_pk AS folderID,
        ztitle2 AS folderTitle
     FROM ziccloudsyncingobject
     WHERE 
         folderTitle IS NOT NULL AND 
         zmarkedfordeletion != 1 
) AS folders ON noteFolderID = folderID
LEFT JOIN (
	SELECT z_uuid FROM z_metadata
)
%s -- either blank to get all notes or WHERE titles/folders like...
ORDER BY %s
`

	FOLDERS_BY_TITLE = `
SELECT 
	ztitle2 AS title,
	'' AS subtitle,
	'x-coredata://' || z_uuid || '/ICFolder/p' || z_pk as URL
FROM ziccloudsyncingobject
LEFT JOIN (
	SELECT z_uuid FROM z_metadata
)
WHERE 
    z_pk > 1 AND -- 1 is the Recently Deleted folder
    title IS NOT NULL AND 
	zmarkedfordeletion != 1 AND
    lower(title) LIKE lower('%%%s%%')
ORDER BY title ASC
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
		Title:        row[TitleKey],
		Subtitle:     row[SubtitleKey],
		Arg:          row[ArgKey] + "?" + escape(query.WordString),
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

func GetCreateNote(query Query) (*alfred.Item, error) {
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

func GetSearchTitleRows(litedb LiteDB, query Query) ([]map[string]string, error) {
	escapedWordString := escape(query.WordString)
	sortByDate := os.Getenv("sortByDate")
	orderBy := "modDate DESC"
	if (sortByDate == "0") {
		orderBy = "lower(noteTitle) ASC"
	}
	searchFolders := os.Getenv("searchFolders")
	foldersLike := fmt.Sprintf("lower(folderTitle) LIKE lower('%%%s%%')", escapedWordString)
	if (searchFolders == "0") {
		 foldersLike = "FALSE"
	}
	searchString := fmt.Sprintf("WHERE (lower(noteTitle) LIKE lower('%%%s%%') OR %s)", escapedWordString, foldersLike)
	rows, err := litedb.Query(fmt.Sprintf(NOTES, searchString, orderBy))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetSearchBodyRows(litedb LiteDB, query Query) ([]map[string]string, error) {
	escapedWordString := escape(query.WordString)
	sortByDate := os.Getenv("sortByDate")
	orderBy := "modDate DESC"
	if (sortByDate == "0") {
		orderBy = "lower(noteTitle) ASC"
	}
	searchString := ""
	rows, err := litedb.Query(fmt.Sprintf(NOTES, searchString, orderBy))
	if err != nil {
		return nil, err
	}	
	
	rowsMatching := []map[string]string{}
	for _, row := range rows {
		decoded, err := hex.DecodeString(row[BodyKey])
		if err != nil {
			continue
		}
		r, err := gzip.NewReader(bytes.NewReader(decoded))
		if err != nil {
			continue
		}
		body, err := ioutil.ReadAll(r)
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(string(body)), strings.ToLower(escapedWordString)) {
			rowsMatching = append(rowsMatching, row)
		}
		r.Close()
	}
	return rowsMatching, nil
}

func GetSearchFolderRows(litedb LiteDB, query Query) ([]map[string]string, error) {
	escapedWordString := escape(query.WordString)
	rows, err := litedb.Query(fmt.Sprintf(FOLDERS_BY_TITLE, escapedWordString))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func main() {
	if len(os.Args) >= 3 {
		litedb, err := NewNotesDB()
		if err != nil {
			panic(err)
		}
		
		query := ParseQuery(os.Args[2])
		searchRows := []map[string]string{}
			
		if os.Args[1] == "title" {
			searchRows, err = GetSearchTitleRows(litedb, query)
			if err != nil {
				panic(err)
			}
			
			if len(searchRows) == 0 {
				createItem, err := GetCreateNote(query)
				if err != nil {
					panic(err)
				}
				alfred.Add(*createItem)
			}
		} else if os.Args[1] == "body" {
			searchRows, err = GetSearchBodyRows(litedb, query)
			if err != nil {
				panic(err)
			}	
		} else if os.Args[1] == "folder" {
			searchRows, err = GetSearchFolderRows(litedb, query)
			if err != nil {
				panic(err)
			}
		}
		
		if len(searchRows) > 0 {
			for _, row := range searchRows {
				alfred.Add(RowToItem(row, query))
			}
		}

		alfred.Run()
	}
}
