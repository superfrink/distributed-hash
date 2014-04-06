package main

// project: distributed hash
// file: hc.go - a client that talks to hash daemons
// purpose: CLI that talks to the nodes in the hash table
// git: FIXME

import "bufio"
//import "encoding/gob"
import "encoding/json"
import "flag"
import "fmt"
import "hash/crc32"
import "io/ioutil"
import "net"
import "os"
import "strings"

type HashServerConfig struct {
	Servers []string
	ServerCount int
}

type HashResponse struct {
	Status string
	Value  string
}

type HashCommand struct {
	Cmd   string
	Key   string
	Value string
	Out   chan HashResponse
}

var debug bool

func read_hash_config(out *HashServerConfig) error {

	config_string, err := ioutil.ReadFile("hc.conf")
	if nil != err {
		return err
	}

	err = json.Unmarshal(config_string, out)
	if nil != err {
		return err
	}

	out.ServerCount = len(out.Servers)

	return nil
}

func select_hash_server(config *HashServerConfig, key string) (server string) {

	h := crc32.NewIEEE()
    h.Write([]byte(key))
    i := h.Sum32()
    fmt.Printf("i: %v\n", i)

	i = i % uint32(config.ServerCount)
    fmt.Printf("i: %v\n", i)

	return config.Servers[i]
}

func hash_read(config *HashServerConfig, key string) (string, error){

	server := select_hash_server(config, key)
	fmt.Printf("selected server: %s\n", server)

	// GOAL : connect to the server and read the value
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return "", err
	}

	//enc := gob.NewEncoder(conn) // FIXME : incomplete

	fmt.Fprintf(conn, "GET %s\n", key)
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Printf("status: %s\n", status);

	value := ""
	if "NOTEXISTS" != status {
		// FIXME : incomplete
	}

	return value, nil
}

func hash_write(config *HashServerConfig, key string, value string) (error) {

	server := select_hash_server(config, key)
	fmt.Printf("selected server: %s\n", server)

	// GOAL : connect to the server and read the value
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return err
	}

	//enc := gob.NewEncoder(conn) // FIXME : incomplete

	fmt.Fprintf(conn, "PUT %s %s\n", key, value)
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("status: %s\n", status);

	return nil
}

func main() {

	// GOAL : process command line arguments

	flag.BoolVar(&debug, "d", false, "enable debug output")
	flag.Parse()

	fmt.Printf("debug: '%v'\n", debug)

	if flag.NArg() < 2 {
	   fmt.Println("FIXME - usage message")
	   os.Exit(1)
	}

	// GOAL : read the command from the user

	var cliRequest HashCommand
	cliRequest.Cmd = flag.Arg(0)
	cliRequest.Key = flag.Arg(1)

	if flag.NArg() >= 3 {
		cliRequest.Value = flag.Arg(2)
	}

	cliRequest.Cmd = strings.ToUpper(cliRequest.Cmd)

	fmt.Printf("cliRequest %+v\n", cliRequest)

	// GOAL : read the hash server config

	var hash_config HashServerConfig
	err := read_hash_config(&hash_config)
	if nil != err {
		fmt.Printf("error unable to read config %v\n", err);
	}
	fmt.Printf("hash_config: %+v\n", hash_config);

	// GOAL : execute the command from the user

	if "GET" == cliRequest.Cmd {
		
		value, err := hash_read(&hash_config, cliRequest.Key)
		if nil != err {
			fmt.Printf("hash_read failed. %v\n", err)
		}
		fmt.Printf("read value: %s\n", value);	

	} else if "PUT" == cliRequest.Cmd {

		err := hash_write(&hash_config, cliRequest.Key, cliRequest.Value)
		if nil != err {
			fmt.Printf("hash_write failed. %v\n", err)
		}
	}
}
