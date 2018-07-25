#!/usr/bin/osascript
on alfred_script(argv)
	tell application "Notes"
		set noteRefs to a reference to every note in default account
		repeat with noteRef in noteRefs
			if name of noteRef contains argv then
				show noteRef
				activate
				exit repeat
			end if
		end repeat
	end tell
end alfred_script
