#!/usr/bin/python
import sqlite3
import zlib
import re
import os


def extractNoteBody(data):
    try:
        # Strip weird characters, title & weird header artifacts, 
        # and replace line breaks with spaces
        data = zlib.decompress(data, 16+zlib.MAX_WBITS).split('\x1a\x10', 1)[0]
        data = re.sub(r'[\x00-\x08\x0b\x0c\x0e-\x1f\x7f-\xff]|^  ', '', data)
        return re.sub('^.*\n|\n', ' ', data)
    except Exception as e:
        return 'Note body could not be extracted: {}'.format(e)


# Sort matches by title or modification date, option to search titles only.
# Edit with Alfred workflow environment variable.
sortId = 2 if os.getenv('sortByDate') == '1' else 0
searchTitlesOnly = (os.getenv('searchTitlesOnly') == '1')
sortInReverse = (sortId == 2)

# Custom icons to look for in folder names
icons = [u'\ud83d\udcd3', u'\ud83d\udcd5', u'\ud83d\udcd7', u'\ud83d\udcd8', 
         u'\ud83d\udcd9']
         
try:
    # Open notes database
    home = os.path.expanduser('~')
    db = home + '/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite'
    conn = sqlite3.connect(db)
    c = conn.cursor()

    # Get uuid string required in full id
    c.execute('SELECT z_uuid FROM z_metadata')
    uuid = str(c.fetchone()[0])

    # Get tuples of note title, folder code, modification date, & id#
    c.execute("""SELECT t1.ztitle1,t1.zfolder,t1.zmodificationdate1,
                        t1.z_pk,t1.znotedata,t2.zdata,t2.z_pk
                 FROM ziccloudsyncingobject AS t1
                 INNER JOIN zicnotedata AS t2
                 ON t1.znotedata = t2.z_pk
                 WHERE t1.ztitle1 IS NOT NULL 
                       AND t1.zmarkedfordeletion IS NOT 1""")
    # Get data and check for d[5] because a New Note with no body can trip us up
    dbItems = [d for d in c.fetchall() if d[5]]
    dbItems = sorted(dbItems, key=lambda d: d[sortId], reverse=sortInReverse)

    # Get ordered lists of folder codes and folder names
    c.execute("""SELECT z_pk,ztitle2 FROM ziccloudsyncingobject
                 WHERE ztitle2 IS NOT NULL 
                       AND zmarkedfordeletion IS NOT 1""")
    folderCodes, folderNames = zip(*c.fetchall())

    conn.close()
    openedDatabase = True
except:
    openedDatabase = False
         
if openedDatabase:
    # Alfred results: title = note title, arg = id to pass on, subtitle = folder name, 
    # match = note contents from gzipped database entries after stripping footers.
    items = [{} for d in dbItems]
    gotOneRealNote = False
    for i, d in enumerate(dbItems):
        try:
            folderName = folderNames[folderCodes.index(d[1])]
            if folderName == 'Recently Deleted':
                continue
            body = extractNoteBody(d[5])
            subtitle = folderName + '  |' + body[:100]
            match = u'{} {} {}'.format(d[0], folderName, '' if searchTitlesOnly else body)
            
            # Custom icons for folder names that start with corresponding emoji
            if any(x in subtitle[:2] for x in icons):
                iconText = subtitle[:2].encode('raw_unicode_escape')
                subtitle = subtitle[3:]
                icon = {'type': 'image', 'path': 'icons/' + iconText + '.png'}
            else:
                icon = {'type': 'default'}
            
            items[i] = {'title': d[0],
                        'subtitle': subtitle,
                        'arg': 'x-coredata://' + uuid + '/ICNote/p' + str(d[3]),
                        'match': match,
                        'icon': icon}
            gotOneRealNote = True
        except Exception as e:
            items[i] = {'title': 'Error getting note', 'subtitle': str(e)}

if openedDatabase and gotOneRealNote:
    import json
    print json.dumps({'items': items})
else:
    import subprocess
    print subprocess.check_output(os.path.dirname(__file__) 
                                  + '/searchNoteTitles.applescript')
