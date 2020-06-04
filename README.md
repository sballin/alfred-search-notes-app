# Search Notes.app with Alfred

<p align="center">
  <img src="https://user-images.githubusercontent.com/2719004/71554521-5c635800-2a31-11ea-97db-1b7c41aaf408.png" width="654" title="screenshot">
</p>

## Usage

### Search keywords

*   **n** to search note titles (this also lets you find notes using the pattern "[folder name] [note name]")
*   **nb** to include note body text in search
*   **nf** to search note folder names

### Result actions

*   **enter** to open the note/folder, or create a new note if none was found
*   **shift+enter** to search for your Alfred query using the Notes in-app search
*   **cmd+enter** to copy the note body to the clipboard
*   **alt+enter** to copy a notes:// URL to the clipboard that can be opened with the included `Note Opener/Note Opener.app`

## Install

If you're on macOS Catalina with Alfred 4, download the [latest version](https://github.com/sballin/alfred-search-notes-app/releases/latest/download/Search.Notes.alfredworkflow) of the workflow.

### Authorizing the binary

The "search" binary needs to be manually authorized to run on your computer. Right-click the Search Notes workflow, click "Open in Finder", open the "search" folder, right-click the "search" binary, click "Open", and you should be all set after it runs once in Terminal.

### Older versions

If you encounter problems or are on an older version of macOS/Alfred, try an older version like [1.4.3](https://github.com/sballin/alfred-search-notes-app/releases/tag/1.4.3).

### Stay up to date

By default, the workflow checks Github for updates every 24 hours. This can be disabled by removing the connections to the "Check for updates" block in the workflow.

You can also be notified of new releases by watching this repo or subscribing to the Alfred forum [thread](https://www.alfredforum.com/topic/11716-search-appleicloud-notes/).

## Customize

Result ordering and title+folder search behavior can be controlled using the [environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

## Contributors

Thank you to all who have submitted pull requests, and to [drgrib](https://github.com/drgrib) for allowing me to build off the [alfred-bear](https://github.com/drgrib/alfred-bear) workflow.

## Donate

If you enjoy using this workflow, consider [donating](http://paypal.me/sbballin)!
