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

Download the latest .alfredworkflow file in [releases](https://github.com/sballin/alfred-search-notes-app/releases) and open it.

Versions >=2.0.0 should work out of the box on Catalina and can be modified to work on earlier OSes. Versions <=1.4.3 should work out of the box on Mojave and earlier macOS versions. 

## Customize

### Result ordering

By default, results are ordered based on the modification date of the note. If you want to order results alphabetically, change the [environment variable](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

### Icons

Custom icons will show up when a folder name begins with an emoji, like `📗 Misc` in the screenshot above. Add your own icons to the workflow's `icons` folder and tweak `get_notes.py` to see them in Alfred. The default icons are from [Emojitwo](https://emojitwo.github.io/).
