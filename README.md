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
*   **alt+enter** to copy a note URL to the clipboard (see the section on [note linking](#note-linking))

## Install

If you're on macOS Catalina with Alfred 4, download the [latest version](https://github.com/sballin/alfred-search-notes-app/releases/latest/download/Search.Notes.alfredworkflow) of the workflow.

### Required setup

1. Try searching for a note
2. When an error message appears, make sure to click "Cancel" in the first dialog, then click "Open System Preferences" in the second dialog
3. Near the bottom of the pane click the "Open Anyway" button for the "search_notes" binary

<img src="https://user-images.githubusercontent.com/2719004/101306597-536b3100-3813-11eb-861e-8cbf74277255.png" width="890">
<img src="https://user-images.githubusercontent.com/2719004/99889399-af826280-2c22-11eb-84dc-0c7972010dfa.png" width="890">

This workflow currently doesn't support notes stored with Google or other internet accounts. Please make sure either iCloud or On My Mac is selected as the default account in the preferences of Notes.app.

## Customize search behavior

Result ordering and title+folder search behavior can be controlled using the [environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

## Note linking

You can generate links to any of your notes and use them on macOS or iOS.

* macOS
    * Copy note URL by pressing alt+enter on an Alfred result
    * Open note URL with "Open Notes URL.app" (included with this workflow) which works automatically when clicking a link
* iOS (must open these links on iOS to install)
    * Copy note URL [shortcut](https://www.icloud.com/shortcuts/556aba9692d64694b7073345ea224dc2) (see image below for usage instructions)
    * Open note URL [shortcut](https://www.icloud.com/shortcuts/825f1ac1d09149689c9d2406c24aef9e) works automatically when clicking a link

<img src="https://user-images.githubusercontent.com/2719004/101307089-a4c7f000-3814-11eb-89cf-f2673b86c543.png" width="890">

## Updates

By default, the workflow checks Github for updates every 24 hours. This can be disabled by removing the connections to the "Check for updates" block in the workflow.

## Compile

If you want to compile the binary yourself, you can go into the "search" folder and do `make`.

## Contributors

Big thanks to...

* All who have submitted pull requests
* drgrib for allowing me to build off the [alfred-bear](https://github.com/drgrib/alfred-bear) workflow
* threeplanetssoftware for the [apple_cloud_notes_parser](https://github.com/threeplanetssoftware/apple_cloud_notes_parser) from which I copied the protobuf handling
* [lslz627](https://github.com/lslz627) for help with protobuf and tables
* [Artem Chistyakov](https://temochka.com/blog/posts/2020/02/22/linking-to-apple-notes.html) for a much improved way to create links to notes

## Donate

If you enjoy using this workflow, consider [donating](http://paypal.me/sbballin)!
