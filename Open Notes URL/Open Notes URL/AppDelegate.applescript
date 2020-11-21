--
--  AppDelegate.applescript
--  Open Notes URL
--
--  Created by Sean on 12/5/18.
--  Copyright © 2018 Sean Ballinger. All rights reserved.
--

script AppDelegate
    property parent : class "NSObject"
    
    -- IBOutlets
    property theWindow : missing value
    
    use framework "Foundation"
    use scripting additions
    property NSString : a reference to current application's NSString
    property NSCharacterSet : a reference to current application's NSCharacterSet
    
    -- Insert code here to initialize your application before any files are opened
    on applicationWillFinishLaunching_(aNotification)
        -- Register the URL Handler stuff
        tell current application's NSAppleEventManager's sharedAppleEventManager() to setEventHandler_andSelector_forEventClass_andEventID_(me, "handleGetURLEvent:", current application's kInternetEventClass, current application's kAEGetURL)
    end applicationWillFinishLaunching_
    
    -- Insert code here to do any housekeeping before your application quits
    on applicationShouldTerminate_(sender)
        return current application's NSTerminateNow
    end applicationShouldTerminate_
    
    -- Handler that runs when the URL is clicked
    on handleGetURLEvent_(ev)
        tell application "Notes"
            activate
            try
                set noteURL to (ev's paramDescriptorForKeyword_(7.57935405E+8)) as string
                set noteName to (NSString's stringWithString:noteURL)
                set noteName to (noteName's stringByRemovingPercentEncoding) as text
                set noteName to text 8 thru (count of noteName) of noteName
                show (first note in default account whose name is noteName)
            on error errorMessage number errorNumber
                set alertMessage to errorMessage & " (" & errorNumber & ")"
                display alert "Open Notes URL error" message alertMessage as critical
            end try
        end tell
        quit
    end handleGetURLEvent_
    
end script
