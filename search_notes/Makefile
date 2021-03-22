make: 
	# -ldflags... reduces binary size by removing debug symbols
	# -trimpath removes any references to absolute paths on the machine that did the compiling
	go build \
	-ldflags "-s -w" \
	-trimpath \
	-o search_notes
