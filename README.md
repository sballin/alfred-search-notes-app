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

*   **enter** to open the note/folder
*   **cmd+enter** to copy the note body to the clipboard
*   **alt+enter** to copy notes://[title of note] to the clipboard

Use `Note Opener/Note Opener.app` to make notes:// urls open the relevant note when clicked.

## Install

If you're on macOS Catalina, download the [latest version](https://github.com/sballin/alfred-search-notes-app/releases/latest/download/Search.Notes.alfredworkflow) of the workflow. If you haven't used python 3 on your computer before, you may need to install the Xcode developer tools by running the following command in a terminal:

    xcode-select --install

If that doesn't work for you, or you're on an older macOS, try version [1.4.3](https://github.com/sballin/alfred-search-notes-app/releases/tag/1.4.3).

### Stay up to date

By default, the workflow checks Github for updates every 24 hours. This can be disabled by removing the connections to the "Check for updates" block in the workflow.

You can also be notified of new releases by watching this repo or subscribing to the Alfred forum [thread](https://www.alfredforum.com/topic/11716-search-appleicloud-notes/).

## Customize

### Result ordering

By default, results are ordered based on the modification date of the note. If you want to order results alphabetically, change the [environment variable](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

### Icons

Custom icons will show up when a folder name begins with an emoji, like `📗 Misc` in the screenshot above. Add your own icons to the workflow's `icons` folder and tweak `get_notes.py` to see them in Alfred. The default icons are from [Emojitwo](https://emojitwo.github.io/).

## Contributors

Thank you to all who have submitted pull requests, and to [drgrib](https://github.com/drgrib) for allowing me to repurpose the [alfred-bear](https://github.com/drgrib/alfred-bear) workflow.

## Donate

If you enjoy using this workflow, consider [donating](http://paypal.me/sbballin)!
