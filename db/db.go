package db

import (
	"database/sql"
	"os/user"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DbPath = "~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite"

	TitleKey  = "noteTitle"
	FolderKey   = "folderName"
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
