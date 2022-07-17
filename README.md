# polyload
Payload is the server that allows to download/upload files to any cloud service provider of choice in the fastest and most convenient way.

# Instructios
1. Install golang
2. Provide configurations in `config.json`
3. Append random strings which can act as token for valid users [this needs to be improved to improve the security]
4. run the server by running `go run cmd/polyload/main.go` 

# Demonstration
Youtube link- https://www.youtube.com/watch?v=ESacD9lM8zA

# Features
1. Upload a file to local storage
`curl -XPOST -F "file=@filename" "ip:port/upload?cloud=local&token=valid_token"`
2. Upload a file to azure storage
`curl -XPOST -F "file=@filename" "ip:port/upload?cloud=azure&token=valid_token"`
3. Download the same file from browser
Visit `http://ip:port/download?cloud=local&file=filename`
or
Visit `http://ip:port/download?cloud=azure&file=filename`
