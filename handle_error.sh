osascript -e "display dialog \"The Search Notes workflow encountered an error. Press OK to open the error log and the Issues page on Github. Please include the error log output in your bug report.

The workflow requires python 3, which can be installed with this Terminal command:

xcode-select --install\""
if [[ $? == 0 ]]; then
    open "https://github.com/sballin/alfred-search-notes-app/issues"
    open error_log.txt
fi
