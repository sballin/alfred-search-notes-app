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
*   **alt+enter** to copy a link to the note to the clipboard

### Note linking

You can generate links to any of your notes and use them on macOS or iOS. Copy the note URL by pressing alt+enter on an Alfred result or paste it using the snippet. This will generate two links. The first one works on macOS Big Sur (11) and newer, and the second one works on iOS.

## Install

Download the [latest version](https://github.com/sballin/alfred-search-notes-app/releases/latest/download/Search.Notes.alfredworkflow) of the workflow if you're on the latest macOS with Alfred >= 4. For older versions of macOS, you may need to use [older versions](https://github.com/sballin/alfred-search-notes-app/releases) of the workflow.

### Required setup

1. Try searching for a note
2. If a warning dialog appears (see image below), click "Cancel" rather than "Move to Trash". Then open System Preferences > Security & Privacy and click the "Open Anyway" button near the bottom
4. Approve additional requests for permission as they appear
5. If there are any other issues, please follow the advice under [troubleshooting](#troubleshooting), if present, is enabled

<img src="https://user-images.githubusercontent.com/2719004/123869471-0b227600-d8ff-11eb-8c20-6537055b1336.png" width="890">

### Troubleshooting

If you get a permission-related error, especially after installing updates to macOS, try disabling and re-enabling the permissions shown below, especially full disk access for Alfred. If that doesn't work, please look through [common issues](https://github.com/sballin/alfred-search-notes-app/issues?q=) before submitting a new one.

<img src="https://github.com/sballin/alfred-search-notes-app/assets/2719004/566002dd-f3be-4b98-88b7-16d6e185d531" width="890">

### Email notes are not supported

This workflow doesn't support notes stored with Google or other internet accounts. Please make sure either iCloud or On My Mac is selected as the default account in the preferences of Notes.app.

## Customize search behavior

Result ordering and title+folder search behavior can be controlled using the [environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

## Updates

By default, the workflow checks Github for updates every 24 hours. This can be disabled by removing the connections to the "Check for updates" block in the workflow.

## Compile

If you want to compile the binary yourself, you can go into the "search_notes" folder and do `make`.

## Contributors

Big thanks to...

* All who have submitted pull requests
* drgrib for allowing me to build off the [alfred-bear](https://github.com/drgrib/alfred-bear) workflow
* threeplanetssoftware for the [apple_cloud_notes_parser](https://github.com/threeplanetssoftware/apple_cloud_notes_parser) from which I copied the protobuf handling
* [lslz627](https://github.com/lslz627) for help with protobuf and tables
* [Artem Chistyakov](https://temochka.com/blog/posts/2020/02/22/linking-to-apple-notes.html) for a much improved way to create links to notes
* [vitorgalvao](https://github.com/vitorgalvao) for the OneUpdater code

## Donate

If you enjoy using this workflow, consider [donating](http://paypal.me/sbballin)!
