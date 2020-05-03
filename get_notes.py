import sqlite3
import zlib
import os
import json


def extractNoteBody(data):
    # Decompress
    try:
        data = zlib.decompress(data, 16+zlib.MAX_WBITS).split(b'\x1a\x10', 1)[0]
    except zlib.error as e:
        return 'Encrypted note'
    # Find magic hex and remove it 
    # Source: https://github.com/threeplanetssoftware/apple_cloud_notes_parser
    index = data.index(b'\x08\x00\x10\x00\x1a')
    index = data.index(b'\x12', index) # starting from index found previously
    # Read from the next byte after magic index
    data = data[index+1:]
    # Convert from bytes object to string
    text = data.decode('utf-8', errors='ignore')
    # Remove title
    lines = text.split('\n')
    if len(lines) > 1:
        return '\n'.join(lines[1:])
    else:
        return ''


def fixStringEnds(text):
    """
    Shortening the note body for a one-line preview can chop two-byte unicode
    characters in half. This method fixes that.
    """
    # This method can chop off the last character of a short note, so add a dummy
    text = text + '.'
    # Source: https://stackoverflow.com/a/30487177
    pos = len(text) - 1
    while pos > -1 and ord(text[pos]) & 0xC0 == 0x80:
        # Character at pos is a continuation byte (bit 7 set, bit 6 not)
        pos -= 1
    return text[:pos]
    
    
def newlinesToSpace(text):
    """
    Replace any number of newlines with a single space character.
    """
    return ' '.join(text.replace('\n', ' ').split())


def readDatabase():
    # Open notes database read-only 
    home = os.path.expanduser('~')
    db = home + '/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite'
    conn = sqlite3.connect('file:' + db + '?mode=ro', uri=True)
    c = conn.cursor()

    # Get uuid string required in x-coredata URL
    c.execute('SELECT z_uuid FROM z_metadata')
    uuid = str(c.fetchone()[0])

    # Get note rows
    c.execute("""SELECT c.ztitle1,            -- note title (str)
                        c.zfolder,            -- folder code (int)
                        c.zmodificationdate1, -- modification date (float)
                        c.z_pk,               -- note id for x-coredata URL (int)
                        n.zdata               -- note body text (str)
                 FROM ziccloudsyncingobject AS c
                 INNER JOIN zicnotedata AS n
                 ON c.znotedata = n.z_pk -- note id (int) distinct from x-coredata one
                 WHERE c.ztitle1 IS NOT NULL AND 
                       c.zfolder IS NOT NULL AND            -- fix issues/21
                       c.zmodificationdate1 IS NOT NULL AND -- fix issues/20
                       c.z_pk IS NOT NULL AND
                       n.zdata IS NOT NULL AND              -- fix issues/3
                       c.zmarkedfordeletion IS NOT 1""")
    dbItems = c.fetchall()

    # Get folder rows
    c.execute("""SELECT z_pk,   -- folder code
                        ztitle2 -- folder name
                 FROM ziccloudsyncingobject
                 WHERE ztitle2 IS NOT NULL AND 
                       zmarkedfordeletion IS NOT 1""")
    folders = {code: name for code, name in c.fetchall()}

    conn.close()
    return uuid, dbItems, folders


def getNotes(searchBodies=False):
    # Custom icons to look for in folder names
    icons = {'📓': 'notebook.png', 
             '📕': 'redbook.png', 
             '📗': 'greenbook.png', 
             '📘': 'bluebook.png', 
             '📙': 'orangebook.png'}

    # Read Notes database and get contents
    uuid, dbItems, folders = readDatabase()
    
    # Sort matches by title or modification date (read Alfred environment variable)
    if os.getenv('sortByDate') == '0':
        sortId = 0
        sortInReverse = False
    else:
        sortId = 2
        sortInReverse = True
    dbItems = sorted(dbItems, key=lambda d: d[sortId], reverse=sortInReverse)

    # Alfred results: title = note title, arg = id to pass on, subtitle = folder name, 
    # match = note contents from gzipped database entries after stripping footers.
    items = [{} for d in dbItems]
    for i, d in enumerate(dbItems):
        title, folderCode, modDate, noteId, bodyData = d
        folderName = folders[folderCode]
        if folderName == 'Recently Deleted':
            continue
            
        try:
            body = extractNoteBody(bodyData)
        except:
            body = ''
            
        try:
            # Replace any number of \ns with a single space for note body preview
            bodyPreview = newlinesToSpace(body[:100])
        except:
            bodyPreview = ''
            
        if bodyPreview:
            subtitle = folderName + ' | ' + bodyPreview
        else:
            subtitle = folderName
            
        if searchBodies:
            bodyMatch = newlinesToSpace(body)
            match = u'{} {} {}'.format(folderName, title, bodyMatch)
        else:
            match = u'{} {}'.format(folderName, title)
            
        try:
            # Custom icons for folder names that start with corresponding emoji
            if folderName[0] in icons.keys():
                icon = {'type': 'image', 'path': 'icons/' + icons[folderName[0]]}
                subtitle = subtitle[2:]
            else:
                icon = {'type': 'default'}
        except: 
            icon = {'type': 'default'}
            
        try:
            subtitle = fixStringEnds(subtitle)
        except:
            subtitle = folderName
        
        items[i] = {'title': title,
                    'subtitle': subtitle,
                    'arg': 'x-coredata://' + uuid + '/ICNote/p' + str(noteId),
                    'match': match,
                    'icon': icon}

    return json.dumps({'items': items}, ensure_ascii=True)


if __name__ == '__main__':
    print(getNotes(searchBodies=False))
