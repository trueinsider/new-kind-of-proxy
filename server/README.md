## Build
Simply run:
```shell
glide install
go build server.go
```

## How to use
Edit `config.json` with your data:
```json
{
  "ListenPort": 8333,
  "DialTimeout": 30,
  "PrivateKey": "cd5fa29ed5b0e951f3d1bce5997458706186320f1dd89156a73d54ed752a7f37",
  "SubscriptionDuration": 65535
}
```
`ListenPort` port to listen for connections  
`DialTimeout` timeout for node-to-HTTPS connection  
`PrivateKey` your private key  
`SubscriptionDuration` duration for subscription in blocks  
`SubscriptionInterval` interval for subscription in seconds  

Run like this:
```shell
./server
```

Then users can connect to your NKN proxy through their New Kind of Proxy client
