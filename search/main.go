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
    "golang.org/x/text/unicode/norm"

    _ "github.com/mattn/go-sqlite3"
    "github.com/sballin/alfred-search-notes-app/alfred"
)

const (
    DbPath = "~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"

    TitleKey  = "title"
    SubtitleKey = "subtitle"
    ArgKey = "URL"
    BodyKey = "noteBodyZipped"

    NotesSQLTemplate = `
SELECT 
    noteTitle as title,
    folderTitle as subtitle,
    'x-coredata://' || z_uuid || '/ICNote/p' || xcoreDataID as URL,
    noteBodyZipped 
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

    FoldersSQLTemplate = `
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
    
    gotReaders := false
    var bytesReader *bytes.Reader
    var gzipReader *gzip.Reader

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
        
        // Skip adding item if note body does not contain search string
        val := columnPointers[len(cols)-1].(*interface{})
        // Type assertion required by bytesReader
        noteBodyZippedBytes, ok := (*val).([]byte)
        if !ok {
            continue
        }
        if gotReaders {
            bytesReader.Reset(noteBodyZippedBytes)
            gzipReader.Reset(bytesReader)
        } else {
            bytesReader = bytes.NewReader(noteBodyZippedBytes)
            gzipReader, err = gzip.NewReader(bytesReader)
            if err != nil {
                continue
            }
            gotReaders = true
        }
        body, err := ioutil.ReadAll(gzipReader)
        if err != nil {
            continue
        }
        if !strings.Contains(strings.ToLower(string(body)), strings.ToLower(search)) {
            continue
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
