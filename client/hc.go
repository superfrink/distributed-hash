package main

// project: distributed hash
// file: hc.go - a client that talks to hash daemons
// purpose: CLI that talks to the nodes in the hash table
// git: FIXME

// import "bytes"
import "encoding/json"
import "flag"
import "fmt"
import "io/ioutil"
// import "net"
import "os"
import "strings"

type HashServerConfig struct {
	Servers []string
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

	return nil
}

func select_hash_server() {
	// FIXME : incomplete
}

func hash_read() {
	// FIXME : incomplete
}

func hash_write() {
	// FIXME : incomplete
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

	// FIXME : read value if NArg == 3

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
		
		// FIXME : incomplete; call hash_read()
	} else if "PUT" == cliRequest.Cmd {

		// FIXME : incomplete; call hash_write()
	}
}
