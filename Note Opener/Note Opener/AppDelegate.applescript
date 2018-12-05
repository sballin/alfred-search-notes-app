--
--  AppDelegate.applescript
--  Note Opener
--
--  Created by Sean on 12/5/18.
--  Copyright © 2018 Sean Ballinger. All rights reserved.
--

script AppDelegate
    property parent : class "NSObject"
    
    -- IBOutlets
    property theWindow : missing value
    
    use scripting additions
    
    use framework "Foundation"
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
        set noteURL to (ev's paramDescriptorForKeyword_(7.57935405E+8)) as string
        
        set noteName to (NSString's stringWithString:noteURL)
        set noteName to (noteName's stringByRemovingPercentEncoding) as text
        set noteName to text 9 thru (count of noteName) of noteName
        
        tell application "Notes"
            activate
            
            set noteRefs to a reference to every note in default account
            set noteNames to name of noteRefs
            
            set noteFound to false
            repeat with i from 1 to count of noteRefs
                if item i of noteNames contains noteName
                    set noteFound to true
                    show item i of noteRefs
                    exit repeat
                end if
            end repeat
        
            if noteFound is false
                display dialog "No note with title \"" & noteName & "\" was found."
            end if
        end tell
        quit
    end handleGetURLEvent_
    
end script

