#!/usr/bin/osascript

use framework "Foundation"

on run argv
	set output to {}
	
	tell application "Notes"
		set folderNames to name of folders in default account
		repeat with folderIndex from 1 to number of folders in default account
			set folderName to item folderIndex of folderNames
			if folderName is not "Recently Deleted" then
				set currentFolder to (a reference to item folderIndex of folders in default account)
				set noteIDs to id of notes of currentFolder
				set noteNames to name of notes of currentFolder
				repeat with i from 1 to count of noteIDs
					set end of output to {title:(item i of noteNames), arg:(item i of noteIDs), subtitle:folderName, icon:{|type|:"fileicon", path:"/Applications/Notes.app"}}
				end repeat
			end if
		end repeat
	end tell
	
	set output to {|items|:output}
	
	-- Source: https://forum.latenightsw.com/t/writing-json-data-with-nsjsonserialization
	set ca to current application
	set output to (ca's NSJSONSerialization)'s dataWithJSONObject:output options:0 |error|:(missing value)
	set output to (ca's NSString's alloc()'s initWithData:output encoding:(ca's NSUTF8StringEncoding)) as text
	return output
end run