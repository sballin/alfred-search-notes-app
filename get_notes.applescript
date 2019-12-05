#!/usr/bin/osascript

use framework "Foundation"

property NSRegularExpressionSearch : a reference to 1024
property NSString : a reference to current application's NSString
property NSJSONSerialization : a reference to current application's NSJSONSerialization
property NSUTF8StringEncoding : a reference to current application's NSUTF8StringEncoding

on run argv
	set output to {}
	tell application "Notes"
		set folderNames to name of folders in default account
		repeat with folderIndex from 1 to number of folders in default account
			set folderName to item folderIndex of folderNames
			if folderName is not "Recently Deleted" then
				set currentFolder to (a reference to item folderIndex of folders in default account)
				set {noteIDs, noteNames, noteBodies} to {id, name, plaintext} of notes of currentFolder
				repeat with i from 1 to count of noteIDs
					set noteBody to item i of noteBodies
					set match to (item i of noteNames) & " " & folderName & " " & noteBody
					if length of noteBody is less than 100 then
						set subtitle to noteBody
					else
						set subtitle to text 1 thru 100 of noteBody
					end if
					set subtitle to folderName & "  |  " & subtitle
					set end of output to {title:(item i of noteNames), arg:(item i of noteIDs), subtitle:subtitle, match:match, uid:(item i of noteNames)}
				end repeat
			end if
		end repeat
	end tell
	
	set output to {|items|:output}
	set output to NSJSONSerialization's dataWithJSONObject:output options:0 |error|:(missing value)
	set output to (NSString's alloc()'s initWithData:output encoding:NSUTF8StringEncoding) as text
	return output
end run
