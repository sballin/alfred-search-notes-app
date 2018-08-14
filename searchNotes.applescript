#!/usr/bin/osascript

use framework "Foundation"

on run argv
	set output to {}
	
	tell application "Notes"
		set allNotes to a reference to every note in default account
		set noteIDs to id of allNotes
		set noteNames to name of allNotes
		repeat with i from 1 to count of allNotes
			set end of output to {title:(item i of noteNames), arg:(item i of noteIDs), icon:{|type|:"fileicon", path:"/Applications/Notes.app"}}
		end repeat
	end tell
	
	set output to {|items|:output}
	
	-- Source: https://forum.latenightsw.com/t/writing-json-data-with-nsjsonserialization
	set ca to current application
	set output to (ca's NSJSONSerialization)'s dataWithJSONObject:output options:0 |error|:(missing value)
	set output to (ca's NSString's alloc()'s initWithData:output encoding:(ca's NSUTF8StringEncoding)) as text
	return output
end run