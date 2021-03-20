make: 
	# -ldflags... reduce binary size by removing debug symbols
	# -trimpath remove any references to absolute paths on the machine that did the compiling
	go build \
	-ldflags "-s -w" \
	-trimpath \
	-o search_notes
