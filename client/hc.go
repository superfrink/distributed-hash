package main

// project: distributed hash
// file: hc.go - a client that talks to hash daemons
// purpose: CLI that talks to the nodes in the hash table
// git: FIXME

// import "bytes"
import "flag"
import "fmt"
// import "net"
import "os"
// import "strings"


type HashResponse struct {
	status string
	value  string
}

type HashCommand struct {
	cmd   string
	key   string
	value string
	out   chan HashResponse
}

var debug bool

func main() {

	// GOAL : process command line arguments

	flag.BoolVar(&debug, "d", false, "enable debug output")
	flag.Parse()

	fmt.Printf("debug: '%v'\n", debug)

	if flag.NArg() < 2 {
	   fmt.Println("FIXME - usage message")
	   os.Exit(1)
	}

	var cliRequest HashCommand
	cliRequest.cmd = flag.Arg(0)
	cliRequest.key = flag.Arg(1)

	// FIXME : read value if NArg == 3

	fmt.Printf("%+v\n", cliRequest)

	// FIXME : incomplete
}
