# Search Notes.app with Alfred

### Search titles or create a new note if none was found

<img src="https://user-images.githubusercontent.com/2719004/83949726-62850e00-a7f3-11ea-99a7-48f8c67cd480.png" width="654">
  
<img src="https://user-images.githubusercontent.com/2719004/83949516-403ec080-a7f2-11ea-940c-1813559ce462.png" width="654">

### Search note titles and bodies

<img src="https://user-images.githubusercontent.com/2719004/83949619-e094e500-a7f2-11ea-8802-7856620d4ec8.png" width="654">

### Search folder names

<img src="https://user-images.githubusercontent.com/2719004/83949622-e25ea880-a7f2-11ea-92fa-b2250e574402.png" width="654">

### Result actions

*   **enter** to open the note/folder or create a new note if none was found
*   **shift+enter** to search for your Alfred query using the Notes in-app search 
*   **cmd+enter** to copy the note body to the clipboard
*   **alt+enter** to copy a notes:// URL to the clipboard that can be opened with the included `Note Opener/Note Opener.app`

## Install

If you're on macOS Catalina with Alfred 4, download the [latest version](https://github.com/sballin/alfred-search-notes-app/releases/latest/download/Search.Notes.alfredworkflow) of the workflow.

### Required setup

The "search" binary needs to be manually authorized to run on your computer. After installing, right-click the Search Notes workflow in Alfred preferences, click "Open in Finder", open the "search" folder, right-click the "search" binary, click "Open", and you should be all set after it runs once in Terminal.

### Customize

Result ordering and title+folder search behavior can be controlled using the [environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

### Updates

By default, the workflow checks Github for updates every 24 hours. This can be disabled by removing the connections to the "Check for updates" block in the workflow.

### Compile

If you want to compile the binary yourself, you can go into the "search" folder and do `go build`.

## Contributors

Thank you to all who have submitted pull requests, and to [drgrib](https://github.com/drgrib) for allowing me to build off the [alfred-bear](https://github.com/drgrib/alfred-bear) workflow.

I'm very new to Go, so if you see anything that can be improved, don't hesitate to submit a pull request.

## Donate

If you enjoy using this workflow, consider [donating](http://paypal.me/sbballin)!
