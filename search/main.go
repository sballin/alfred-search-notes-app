package main

import (
    "os"
    "os/user"
    "strings"
    "fmt"
    "path/filepath"
    "database/sql"
    "bytes"
    "compress/gzip"
    "io/ioutil"
    "google.golang.org/protobuf/proto"
    "golang.org/x/text/unicode/norm"
    "unicode"

    _ "github.com/mattn/go-sqlite3"
    "github.com/sballin/alfred-search-notes-app/alfred"
    notestore "github.com/sballin/alfred-search-notes-app/proto"
)

const (
    DbPath = "~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"

    TitleKey  = "title"
    SubtitleKey = "subtitle"
    ArgKey = "url"
    BodyKey = "noteBodyZipped"
    TableTextKey = "tableText"

    NotesSQLTemplate = `
SELECT 
    noteTitle AS title,
    folderTitle AS subtitle,
    'x-coredata://' || z_uuid || '/ICNote/p' || xcoreDataID AS url,
    noteBodyZipped,
    tableText
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
%s -- either blank to get all notes or WHERE titles/folders like...
ORDER BY %s
`

    FoldersSQLTemplate = `
SELECT 
    ztitle2 AS title,
    '' AS subtitle,
    'x-coredata://' || z_uuid || '/ICFolder/p' || z_pk AS url
FROM ziccloudsyncingobject
LEFT JOIN (
    SELECT z_uuid FROM z_metadata
)
WHERE 
    title IS NOT NULL AND 
    zmarkedfordeletion != 1 AND
    lower(title) LIKE lower(?)
ORDER BY title ASC
`)

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

func (lite LiteDB) Query(sqlQuery string, sqlArg string) ([]map[string]string, error) {
    results := []map[string]string{}
    rows, err := lite.db.Query(sqlQuery, sqlArg)
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

func BuildSubtitleAddition(body string, bodyLower string, searchLower string) string {
    subtitleAddition := " | …"
    i := 0
    j := 0
    for i >= 0 && j >= 0 && len(subtitleAddition) < 400 {
        j = strings.Index(bodyLower[i:], searchLower)
        if j >= 0 {
            // Include context around match up to rb or next newline
            rb := min(len(body), i+j+len(searchLower)+25)
            nextNewline := strings.Index(body[i+j:rb], "\n")
            if nextNewline > 0 {
                rb = i+j+nextNewline
            }
            match := strings.ToValidUTF8(strings.Trim(body[i+j:rb], " "), "")
            subtitleAddition += match + "…"
            i = rb
        }
    }
    return subtitleAddition
}

func (lite LiteDB) QueryThenSearch(q string, search string) ([]map[string]string, error) {    
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
    
    searchLower := strings.ToLower(search)
    gzipHeader := []byte{31,139,8,0,0,0,0,0,0,19}
    bytesReader := bytes.NewReader(gzipHeader)
    gzipReader, err := gzip.NewReader(bytesReader)
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
        
        subtitleAddition := ""
        if len(search) > 0 {
            // Get note title and check for search string
            valTitle := columnPointers[0].(*interface{})
            title, ok := (*valTitle).(string)
            if !ok {
                continue
            }
            titleContains := strings.Contains(strings.ToLower(string(title)), searchLower)
            // Decompress note body data
            valBody := columnPointers[3].(*interface{})
            noteDataZippedBytes, ok := (*valBody).([]byte)
            if ok {
                bytesReader.Reset(noteDataZippedBytes)
                errReset := gzipReader.Reset(bytesReader)
                if errReset == nil {
                    noteBytes, errRead := ioutil.ReadAll(gzipReader)
                    if errRead == nil {
                        // Get plaintext of any tables in this note
                        valTableText := columnPointers[4].(*interface{})
                        tableText, ok := (*valTableText).(string)
                        if !ok {
                            tableText = ""
                        }
                        // Extract protobuf-format data from unzipped note
                        body := GetNoteBody(noteBytes)
                        body += tableText
                        bodyLower := strings.ToLower(body)
                        if strings.Contains(bodyLower, searchLower) {
                            subtitleAddition = BuildSubtitleAddition(body, bodyLower, searchLower)
                        } else if !titleContains {
                            continue
                        }
                    } else if !titleContains {
                        continue
                    }
                } else if !titleContains {
                    continue
                }
            } else if !titleContains {
                continue
            }
        }
        
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
        m[SubtitleKey] += subtitleAddition
        results = append(results, m)
    }
    return results, err
}

func RowToItem(row map[string]string, userQuery UserQuery) alfred.Item {
    return alfred.Item{
        Title:        row[TitleKey],
        Subtitle:     row[SubtitleKey],
        Arg:          row[ArgKey] + "?" + Escape(userQuery.WordString),
        QuicklookURL: nil,
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

func GetOrderPreference() string {
    if (os.Getenv("sortByDate") != "0") {
        return "modDate DESC"
    } else {
        return "lower(noteTitle) ASC"
    }
}

func GetSearchTitleRows(litedb LiteDB, userQuery UserQuery) ([]map[string]string, error) {
    likeString := "%%" + userQuery.WordString + "%%"
    orderBy := GetOrderPreference()
    var where string
    if (os.Getenv("searchFolders") != "0") {
        where = "WHERE (lower(folderTitle) || ' ' || lower(noteTitle)) LIKE lower(?)"
    } else {
        where = "WHERE lower(noteTitle) LIKE lower(?)"
    }
    rows, err := litedb.Query(fmt.Sprintf(NotesSQLTemplate, where, orderBy), likeString)
    if err != nil {
        return nil, err
    }
    return rows, nil
}

func GetSearchBodyRows(litedb LiteDB, userQuery UserQuery) ([]map[string]string, error) {
    orderBy := GetOrderPreference()
    where := ""
    rows, err := litedb.QueryThenSearch(fmt.Sprintf(NotesSQLTemplate, where, orderBy), userQuery.WordString)
    if err != nil {
        return nil, err
    }   

    return rows, nil
}

func GetSearchFolderRows(litedb LiteDB, userQuery UserQuery) ([]map[string]string, error) {
    likeString := "%%" + userQuery.WordString + "%%"
    rows, err := litedb.Query(FoldersSQLTemplate, likeString)
    if err != nil {
        return nil, err
    }
    return rows, nil
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
        
        userQuery := ParseUserQuery(os.Args[2])
        searchRows := []map[string]string{}
            
        if os.Args[1] == "title" {
            searchRows, err = GetSearchTitleRows(litedb, userQuery)
            PanicOnErr(err)
            
            if len(searchRows) == 0 {
                createItem, err := CreateNoteItem(userQuery)
                PanicOnErr(err)
                alfred.Add(*createItem)
            }
        } else if os.Args[1] == "body" {
            searchRows, err = GetSearchBodyRows(litedb, userQuery)
            PanicOnErr(err)
        } else if os.Args[1] == "folder" {
            searchRows, err = GetSearchFolderRows(litedb, userQuery)
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
