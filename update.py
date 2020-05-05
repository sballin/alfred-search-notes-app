import os
import time
import json
import plistlib
from distutils.version import StrictVersion


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


def userWantsUpdate(updateNotes):
    '''
    Show user a confirmation dialog.
    '''
    retval = os.system("""
osascript -e 'display dialog "An update is available for the Search Notes workflow. Press OK to download and open it. Notes for this release:

%s

Daily update checks can be disabled by editing the workflow."' 2>/dev/null
""" % updateNotes)
    if retval == 0:
        return True
    else:
        return False
        
        
def update(updateUrl):
    '''
    Download and open new version of workflow.
    '''
    updateFile = '/tmp/Search.Notes.alfredworkflow'
    # --location is required in order to follow redirects
    curlRet = os.system('curl --silent --location --output %s %s' % (updateFile, updateUrl))
    openRet = 1
    if curlRet == 0:
        openRet = os.system('open ' + updateFile)
    if curlRet != 0 or openRet != 0:
        os.system("osascript -e 'display dialog \"The Search Notes workflow failed to update.\"' 2>/dev/null")
    

if oneDaySinceLastCheck():
    latestUrl = 'https://api.github.com/repos/sballin/alfred-search-notes-app/releases/latest'
    latestFile = 'latest_release.json'
    retval = os.system('curl --silent --max-time 30 --output %s %s' % (latestFile, latestUrl))
    if retval == 0:
        with open(latestFile, 'r') as f:
            latest = json.load(f)
        latestVersion = latest['tag_name']
        updateNotes = latest['body']
        updateUrl = latest['assets'][0]['browser_download_url']
        
        if updateAvailable(latestVersion): 
            if userWantsUpdate(updateNotes):
                update(updateUrl)
