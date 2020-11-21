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
*   **alt+enter** to copy a note:// URL to the clipboard that can be opened with the included `Open Notes URL/Open Notes URL.app`

## Install

If you're on macOS Catalina with Alfred 4, download the [latest version](https://github.com/sballin/alfred-search-notes-app/releases/latest/download/Search.Notes.alfredworkflow) of the workflow.

### Required setup

1. Try searching for a note
2. When an error message appears, press the "Open System Preferences" button
3. Near the bottom of the pane click the "Open Anyway" button for the "search_notes" binary

<img src="https://user-images.githubusercontent.com/2719004/99889399-af826280-2c22-11eb-84dc-0c7972010dfa.png" width="890">

This workflow currently doesn't support notes stored with Google or other internet accounts. Please make sure either iCloud or On My Mac is selected as the default account in the preferences of Notes.app.

### Customize

Result ordering and title+folder search behavior can be controlled using the [environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

### Updates

By default, the workflow checks Github for updates every 24 hours. This can be disabled by removing the connections to the "Check for updates" block in the workflow.

### Compile

If you want to compile the binary yourself, you can go into the "search" folder and do `make`.

## Contributors

Thank you to 

* All who have submitted pull requests
* drgrib for allowing me to build off the [alfred-bear](https://github.com/drgrib/alfred-bear) workflow
* threeplanetssoftware for the [apple_cloud_notes_parser](https://github.com/threeplanetssoftware/apple_cloud_notes_parser) from which I copied the protobuf handling
* [lslz627](https://github.com/lslz627) for help with protobuf and tables

## Donate

If you enjoy using this workflow, consider [donating](http://paypal.me/sbballin)!
