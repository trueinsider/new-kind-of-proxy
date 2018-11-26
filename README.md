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
  "SeedNode": "http://testnet-seed-0001.nkn.org:30003",
  "Listener": ":8888",
  "PublicKey": ""
}
```
`SeedNode` seed node to connect to NKN  
`Listener` port to listen for connections  
`PublicKey` your public key

Run like this:
```shell
./proxy
```

Then you can set HTTPS proxy address in your browser (`127.0.0.1:8888` for example)
