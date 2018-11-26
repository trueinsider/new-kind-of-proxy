package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/xtaci/smux"
)

var nodeConn net.Conn
var nodeSession *smux.Session
var config = &Configuration{}

type Configuration struct {
	SeedNode  string `json:"SeedNode"`
	Listener  string `json:"Listener"`
	PublicKey string `json:"PublicKey"`
}

type RPCResponse struct {
	Result  string `json:"result"`
}

func pipe(dest io.WriteCloser, src io.ReadCloser) {
	defer dest.Close()
	defer src.Close()
	io.Copy(dest, src)
}

func connectToNode(force bool) (net.Conn, error) {
	if nodeConn == nil || force {
		data := []byte(`{"jsonrpc":"2.0","method":"gethttpproxyaddr","params":{"address":"` + config.PublicKey + `"}}`)
		r := bytes.NewReader(data)
		resp, err := http.Post(config.SeedNode, "application/json", r)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		rpcResp := &RPCResponse{}
		err = json.Unmarshal(body, rpcResp)
		if err != nil {
			return nil, err
		}

		nodeConn, err = net.Dial("tcp", rpcResp.Result) //TODO: add timeout
		if err != nil {
			return nil, err
		}
	}

	return nodeConn, nil
}

func getSession(force bool) (*smux.Session, error) {
	if nodeSession == nil || nodeSession.IsClosed() || force {
		nodeConn, err := connectToNode(force)
		if err != nil {
			return nil, err
		}
		nodeSession, err = smux.Client(nodeConn, nil)
		if err != nil {
			if !force {
				return getSession(true)
			} else {
				return nil, err
			}
		}
	}

	return nodeSession, nil
}

func openStream(force bool) (*smux.Stream, error) {
	nodeSession, err := getSession(force)
	if err != nil {
		return nil, err
	}
	nodeStream, err := nodeSession.OpenStream()
	if err != nil {
		if !force {
			return openStream(true)
		} else {
			return nil, err
		}
	}
	return nodeStream, err
}

func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panicln("Couldn't read config:", err)
	}

	err = json.Unmarshal(file, config)
	if err != nil {
		log.Panicln("Couldn't unmarshal config:", err)
	}

	browserListener, err := net.Listen("tcp", config.Listener)
	if err != nil {
		log.Panicln("Couldn't bind listener:", err)
	}

	for {
		browserConn, err := browserListener.Accept()
		if err != nil {
			log.Println("Couldn't accept browser connection:", err)
			closeConnection(browserConn)
			continue
		}

		nodeStream, err := openStream(false)
		if err != nil {
			log.Println("Couldn't open stream:", err)
			closeConnection(browserConn)
			continue
		}

		go pipe(nodeStream, browserConn)
		go pipe(browserConn, nodeStream)
	}
}
