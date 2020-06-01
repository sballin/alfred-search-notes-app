package core

import (
	"fmt"
	"strings"
	"os"
	
	"golang.org/x/text/unicode/norm"

	"github.com/drgrib/alfred-bear/db"
	"github.com/drgrib/alfred-bear/alfred"
)

const argSplit = "|"

func RowToItem(row map[string]string, query Query) alfred.Item {
	return alfred.Item{
		Title:    row[db.TitleKey],
		Subtitle: row[db.FolderKey],
		Arg:      row[db.NoteIDKey] + "?" + escape(query.WordString),
		QuicklookURL: nil,
	}
}

func AddNoteRowsToAlfred(rows []map[string]string, query Query) {
	for _, row := range rows {
		item := RowToItem(row, query)
		alfred.Add(item)
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

func GetSearchRows(litedb db.LiteDB, query Query) ([]map[string]string, error) {
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
	rows, err := litedb.Query(fmt.Sprintf(db.NOTES_BY_QUERY, escapedWordString, foldersLike, orderBy))
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
