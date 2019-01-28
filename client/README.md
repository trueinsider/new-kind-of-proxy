## Build
Simply run:
```shell
glide install
go build client.go
```

## How to use
Edit `config.json` with your data:
```json
{
  "Listener": ":8888",
  "NodeDialTimeout": 30,
  "PublicKey": ""
}
```
`Listener` port to listen for connections  
`NodeDialTimeout` timeout for NKN node connection  
`PrivateKey` your private key  

Run like this:
```shell
./client
```

Then you can set HTTPS proxy address in your browser (`127.0.0.1:8888` for example)
