// Runs a HTTP server in the current directory.
// David Leadbeater, 2014. http://dgl.cx/licence (WTFPL).
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	quiet *bool
	port  *string
	safe  *bool
)

func init() {
	var default_port string
	if os.Getuid() == 0 {
		default_port = ":http"
	} else {
		default_port = ":8080" // http-alt, but not all systems know this
	}
	quiet = flag.Bool("quiet", false, "Don't log except for errors")
	port = flag.String("port", default_port, "[ip]:port to listen on")
	safe = flag.Bool("safe", true, "Only serve world readable files. "+
		"This is merely a safety measure to avoid accidents, do not rely on it.")
	flag.Parse()
}

func main() {
	http.Handle("/", http.FileServer(dirWrapper("")))
	if !*quiet {
		log.Print("Opening server on ", *port, " available at:")
		printAddrs()
	}
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func printAddrs() {
	ip, port, err := net.SplitHostPort(*port)
	if err != nil {
		log.Print("Unable to split host/port: ", err)
		return
	}
	if ip == "" {
		interfaces, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range interfaces {
				printAddr(addr.(*net.IPNet).IP.String(), port)
			}
		} else {
			log.Print("Unable to query interfaces: ", err)
		}
	} else {
		printAddr(ip, port)
	}
}

func printAddr(addr, port string) {
	if port == "http" {
		port = "80"
	}
	fmt.Printf("  - http://%v/\n", net.JoinHostPort(addr, port))
}

// Wrap a http.Dir with a very basic permission check. Beware, this only checks
// one level (not if it can recurse into the directory it is inside) this is
// designed to avoid accidents like serving your ssh private key or something.
type dirWrapper string
func (d dirWrapper) Open(path string) (http.File, error) {
	if !*quiet {
		log.Print(path)
	}
	file := path
	// Canonicalise, for some reason we get a request for "//index.html"
	for len(file) > 0 && file[0] == '/' {
		file = file[1:]
	}
	if file == "" {
		file = "."
	}
	fi, err := os.Stat(file)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if *safe && fi.Mode().Perm()&4 == 0 {
		log.Print(d, path, " is not world readable")
		return nil, errors.New("Not world readable")
	}
	return http.Dir(d).Open(path)
}
