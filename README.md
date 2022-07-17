# polyload
Payload is the server that allows to download/upload files to any cloud service provider of choice in the fastest and most convenient way.

# Features
1. Upload a file to local storage
`curl -XPOST -F "file=@filename" "ip:port/upload?cloud=local&token=valid_token"`
2. Upload a file to azure storage
`curl -XPOST -F "file=@filename" "ip:port/upload?cloud=azure&token=valid_token"`
3. Download the same file from browser
Visit `http://ip:port/download?cloud=local&file=filename`
or
Visit `http://ip:port/download?cloud=azure&file=filename`
