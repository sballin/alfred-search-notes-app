## Acknowledgement

The notestore.proto file is copied from the [apple_cloud_notes_parser](https://github.com/threeplanetssoftware/apple_cloud_notes_parser/tree/master/proto) project.

## Building

To get the protoc command, do

    brew install protobuf
    go install google.golang.org/protobuf/cmd/protoc-gen-go

Then in this folder, do

    protoc -I=. --go_out=. notestore.proto

## Development

To play around with parsing protobuf, see [protobuf-inspector](https://github.com/jmendeth/protobuf-inspector) and this [protobuf_config.py](https://github.com/threeplanetssoftware/apple_cloud_notes_parser/blob/master/proto/protobuf_config.py) file. Use DB Browser for SQLite.app to save some zipped note data to disk, rename the file as .gz and do gunzip, then in protobuf-inspector do

    ./main.py < file
