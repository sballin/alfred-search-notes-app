#!/usr/bin/osascript

use framework "Foundation"

property NSString : a reference to current application's NSString
property NSJSONSerialization : a reference to current application's NSJSONSerialization
property NSUTF8StringEncoding : a reference to current application's NSUTF8StringEncoding

on run argv
	set output to {}
	
	tell application "Notes"
		set allNotes to a reference to every note in default account
		set noteIDs to id of allNotes
		set noteNames to name of allNotes
		repeat with i from 1 to count of allNotes
			set end of output to {title:(item i of noteNames), arg:(item i of noteIDs), subtitle:"Notes.app"}
		end repeat
	end tell
	
	set output to {|items|:output}
	set output to NSJSONSerialization's dataWithJSONObject:output options:0 |error|:(missing value)
	set output to (NSString's alloc()'s initWithData:output encoding:NSUTF8StringEncoding) as text
	return output
end run