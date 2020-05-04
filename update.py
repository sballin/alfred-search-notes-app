import os
import time
import json
import plistlib
from distutils.version import StrictVersion
from urllib import request


def oneDaySinceLastCheck():
    '''
    Check whether it's been 24 hours since the last github query.
    We keep track of this through the modification time of this file,
    which is updated every time we query github.
    '''
    lastCheck = os.path.getmtime(__file__)
    if time.time() > 24*60*60 + lastCheck:
        os.system('touch ' + __file__)
        return True
    else:
        return False


def updateAvailable(latestVersion):
    '''
    Check whether the latest version is ahead of the current version.
    '''
    with open('info.plist', 'rb') as f:
        currentVersion = plistlib.load(f)['version']    
    if StrictVersion(currentVersion) < StrictVersion(latestVersion):
        return True
    else:
        return False


def userWantsUpdate(updateUrl):
    retval = os.system("""
osascript -e 'display dialog "An update is available for the Search Notes workflow. Press OK to download and open this file:

%s

Daily update checks can be disabled by editing the workflow."'
""" % updateUrl)
    if retval == 0:
        return True
    else:
        return False
        
        
def update(updateUrl):
    '''
    Download and open new version of workflow.
    '''
    r = request.urlopen(updateUrl, timeout=60)
    if r.status == 200:
        updateFile = '/tmp/Search.Notes.alfredworkflow'
        with open(updateFile, 'wb') as f:
            f.write(r.read())
        os.system('open ' + updateFile)
    

try:
    if oneDaySinceLastCheck():
        latestUrl = 'https://api.github.com/repos/sballin/alfred-search-notes-app/releases/latest'
        r = request.urlopen(latestUrl, timeout=60)
        if r.status == 200:
            body = r.read()
            latest = json.loads(body.decode('utf-8'))
            latestVersion = latest['tag_name']
            updateUrl = latest['assets'][0]['browser_download_url']
            
            if updateAvailable(latestVersion): 
                if userWantsUpdate(updateUrl):
                    update(updateUrl)
except:
    pass
    
