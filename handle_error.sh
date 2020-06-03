osascript -e "display dialog \"The Search Notes workflow encountered an error. 

You may need to right click the 'search' binary and click Open for macOS to trust it. Press OK to open the folder containing the binary.\""
if [[ $? == 0 ]]; then
    open ./search
fi
