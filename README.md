# polyload
Payload is the server that allows to download/upload files to any cloud service provider of choice in the fastest and most convenient way.

# Instructios
1. Install golang
2. Provide configurations in `config.json`
3. Append random strings which can act as token for valid users [this needs to be improved to improve the security]
4. run the server by running `go run cmd/polyload/main.go` 

# Demonstration
Youtube link- https://youtu.be/ZDYmxGZ-MjI 
# Features
Please mail souvikhaldar32@gmail.com to get a valid token for testing.  
1. Upload a file to local storage
`curl -XPOST -F "file=@filename" "http://souvikhaldar.in/upload?cloud=local&token=valid_token"`
2. Upload a file to azure storage
3. Upload a file to aws storage
`curl -XPOST -F "file=@filename" "http://souvikhaldar.in/upload?cloud=azure&token=valid_token"`
`curl -XPOST -F "file=@filename" "http://souvikhaldar.in/upload?cloud=aws&token=valid_token"`
3. Download the same file from browser
Visit `http://http://souvikhaldar.in/download?cloud=local&file=filename`  
or  
Visit `http://http://souvikhaldar.in/download?cloud=azure&file=filename`  
or  
Visit `http://http://souvikhaldar.in/download?cloud=aws&file=filename`
