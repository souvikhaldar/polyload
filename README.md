# polyload
Payload is the server that allows to download/upload files to any cloud service provider of choice in the fastest and most convenient way.

# Instructios
1. Install golang
2. Provide configurations in `config.json`  
An example of `config.json` in `~/.polyload/`:  
```
{
        "port":"8192",
        "max_upload_size": 1073741824,
        "upload_dir": "/tmp/polyload",
        "download_dir":"~/polyload-download/",
        "registered_users":"$HOME/.polyload/registeredUsers.txt",
        "aws":{
                "s3_bucket_name":"polyload",
                "access_key_id":"",
                "secret_access_key":"",
                "s3_region":"ap-south-1"
        },
        "azure":{
                "object_id":"",
                "blob_container_name":"polyload",
                "blob_storage_account_name":"polyload",
                "storage_account_resource_id":""
        }
}

```  
In `registeredUsers.txt` file in this directory we provide the token of valid users in each line.  

3. Append random strings which can act as token for valid users [this needs to be improved to improve the security]
4. run the server by running `go run cmd/polyload/main.go` 

# Demonstration
Youtube link- https://youtu.be/ZDYmxGZ-MjI 
# Features
Please mail souvikhaldar32@gmail.com to get a valid token for testing.  
1. Upload a file to local storage
`curl -XPOST -F "file=@filename" "http://souvikhaldar.in/upload?cloud=local&token=valid_token"`  
2. Upload a file to azure storage
`curl -XPOST -F "file=@filename" "http://souvikhaldar.in/upload?cloud=azure&token=valid_token"`  
3. Upload a file to aws storage
`curl -XPOST -F "file=@filename" "http://souvikhaldar.in/upload?cloud=aws&token=valid_token"`  
4. Download the same file from browser
Visit `http://http://souvikhaldar.in/download?cloud=local&file=filename&token=valid_token`  
or  
Visit `http://http://souvikhaldar.in/download?cloud=azure&file=filename&token=valid_token`  
or  
Visit `http://http://souvikhaldar.in/download?cloud=aws&file=filename&token=valid_token`
