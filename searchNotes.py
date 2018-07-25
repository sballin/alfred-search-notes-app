#!/usr/bin/python
import sqlite3
import json


# Open notes database
home = '/'.join(__file__.split('/')[:3])
conn = sqlite3.connect(home + '/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite')
c = conn.cursor()

# Get uuid string required in full id
c.execute("SELECT z_uuid FROM z_metadata")
uuid = str(c.fetchone()[0])

# Get tuples of note title, folder code, snippet, modification date, and id number 
c.execute("""SELECT ztitle1,zfolder,zsnippet,zmodificationdate1,z_pk
FROM ziccloudsyncingobject 
WHERE ztitle1 IS NOT NULL AND zmarkedfordeletion IS NOT 1""")
matches = c.fetchall()
# Sort by modification date
matches = sorted(matches, key=lambda m: m[3], reverse=True)

# Get ordered lists of folder codes and folder names
c.execute("""SELECT z_pk,ztitle2 
FROM ziccloudsyncingobject 
WHERE ztitle2 IS NOT NULL AND zmarkedfordeletion IS NOT 1""")
folderCodes, folderNames = zip(*c.fetchall())

conn.close()

# Alfred results: title=note title, arg=id to pass on, subtitle="folder name | snippet"
items = [{"title":    m[0], 
          "arg":      "x-coredata://" + uuid + "/ICNote/p" + str(m[4]), 
          "subtitle": folderNames[folderCodes.index(m[1])] + " | " + m[2][:80]}
         for m in matches]

# Custom icons for folder names that start with corresponding emoji
icons = [u'\ud83d\udcd3', u'\ud83d\udcd5', u'\ud83d\udcd7', u'\ud83d\udcd8', u'\ud83d\udcd9']
for i in items:
    if any(x in i['subtitle'] for x in icons):
        subtitle = i['subtitle']
        icon = subtitle[:2]
        i['subtitle'] = subtitle[3:]
        i['icon'] = {'type': 'image', 'path': icon.encode('raw_unicode_escape') + '.png'}

output = {"items": items}
print json.dumps(output)
