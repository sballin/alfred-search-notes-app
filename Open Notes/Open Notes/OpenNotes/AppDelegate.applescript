--
--  AppDelegate.applescript
--  openNotes
--
--  Created by Sean on 12/1/18.
--  Copyright © 2018 Sean Ballinger. All rights reserved.
--

script AppDelegate
	property parent : class "NSObject"
	
	-- IBOutlets
	--property theWindow : missing value
	
	on applicationWillFinishLaunching_(aNotification)
        -- Insert code here to initialize your application before any files are opened
        -- Register the URL Handler stuff
        tell current application's NSAppleEventManager's sharedAppleEventManager() to setEventHandler_andSelector_forEventClass_andEventID_(me, "handleGetURLEvent:", current application's kInternetEventClass, current application's kAEGetURL)
    end applicationWillFinishLaunching_
    
	on applicationShouldTerminate_(sender)
		-- Insert code here to do any housekeeping before your application quits 
		return current application's NSTerminateNow
	end applicationShouldTerminate_
    
    -- handler that runs when the URL is clicked
    on handleGetURLEvent_(ev)
        set shortURL to (ev's paramDescriptorForKeyword_(7.57935405E+8)) as string
        tell application "Notes"
            set accountURL to id of default account
            set noteURL to (get characters 1 thru 49 of accountURL) as string
            set noteNumber to (get characters 9 thru (length of shortURL) of shortURL) as string
            set noteURL to noteURL & "/ICNote/" & noteNumber
            show note id noteURL
        end tell
        quit
    end handleGetURLEvent_
	
end script
