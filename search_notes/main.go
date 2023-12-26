package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
	"golang.org/x/text/unicode/norm"
	"google.golang.org/protobuf/proto"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sballin/alfred-search-notes-app/alfred"
	notestore "github.com/sballin/alfred-search-notes-app/proto"
)

var matcher = search.New(language.Und, search.Loose)

const (
	DbPath           = "~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"
	TitleKey         = "title"    // titles of rows in Alfred
	SubtitleKey      = "subtitle" // subtitles of rows in Alfred
	ArgKey           = "arg"      // comma-separated lists of identifiers that wind up in Alfred "arg" fields
	BodyKey          = "noteBodyZipped"
	TableTextKey     = "tableText"
	NotesSQLTemplate = `
SELECT 
    noteTitle AS title,
    folderTitle AS subtitle,
    identifier || ',' || -- note ID used in notes:// and applenotes:// URI schemes
        'x-coredata://' || z_uuid || '/ICNote/p' || xcoreDataID || ',' || -- applescript ID of note
        IFNULL('x-coredata://' || z_uuid || '/ICAccount/p' || accountID, 'null') || ',' || -- applescript ID of account
        %s -- applescript ID of folder that note is in or "null"
        AS arg,
    noteBodyZipped,
    tableText,
    CAST(xcoreDataID AS TEXT)
FROM (
    SELECT
        c.ztitle1 AS noteTitle,
        c.zfolder AS noteFolderID,
        c.zmodificationdate1 AS modDate,
        c.z_pk AS xcoredataID,
        c.zaccount3 AS accountID,
        c.zidentifier AS identifier,
        n.zdata AS noteBodyZipped
    FROM 
        ziccloudsyncingobject AS c
        INNER JOIN zicnotedata AS n ON c.znotedata = n.z_pk -- note id (int) distinct from xcoredataID
    WHERE 
        noteTitle IS NOT NULL AND 
        modDate IS NOT NULL AND
        xcoredataID IS NOT NULL AND
        noteBodyZipped IS NOT NULL AND
        c.zmarkedfordeletion != 1
) AS notes
INNER JOIN (
    SELECT
        z_pk AS folderID,
        ztitle2 AS folderTitle,
        zfoldertype AS isRecentlyDeletedFolder
    FROM ziccloudsyncingobject
    WHERE 
        folderTitle IS NOT NULL AND 
        isRecentlyDeletedFolder != 1 AND
        zmarkedfordeletion != 1 
) AS folders ON noteFolderID = folderID
LEFT JOIN (
    SELECT 
        GROUP_CONCAT(zsummary, '') AS tableText,
        znote
    FROM ziccloudsyncingobject
    WHERE ztypeuti = 'com.apple.notes.table'
    GROUP BY znote
) AS tables ON znote = xcoreDataID
LEFT JOIN (
    SELECT z_uuid FROM z_metadata
)
ORDER BY %s
`

	OCRsSQL = `
SELECT 
    CAST(znote AS TEXT),
    IFNULL(GROUP_CONCAT(zhandwritingsummary, ''), '') || IFNULL(GROUP_CONCAT(zocrsummary, ''), '')
FROM ziccloudsyncingobject
WHERE (zocrsummary IS NOT NULL OR
       zhandwritingsummary IS NOT NULL) AND
      zmarkedfordeletion != 1 
GROUP BY znote
`

	HashtagsSQL = `
SELECT 
    CAST(znote1 AS TEXT),
    GROUP_CONCAT(zalttext, ' ')
FROM ziccloudsyncingobject
WHERE ztypeuti1 = 'com.apple.notes.inlinetextattachment.hashtag'
GROUP BY znote1
`

	FoldersSQLTemplate = `
SELECT 
    ztitle2 AS title,
    '' AS subtitle,
    zidentifier || ',x-coredata://' || z_uuid || '/ICFolder/p' || z_pk || ',' || IFNULL('x-coredata://' || z_uuid || '/ICAccount/p' || zaccount4, 'null') AS arg
FROM ziccloudsyncingobject
LEFT JOIN (
    SELECT z_uuid FROM z_metadata
)
WHERE 
    title IS NOT NULL AND 
    zmarkedfordeletion != 1 AND
    zneedsinitialfetchfromcloud != 1 -- some phantom folders can show up in the results
ORDER BY title ASC
`
)

type LiteDB struct {
	db *sql.DB
}

type UserQuery struct {
	Tokens     []string
	WordString string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Escape(s string) string {
	return strings.Replace(s, "'", "''", -1)
}

func Expanduser(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path[:2] == "~/" {
		path = filepath.Join(dir, path[2:])
	}
	return path
}

func NewLiteDB(path string) (LiteDB, error) {
	db, err := sql.Open("sqlite3", path)
	litedb := LiteDB{db}
	return litedb, err
}

func NewNotesDB() (LiteDB, error) {
	path := Expanduser(DbPath)
	litedb, err := NewLiteDB("file:" + path + "?mode=ro&_query_only=true")
	return litedb, err
}

func SafeUnicode(r rune) rune {
	// Keep graphic characters and newlines
	if unicode.IsGraphic(r) || r == '\n' {
		return r
	} else {
		return -1
	}
}

func GetNoteBody(noteBytes []byte) string {
	body := ""
	note := &notestore.NoteStoreProto{}
	err := proto.Unmarshal(noteBytes, note)
	if err != nil {
		return ""
	}
	if note.Document.Note.NoteText != nil {
		body += *note.Document.Note.NoteText
		for _, a := range note.Document.Note.AttributeRun {
			if a.Link != nil {
				body += "\n" + *a.Link
			}
		}
		// Remove title from body
		bodyStart := strings.Index(body, "\n")
		if bodyStart > 0 {
			body = body[bodyStart:]
		}
		// Remove object substitution character
		body = strings.ReplaceAll(body, string([]byte{239, 191, 188}), "")
		// Remove any weird characters that might be left over
		body = strings.Map(SafeUnicode, body)
		body = norm.NFC.String(strings.ToValidUTF8(body, ""))
	}
	return body
}

func SubtitleMatchSummary(body string, search string) string {
	matchSummary := " | …"
	i := 0
	j := 0
	k := 0
	for i >= 0 && j >= 0 && len(matchSummary) < 400 {
		j, k = matcher.IndexString(body[i:], search)
		if j >= 0 {
			// Include context around match up to rb or next newline
			rb := min(len(body), i+k+25)
			nextNewline := strings.Index(body[i+j:rb], "\n")
			if nextNewline > 0 {
				rb = i + j + nextNewline
			}
			match := strings.ToValidUTF8(strings.Trim(body[i+j:rb], " "), "")
			matchSummary += match + "…"
			i = rb
		}
	}
	return matchSummary
}

func (lite LiteDB) GetSpecialColumn(query string) (map[string]string, error) {
	results := map[string]string{}
	rows, err := lite.db.Query(query)
	if err != nil {
		return results, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return results, err
	}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}
		valNoteID := columnPointers[0].(*interface{})
		noteID, ok := (*valNoteID).(string)
		if !ok {
			continue
		}
		valSpecialColumn := columnPointers[1].(*interface{})
		specialColumn, ok := (*valSpecialColumn).(string)
		if !ok {
			continue
		}
		results[noteID] = specialColumn
	}
	return results, nil
}

func (lite LiteDB) GetResults(search string, scope string) ([]map[string]string, error) {
	// Format SQL query
	sqlQuery := fmt.Sprintf(NotesSQLTemplate, GetEnclosingFolderPreference(), GetOrderPreference())
	if scope == "folder" {
		sqlQuery = FoldersSQLTemplate
	}

	// Get OCR text
	OCRs := map[string]string{}
	if scope == "body" {
		OCRs, _ = lite.GetSpecialColumn(OCRsSQL)
	}

	// Get hashtags
	hashtags := map[string]string{}
	if scope == "body" || scope == "hashtag" {
		hashtags, _ = lite.GetSpecialColumn(HashtagsSQL)
	}

	// Run query to get all results
	results := []map[string]string{}
	rows, err := lite.db.Query(sqlQuery)
	if err != nil {
		return results, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return results, err
	}

	searchFolders := false
	if os.Getenv("searchFolders") != "0" {
		searchFolders = true
	}

	gzipHeader := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 19}
	bytesReader := bytes.NewReader(gzipHeader)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil && scope == "body" {
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
			continue
		}

		hashtagText := ""
		scopeText := ""
		matchSummary := ""
		if len(search) > 0 {
			ocrText := ""
			if scope == "body" || scope == "hashtag" {
				// Get note ID and OCR text
				valNoteID := columnPointers[5].(*interface{})
				noteID, ok := (*valNoteID).(string)
				if ok {
					ocrText = OCRs[noteID]
					hashtagText = hashtags[noteID]
				}
			}

			// Add note/folder title to search scope
			valTitle := columnPointers[0].(*interface{})
			title, ok := (*valTitle).(string)
			if !ok {
				continue
			}
			scopeText = title
			if searchFolders {
				// Add folder of note object to search scope (this field is empty for folder objects)
				valFolder := columnPointers[1].(*interface{})
				folder, ok := (*valFolder).(string)
				if !ok {
					folder = ""
				}
				scopeText += " " + folder
			}
			if scope == "hashtag" {
				// Return notes matching every hash tag provided
				searchTags := strings.Split(search, " ")
				containsSearchTag := true
				for _, searchTag := range searchTags {
					firstMatch, _ := matcher.IndexString(hashtagText, "#"+searchTag)
					if firstMatch == -1 {
						containsSearchTag = false
					}
				}
				if !containsSearchTag {
					continue
				}
			} else {
				if scope == "body" {
					// Decompress note body data
					valBody := columnPointers[3].(*interface{})
					noteDataZippedBytes, ok := (*valBody).([]byte)
					if ok {
						bytesReader.Reset(noteDataZippedBytes)
						errReset := gzipReader.Reset(bytesReader)
						if errReset == nil {
							noteBytes, errRead := io.ReadAll(gzipReader)
							if errRead == nil {
								// Get plaintext of any tables in this note
								valTableText := columnPointers[4].(*interface{})
								tableText, ok := (*valTableText).(string)
								if !ok {
									tableText = ""
								}
								// Extract protobuf-format data from unzipped note and add other text
								body := hashtagText + " " + GetNoteBody(noteBytes) + " " + tableText + " " + ocrText
								// Add body text to search scope
								scopeText += " " + body
								// Prepare result summary for subtitle string
								firstMatch, _ := matcher.IndexString(body, search)
								if firstMatch >= 0 {
									matchSummary = SubtitleMatchSummary(body, search)
								}
							}
						}
					}
				}
				firstMatch, _ := matcher.IndexString(scopeText, search)
				if firstMatch == -1 {
					continue
				}
			}
		}

		// If we get here, the note/folder contains a match. Add it to the Alfred results.
		for i, colName := range cols {
			// Don't add note body data to future alfred row
			if colName == BodyKey {
				continue
			}

			val := columnPointers[i].(*interface{})
			s, ok := (*val).(string)
			if ok {
				m[colName] = s
			} else {
				m[colName] = ""
			}
		}

		// Add additional text to subtitle string
		if hashtagText != "" {
			hashtagText = " " + hashtagText
		}
		m[SubtitleKey] += hashtagText
		m[SubtitleKey] += matchSummary

		results = append(results, m)
	}
	return results, err
}

func RowToItem(row map[string]string, userQuery UserQuery) alfred.Item {
	return alfred.Item{
		Title:        row[TitleKey],
		Subtitle:     row[SubtitleKey],
		Arg:          row[ArgKey] + "," + Escape(userQuery.WordString),
		QuicklookURL: " ",
	}
}

func CreateNoteItem(userQuery UserQuery) (*alfred.Item, error) {
	title := userQuery.WordString
	item := alfred.Item{
		Title:    title,
		Arg:      title,
		Subtitle: "Create new note",
	}
	return &item, nil
}

func ParseUserQuery(arg string) UserQuery {
	userQuery := UserQuery{}
	userQuery.Tokens = strings.Split(norm.NFC.String(arg), " ")
	words := []string{}
	for _, e := range userQuery.Tokens {
		if e != "" {
			words = append(words, e)
		}
	}
	userQuery.WordString = strings.Join(words, " ")
	return userQuery
}

func GetEnclosingFolderPreference() string {
	if os.Getenv("showEnclosingFolder") != "0" {
		return "'x-coredata://' || z_uuid || '/ICFolder/p' || noteFolderID"
	} else {
		return "'null'"
	}
}

func GetOrderPreference() string {
	if os.Getenv("sortByDate") != "0" {
		return "modDate DESC"
	} else {
		return "lower(noteTitle) ASC"
	}
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) >= 3 {
		litedb, err := NewNotesDB()
		PanicOnErr(err)

		scope := os.Args[1]
		userQuery := ParseUserQuery(os.Args[2])

		searchRows, err := litedb.GetResults(userQuery.WordString, scope)
		PanicOnErr(err)

		if os.Getenv("fallbackCreateNew") == "1" && (scope == "title" || scope == "body") && len(searchRows) == 0 {
			createItem, err := CreateNoteItem(userQuery)
			PanicOnErr(err)
			alfred.Add(*createItem)
		}

		if os.Getenv("fallbackSearchBody") == "1" && scope == "title" && len(searchRows) == 0 {
			searchRows, err = litedb.GetResults(userQuery.WordString, "body")
			PanicOnErr(err)
		}

		if len(searchRows) > 0 {
			for _, row := range searchRows {
				alfred.Add(RowToItem(row, userQuery))
			}
		}

		alfred.Run()
	}
}
