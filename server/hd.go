package main

// project: distributed hash
// file: hd.go - a single hash daemon
// purpose: the daemon for a single node in a distributed hash table
// git: https://github.com/superfrink/distributed-hash.git

import "encoding/gob"
import "flag"
import "fmt"
import "io"
import "log"
import "net"
import "os"

type HashResponse struct {
	status string
	value  string
}
type HashRequest struct {
	cmd   string
	key   string
	value string
	out   chan HashResponse
}

type HashWireMessage struct {
	Cmd    string
	Key    string
	Value  string
	Status string
}

var debug bool
var listenPort uint

func handleClientRequest(conn net.Conn, hashAccessor chan HashRequest) {

	var err error

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	var wireRequest HashWireMessage
	var wireResponse HashWireMessage

	for { // read and process a command loop

		wireRequest = HashWireMessage{}
		wireResponse = HashWireMessage{}

		err = dec.Decode(&wireRequest)
		if io.EOF == err {
			return
		}
		if err != nil {
			fmt.Printf("decode error: ", err)
			os.Exit(1) // FIXME : don't exit
		}
		if debug {
			fmt.Printf("wireRequest: %+v\n", wireRequest)
		}

		if "GET" == wireRequest.Cmd {
			responseChannel := make(chan HashResponse)

			var request HashRequest
			request.cmd = wireRequest.Cmd
			request.key = wireRequest.Key
			request.out = responseChannel
			hashAccessor <- request

			response := <-responseChannel
			if debug {
				fmt.Printf("response: '%+v'\n", response)
			}

			wireResponse.Cmd = request.cmd
			wireResponse.Key = request.key

			if "EXISTS" == response.status {
				wireResponse.Value = response.value
				wireResponse.Status = "EXISTS"

			} else {
				wireResponse.Status = "NOEXISTS"
			}

		} else if "PUT" == wireRequest.Cmd {

			responseChannel := make(chan HashResponse)

			var request HashRequest
			request.cmd = wireRequest.Cmd
			request.key = wireRequest.Key
			request.value = wireRequest.Value
			request.out = responseChannel
			hashAccessor <- request

			response := <-responseChannel
			if debug {
				fmt.Printf("response: '%+v'\n", response)
			}

			wireResponse.Cmd = request.cmd
			wireResponse.Key = request.key
			wireResponse.Status = "ERROR"
			if "OK" == response.status {
				wireResponse.Status = "OK"
			}

		} else {
			// wireRequest.Cmd is not valid
			wireResponse.Cmd = wireRequest.Cmd
			wireResponse.Status = "INVALIDCMD"
		}

		if debug {
			fmt.Printf("wireResponse: %+v\n", wireResponse)
		}
		err = enc.Encode(wireResponse)
		if err != nil {
			fmt.Printf("encode error:", err)
			os.Exit(1) // FIXME don't exit
		}
	}
}

func createHashAccessor(table map[string]string) chan HashRequest {
	requestChannel := make(chan HashRequest)

	go func() {
		for {
			request := <-requestChannel
			if debug {
				fmt.Printf("hashAccessor - cmd %s\n", request.cmd)
			}
			switch request.cmd {
			case "GET":
				var response HashResponse
				var ok bool

				if response.value, ok = table[request.key]; ok {
					response.status = "EXISTS"
				} else {
					response.status = "NOTEXISTS"
				}

				go func() { request.out <- response }()

			case "PUT":
				table[request.key] = request.value
				var response HashResponse
				response.status = "OK"
				go func() { request.out <- response }()

			default:
				fmt.Printf("hashAccessor - unexpected cmd %s\n", request.cmd)
				os.Exit(1) // FIXME don't exit
			}
		}
	}()

	return requestChannel
}

func usage_message() string {
	str := `
Usa: hd [OPTIONS]

  OPTIONS are:
    -d           Enable debugging output.
    -h           Show this usage message.
    -p           Port number to listen on.  Defaults to 1742.

Examples:

  hd
  hd -p 1234

`
	return str
}

func main() {

	// GOAL : process command line arguments

	var help bool
	flag.BoolVar(&debug, "d", false, "enable debug output")
	flag.BoolVar(&help, "h", false, "show usage message")
	flag.UintVar(&listenPort, "p", 1742, "port to listen on")
	flag.Parse()

	if help {
		fmt.Print(usage_message())
		os.Exit(1)
	}

	fmt.Printf("Listening on port %d.\n", listenPort)

	// GOAL : create the storage for the hash table

	hashTable := make(map[string]string)
	hashAccessor := createHashAccessor(hashTable)

	// GOAL : accept connections from the network

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		log.Println("Error on listen()\n")
		log.Fatal(err)
		// FIXME handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error on accept()\n")
			log.Println(err)
			// FIXME handle error
			continue
		}
		go handleClientRequest(conn, hashAccessor)
	}
}
