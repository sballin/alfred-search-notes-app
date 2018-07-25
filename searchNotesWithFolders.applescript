#!/usr/bin/osascript
on run argv
	tell application "Notes"
		set noteRefs to a reference to every note in default account
		set noteNames to name of noteRefs
		set containers to container of noteRefs
	end tell

	set output to "{\"items\":["
	repeat with i from 1 to count of noteRefs
		set output to output & "{\"title\":\"" & (item i in noteNames)  & "\",\"arg\":\"" & (item i in noteNames) & "\",\"subtitle\": \"" & (name of item i of containers) & "\"},"
	end repeat
	return output & "]}" 
end run
