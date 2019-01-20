# new-kind-of-proxy
Allows to directly connect to NKN node and browse HTTPS websites through it.  
Node will receive payment for proxied traffic directly from user (not implemented yet).

**Note:** HTTP (non-secure) won't be proxied because of security reasons

## Build
Simply run:
```shell
glide install
go build proxy.go
```

## How to use
Edit `config.json` with your data:
```json
{
  "SeedList": [
    "http://testnet-seed-0001.nkn.org:30003",
    "http://testnet-seed-0002.nkn.org:30003",
    "http://testnet-seed-0003.nkn.org:30003",
    "http://testnet-seed-0004.nkn.org:30003",
    "http://testnet-seed-0005.nkn.org:30003",
    "http://testnet-seed-0006.nkn.org:30003",
    "http://testnet-seed-0007.nkn.org:30003",
    "http://testnet-seed-0008.nkn.org:30003"
  ],
  "Listener": ":8888",
  "NodeDialTimeout": 30,
  "PublicKey": ""
}
```
`SeedList` list of seed nodes to connect to NKN  
`Listener` port to listen for connections  
`NodeDialTimeout` timeout for NKN node connection  
`PublicKey` your public key

Run like this:
```shell
./proxy
```

Then you can set HTTPS proxy address in your browser (`127.0.0.1:8888` for example)
