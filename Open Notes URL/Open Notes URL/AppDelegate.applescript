script AppDelegate
    property parent : class "NSObject"
    
    -- IBOutlets
    property theWindow : missing value
    
    use framework "Foundation"
    use scripting additions
    property NSString : a reference to current application's NSString
    property NSCharacterSet : a reference to current application's NSCharacterSet
    property NSDate : a reference to current application's NSDate
    
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
                set timestampArg to characters 45 thru -1 of noteURL as text
                set stringTimestamp to (NSString's stringWithString:timestampArg)
                set doubleTimestamp to stringTimestamp's doubleValue
                set creationDate to (NSDate's dateWithTimeIntervalSince1970:doubleTimestamp) as date
                show first note in default account whose creation date ≥ creationDate and creation date < (creationDate + 1)
            on error errorMessage number errorNumber
                set alertMessage to errorMessage & " (" & errorNumber & ")"
                display alert "Open Notes URL error" message alertMessage as critical
            end try
        end tell
        quit
    end handleGetURLEvent_
    
end script
