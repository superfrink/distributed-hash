package main

// project: distributed hash
// file: hd.go - a single hash daemon
// purpose: the daemon for a single node in a distributed hash table
// git: 

import "bytes"
import "flag"
import "fmt"
import "log"
import "net"
import "strings"

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

var debug bool
var listenPort uint

func handleClientRequest(conn net.Conn, hashAccessor chan HashRequest) {

	maxBufSize := 1024
	var err error

	var totalBuf bytes.Buffer
	for { // read and process a command loop

		requestString := ""

		for strings.Index(requestString, "\n") < 0 { // read until we find a single command loop

			if bytes.IndexByte(totalBuf.Bytes(), 0x0A) < 0 {

				buf := make([]byte, maxBufSize)
				byteCount, err := conn.Read(buf)
				if err != nil {
					log.Println(err)
					conn.Close()
					return
				}

				if byteCount >= maxBufSize {

					byteCount, err = conn.Write([]byte(fmt.Sprintf("Reuqests must be shorter than %d bytes.\n", maxBufSize)))
					if err != nil {
						log.Println(err)
						conn.Close()
						return
					}
					log.Println("request greater than buffer size")
					conn.Close()
					return
				}

				buf = bytes.Trim(buf, string(0))
				totalBuf.Write(buf)

			}
			//log.Printf("totalBuf: %v", totalBuf)
			if debug {
				log.Printf("Bytes(): %v", totalBuf.Bytes())
			}

			var cmd []byte
			// DOC: 0x0A is "\n"
			if bytes.IndexByte(totalBuf.Bytes(), 0x0A) > -1 {
				cmd, err = totalBuf.ReadBytes(0x0A)
				if err != nil {
					log.Println(err)
					conn.Close()
					return
				}
				if debug {
					log.Printf("cmd found: %s\n", cmd)
					log.Printf("Bytes(): %v", totalBuf.Bytes())
				}
				requestString = string(cmd)
				break
			}
		}

		requestString = strings.Trim(requestString, string(0)) // trim buffer null bytes
		requestString = strings.TrimSpace(requestString)       // trim trailing newlines, github issue #1

		// FIXME : replace the Index() with a split to get the command name and the rest

		// GOAL : get the command from the request string
		if 0 == strings.Index(requestString, "GET ") {
			//conn.Write([]byte("get\n"));

			parts := strings.SplitN(requestString, " ", 2)
			cmd := parts[0]
			key := parts[1]
			conn.Write([]byte(fmt.Sprintf("cmd: %s\n", cmd)))
			conn.Write([]byte(fmt.Sprintf("key: %s\n", key)))
			conn.Write([]byte(fmt.Sprintf("len: %v\n", len(key))))

			if 0 == len(key) {
				conn.Write([]byte("Missing key in GET\n"))
				log.Println("missing key in GET")
				conn.Close()
				return
			}

			responseChannel := make(chan HashResponse)

			var request HashRequest
			request.cmd = cmd
			request.key = key
			request.out = responseChannel
			hashAccessor <- request

			response := <-responseChannel
			fmt.Printf("response: '%+v'\n", response)
			if "EXISTS" == response.status {
				conn.Write([]byte(fmt.Sprintf("%s %d %s\n", response.status, len(response.value), response.value)))
			} else {
				conn.Write([]byte(fmt.Sprintf("%s\n", response.status)))
			}

		} else if 0 == strings.Index(requestString, "PUT ") {
			//conn.Write([]byte("put\n"));

			if strings.Count(requestString, " ") < 2 {
				conn.Write([]byte("Missing value in PUT\n"))
				log.Println("missing value in PUT")
				conn.Close()
				return
			}

			parts := strings.SplitN(requestString, " ", 3)
			cmd := parts[0]
			key := parts[1]
			val := parts[2]
			conn.Write([]byte(fmt.Sprintf("cmd: %s\n", cmd)))
			conn.Write([]byte(fmt.Sprintf("key: %s\n", key)))
			conn.Write([]byte(fmt.Sprintf("val: %s\n", val)))

			responseChannel := make(chan HashResponse)

			var request HashRequest
			request.cmd = cmd
			request.key = key
			request.value = val
			request.out = responseChannel
			hashAccessor <- request

			response := <-responseChannel
			fmt.Printf("response: '%+v'\n", response)

			// FIXME : incomplete
		}

	}
}

func createHashAccessor(table map[string]string) chan HashRequest {
	requestChannel := make(chan HashRequest)

	go func() {
		for {
			request := <-requestChannel
			fmt.Printf("hashAccessor - cmd %s\n", request.cmd)
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
				// FIXME : incomplete
			}
		}
	}()

	return requestChannel
}

func main() {

	// GOAL : process command line arguments

	flag.BoolVar(&debug, "d", false, "enable debug output")
	flag.UintVar(&listenPort, "p", 1742, "port to listen on")
	flag.Parse()

	fmt.Printf("debug: '%v'\n", debug)
	fmt.Printf("Listen port: '%d'\n", listenPort)

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
