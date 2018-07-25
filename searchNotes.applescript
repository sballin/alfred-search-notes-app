#!/usr/bin/osascript
on run argv
	tell application "Notes"
		set noteRefs to a reference to every note in default account
	end tell
	set noteNames to name of noteRefs

	set output to "{\"items\":["
	repeat with name in noteNames
		set output to output & "{\"title\":\"" & name  & "\",\"arg\":\"" & name & "\"},"
	end repeat
	return output & "]}" 
end run
