#!/usr/bin/python
import sqlite3
import json
import zlib
import re
import os

# Sort matches by title or modification date. Edit w/environment variable.
sortId = 3 if os.getenv('sortBy') != 'Title' else 0
sortInReverse = sortId is 3

# Open notes database
home = '/'.join(__file__.split('/')[:3])
conn = sqlite3.connect(
    home + '/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite')
c = conn.cursor()

# Get uuid string required in full id
c.execute("SELECT z_uuid FROM z_metadata")
uuid = str(c.fetchone()[0])

# Get tuples of note title, folder code, snippet, modification date, & id#
# 432 is the zfolder id for 'Recently Deleted'
c.execute("""SELECT t1.ztitle1,t1.zfolder,t1.zsnippet,t1.zmodificationdate1,t1.z_pk,t1.znotedata,t2.zdata,t2.z_pk
FROM ziccloudsyncingobject AS t1
INNER JOIN zicnotedata AS t2
ON t1.znotedata = t2.z_pk
WHERE t1.ztitle1 IS NOT NULL AND t1.zfolder IS NOT 432 AND t1.zmarkedfordeletion IS NOT 1""")
matches = c.fetchall()
matches = sorted(matches, key=lambda m: m[sortId], reverse=sortInReverse)

# Get ordered lists of folder codes and folder names
c.execute("""SELECT z_pk,ztitle2
FROM ziccloudsyncingobject
WHERE ztitle2 IS NOT NULL AND zmarkedfordeletion IS NOT 1""")
folderCodes, folderNames = zip(*c.fetchall())

conn.close()

# Alfred results: title = note title, arg = id to pass on, subtitle = folder name, match = the note contents
items = [{"title": m[0],
          "arg": "x-coredata://" + uuid + "/ICNote/p" + str(m[4]),
          "subtitle": folderNames[folderCodes.index(m[1])],
          #  + ("  |  " + m[2] if type(m[2]) is unicode and len(m[2]) > 0 else ""),
          #
          #  decompress gzipped notes from the sqlite database, strip out gobbledygook footers.
          "match": zlib.decompress(m[6], 16+zlib.MAX_WBITS).split('\x1a\x10', 1)[0]}
         for m in matches]

# Do further clean up and additions to the match and subtitle fields.
for i, item in enumerate(items):
    # strip weird characters, title & weird header artifacts,
    # replace line breaks with spaces.
    txt = re.sub('^  ', '', re.sub('\n', ' ', re.sub('^.*\n', ' ', re.sub(r'[\x00-\x08\x0b\x0c\x0e-\x1f\x7f-\xff]', '', items[i]['match']))))
    items[i]['match'] = items[i]['title'] + " " + items[i]['subtitle'] + " " + txt
    items[i]['subtitle'] += "  |  " + txt[:100]

# Custom icons for folder names that start with corresponding emoji
icons = [u'\ud83d\udcd3', u'\ud83d\udcd5', u'\ud83d\udcd7', u'\ud83d\udcd8', u'\ud83d\udcd9']
for i in items:
    if any(x in i['subtitle'] for x in icons):
        subtitle = i['subtitle']
        icon = subtitle[:2]
        i['subtitle'] = subtitle[3:]
        i['icon'] = {'type': 'image', 'path': 'icons/' + icon.encode('raw_unicode_escape') + '.png'}

output = {"items": items}
print json.dumps(output, sort_keys=True, indent=4, separators=(',', ': '))
