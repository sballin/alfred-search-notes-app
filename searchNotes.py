#!/usr/bin/python
import sqlite3
import json
import zlib
import re
import os

# Sort matches by title or modification date, option to search titles only.
# Edit with Alfred workflow environment variable.
sortId = 2 if os.getenv('sortByDate') == '1' else 0
searchTitlesOnly = (os.getenv('searchTitlesOnly') == '1')
sortInReverse = (sortId == 2)

# Open notes database
home = '/'.join(__file__.split('/')[:3])
db = home + '/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite'
conn = sqlite3.connect(db)
c = conn.cursor()

# Get uuid string required in full id
c.execute('SELECT z_uuid FROM z_metadata')
uuid = str(c.fetchone()[0])

# Get tuples of note title, folder code, modification date, & id#
# 432 is for at least one person the zfolder id for 'Recently Deleted'
c.execute("""SELECT t1.ztitle1,t1.zfolder,t1.zmodificationdate1,
                    t1.z_pk,t1.znotedata,t2.zdata,t2.z_pk
FROM ziccloudsyncingobject AS t1
INNER JOIN zicnotedata AS t2
ON t1.znotedata = t2.z_pk
WHERE t1.ztitle1 IS NOT NULL AND t1.zfolder IS NOT 432 
      AND t1.zmarkedfordeletion IS NOT 1""")
# Get and check for d[5] because a New Note with no body can trip us up
dbItems = [d for d in c.fetchall() if d[5]]
dbItems = sorted(dbItems, key=lambda d: d[sortId], reverse=sortInReverse)

# Get ordered lists of folder codes and folder names
c.execute("""SELECT z_pk,ztitle2 FROM ziccloudsyncingobject
WHERE ztitle2 IS NOT NULL AND zmarkedfordeletion IS NOT 1""")
folderCodes, folderNames = zip(*c.fetchall())

conn.close()

# Custom icons to look for in folder names
icons = [u'\ud83d\udcd3', u'\ud83d\udcd5', u'\ud83d\udcd7', u'\ud83d\udcd8', 
         u'\ud83d\udcd9']
         
# Alfred results: title = note title, arg = id to pass on, subtitle = folder name, 
# match = note contents from gzipped database entries after stripping footers.
items = [{} for d in dbItems]
for i, d in enumerate(dbItems):
    # In body, strip weird characters, title & weird header artifacts, 
    # and replace line breaks with spaces
    body = zlib.decompress(d[5], 16+zlib.MAX_WBITS).split('\x1a\x10', 1)[0]
    body = re.sub(r'[\x00-\x08\x0b\x0c\x0e-\x1f\x7f-\xff]', '', body)
    body = re.sub('^.*\n', ' ', body)
    body = re.sub('\n', ' ', body)
    body = re.sub('^  ', '', body)
    
    subtitle = folderNames[folderCodes.index(d[1])] + " | " + body[:100]
    
    # Custom icons for folder names that start with corresponding emoji
    if any(x in subtitle[:2] for x in icons):
        iconText = subtitle[:2].encode('raw_unicode_escape')
        subtitle = subtitle[3:]
        icon = {'type': 'image', 'path': 'icons/' + iconText + '.png'}
    else:
        icon = {'type': 'fileicon', 'path': '/Applications/Notes.app'}
    
    items[i] = {'title': d[0],
                'subtitle': subtitle,
                'arg': 'x-coredata://' + uuid + '/ICNote/p' + str(d[3]),
                'match': d[0] if searchTitlesOnly else d[0] + body,
                'icon': icon}

print json.dumps({'items': items}, sort_keys=True)
