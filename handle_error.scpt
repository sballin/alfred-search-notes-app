#!/usr/bin/osascript
set alertTitle to "Alfred Search Notes workflow error."
set alertMessage to "This error may be due to the \"search\" binary not being authorized to run. To manually authorize it, right click the binary and click \"Open\"."
display alert alertTitle message alertMessage as critical buttons {"Cancel", "Open Error Log", "Show Binary"} default button "Show Binary" cancel button "Cancel"
if button returned of result = "Show Binary" then
    do shell script "open ./search"
else if button returned of result = "Open Error Log" then
    do shell script "open ./error_log.txt"
end if
