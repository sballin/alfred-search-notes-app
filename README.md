# Search Notes.app with Alfred

<p align="center">
  <img src="https://user-images.githubusercontent.com/2719004/45398015-8a404580-b610-11e8-9879-93c9c002a375.png" width="654" title="screenshot">
</p>

## Usage

Type `n[part of note]` and press enter.

## Install

Download the Alfred workflow file from Releases and open it.

## Customize

### Search 

By default, this searches the title + full text of the notes and orders results based on the last modification date of the note. If you want to search titles only or order results alphabetically, change the [environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment).

### Icons

Icons are from [Emojitwo](https://emojitwo.github.io/) and will show up when they are the first character in the name of a folder, like `📘 GPI` in the screenshot above. Add your own icons to the workflow's `icons` folder and tweak `searchNotes.py` to see them in Alfred.

## Compatibility

The default search method is only tested in High Sierra. If it's not working for you and/or you're on a different version of macOS, try the AppleScript search methods using the keywords `a` and `b`.

