make: 
	# -s and -w reduce binary size by removing debug symbols
	go build \
	-ldflags "-s -w" \
	-gcflags=-trimpath="$(shell pwd)" \
	-o search_notes
