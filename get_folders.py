#!/usr/bin/env python3
import json
from get_notes import readDatabase


uuid, dbItems, folders = readDatabase()

items = []
for folderCode in folders:
    name = folders[folderCode]
    if name != 'New Folder':
        items.append({'title': name,
                      'subtitle': 'Folder',
                      'arg':'x-coredata://' + uuid + '/ICFolder/p' + str(folderCode)})

print(json.dumps({'items': items}))
