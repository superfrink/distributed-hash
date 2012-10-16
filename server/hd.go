package main
// project: distributed hash
// file: hd.go - a single hash daemon
// purpose: the daemon for a single node in a distributed hash table
// git: 

import "flag"
import "fmt"
import "log"
import "net"
import "strings"

type HashResponse struct {
	status string
	value string
}
type HashRequest struct {
	cmd string
	key string
	value string
	out chan HashResponse
}

var debug bool
var listenPort uint

func handleClientRequest(conn net.Conn, hashAccessor chan HashRequest) {

	maxBufSize := 1024  // FIXME: is this size good?

	for {
		buf := make([]byte, maxBufSize)

		byteCount, err := conn.Read(buf);
		if err != nil {
			log.Println(err)
			conn.Close();
			return
		}

		// FIXME : it would be better to append to a buffer for long requests
		if byteCount >= maxBufSize {

			byteCount, err = conn.Write([]byte(fmt.Sprintf("Reuqests must be shorter than %d bytes.\n", maxBufSize)));
			if err != nil {
				log.Println(err)
				conn.Close();
				return
			}
			log.Println("request greater than buffer size")
			conn.Close();
			return
		}

		// // echo back the number of bytes read and the message.
		// byteCount, err = conn.Write([]byte(fmt.Sprintf("%d\n",byteCount)));
		// if err != nil {
		// 	log.Println(err)
		// 	conn.Close();
		// 	return
		// }
		// byteCount, err = conn.Write(buf);
		// if err != nil {
		// 	log.Println(err)
		// 	conn.Close();
		// 	return
		// }

		requestString := string(buf)

		// FIXME : replace the Index() with a split to get the command name and the rest

		// GOAL : get the command from the request string
		if 0 == strings.Index(requestString, "GET ") {
			//conn.Write([]byte("get\n"));

			parts := strings.SplitN(requestString, " ", 2)
			cmd := parts[0]
			key := parts[1]
			key = strings.Trim(key, string(0)) // trim buffer null bytes
			conn.Write([]byte(fmt.Sprintf("cmd: %s\n", cmd)));
			conn.Write([]byte(fmt.Sprintf("key: %s\n", key)));
			conn.Write([]byte(fmt.Sprintf("len: %v\n", len(key))));

			if 0 == len(key) {
				conn.Write([]byte("Missing key in GET\n"));
				log.Println("missing key in GET")
				conn.Close();
				return
			}

			responseChannel := make(chan HashResponse)

			var request HashRequest
			request.cmd = cmd
			request.key = key
			request.out = responseChannel
			hashAccessor <- request

			response := <- responseChannel
			fmt.Printf("response: '%+v'\n", response)
			conn.Write([]byte(fmt.Sprintf("%s %d%s\n", response.status, len(response.value), response.value)));

		} else if 0 == strings.Index(requestString, "PUT ") {
			//conn.Write([]byte("put\n"));

			if strings.Count(requestString, " ") < 2 {
				conn.Write([]byte("Missing value in PUT\n"));
				log.Println("missing value in PUT")
				conn.Close();
				return
			}

			parts := strings.SplitN(requestString, " ", 3)
			cmd := parts[0]
			key := parts[1]
			val := parts[2]
			val = strings.Trim(val, string(0)) // trim buffer null bytes
			conn.Write([]byte(fmt.Sprintf("cmd: %s\n", cmd)));
			conn.Write([]byte(fmt.Sprintf("key: %s\n", key)));
			conn.Write([]byte(fmt.Sprintf("val: %s\n", val)));

			responseChannel := make(chan HashResponse)

			var request HashRequest
			request.cmd = cmd
			request.key = key
			request.value = val
			request.out = responseChannel
			hashAccessor <- request

			response := <- responseChannel
			fmt.Printf("response: '%+v'\n", response)

			// FIXME : incomplete
		}

	}
}

func createHashAccessor(table map[string]string) chan HashRequest {
	requestChannel := make(chan HashRequest)

	go func() {
		for {
			request := <- requestChannel
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

	return requestChannel;
}

func main() {

	// GOAL : process command line arguments

	flag.BoolVar(&debug, "d", false, "enable debug output")
	flag.UintVar(&listenPort, "p", 1742, "port to listen on")
	flag.Parse()

	fmt.Printf("debug: '%v'\n", debug);
	fmt.Printf("Listen port: '%d'\n", listenPort);

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
