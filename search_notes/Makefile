# -ldflags... reduces binary size by removing debug symbols
# -trimpath removes any references to absolute paths on the machine that did the compiling
# CGO_ENABLED=1 is required by go-sqlite3
FLAGS = -ldflags "-s -w" -trimpath

.PHONY: universal current clean
	
universal: search_notes_intel search_notes_apple
	# Create universal binary that works on both Intel and Apple CPUs
	lipo -create -output search_notes search_notes_intel search_notes_apple
	rm -f search_notes_intel search_notes_apple

search_notes_apple:
	# Build for Apple CPUs
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build $(FLAGS) -o search_notes_apple
	
search_notes_intel:
	# Build for Intel CPUs
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build $(FLAGS) -o search_notes_intel

current:
	# Quicker build only for the current machine's architecture
	go build $(FLAGS) -o search_notes

clean:
	rm -f search_notes search_notes_intel search_notes_apple
