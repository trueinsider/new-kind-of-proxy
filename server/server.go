package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	. "github.com/nknorg/nkn-sdk-go"
	"github.com/nknorg/nkn/vault"
	"github.com/rdegges/go-ipify"
	"github.com/xtaci/smux"
)

var config = &Configuration{}

type Configuration struct {
	Hostname             string `json:"Hostname"`
	ListenPort           int    `json:"ListenPort"`
	DialTimeout          uint16 `json:"DialTimeout"`
	PrivateKey           string `json:"PrivateKey"`
	SubscriptionDuration uint32 `json:"SubscriptionDuration"`
}

func pipe(dest io.WriteCloser, src io.ReadCloser) {
	defer dest.Close()
	defer src.Close()
	io.Copy(dest, src)
}

type HTTPProxy struct {
	// defines a listeners for http proxy, such as "127.0.0.1:30004"
	Listener string
	timeout  time.Duration
}

func NewServer() *HTTPProxy {
	return &HTTPProxy{
		Listener: config.Hostname + ":" + strconv.Itoa(config.ListenPort),
		timeout:  time.Duration(config.DialTimeout),
	}
}

// parseRequestLine parses "GET /foo HTTP/1.1" into its three parts.
func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println("Error while closing connection:", err)
	}
}

func (s *HTTPProxy) handleSession(conn net.Conn, session *smux.Session) {
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Println("Couldn't accept stream:", err)
			break
		}

		tp := textproto.NewReader(bufio.NewReader(stream))

		line, err := tp.ReadLine()
		if err != nil {
			log.Println("Couldn't read line:", err)
			closeConnection(stream)
			continue
		}

		method, host, _, _ := parseRequestLine(line)

		// won't proxy HTTP due to security reasons
		if method != http.MethodConnect {
			log.Println("Only CONNECT HTTP method supported")
			stream.Write([]byte("HTTP/1.1 403 Forbidden\r\n\r\n"))
			closeConnection(stream)
			continue
		}

		destConn, err := net.DialTimeout("tcp", host, s.timeout*time.Second)
		if err != nil {
			log.Println("Couldn't connect to host", host, "with error:", err)
			closeConnection(stream)
			continue
		}

		stream.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

		go pipe(destConn, stream)
		go pipe(stream, destConn)
	}

	session.Close()
	conn.Close()
}

func (s *HTTPProxy) Start() {
	listener, err := net.Listen("tcp", s.Listener)
	if err != nil {
		log.Println("Couldn't bind HTTP proxy port:", err)
	}

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Couldn't accept client connection:", err)
			closeConnection(clientConn)
			continue
		}

		clientSession, err := smux.Server(clientConn, nil)
		if err != nil {
			log.Println("Couldn't create smux session:", err)
			closeConnection(clientConn)
			continue
		}

		go s.handleSession(clientConn, clientSession)
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

	if config.Hostname == "" {
		ip, err := ipify.GetIp()
		if err != nil {
			log.Panicln("Couldn't get IP:", err)
		}

		config.Hostname = ip
	}

	Init()

	privateKey, _ := hex.DecodeString(config.PrivateKey)
	account, err := vault.NewAccountWithPrivatekey(privateKey)
	if err != nil {
		log.Panicln("Couldn't load account:", err)
	}

	w := NewWalletSDK(account)

	s := NewServer()

	// retry subscription once a minute (regardless of result)
	go func() {
		for {
			txid, err := w.SubscribeToFirstAvailableBucket("", "proxyhttp", config.SubscriptionDuration, s.Listener)
			if err != nil {
				log.Println("Couldn't subscribe:", err)
			} else {
				log.Println("Subscribed to topic successfully:", txid)
			}

			time.Sleep(time.Minute)
		}
	}()

	s.Start()
}
