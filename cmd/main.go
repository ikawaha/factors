package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp/syntax"

	"github.com/ikawaha/factors/factors"
)

const defaultAddr = ":6060"

var httpAddr = flag.String("http", "", "HTTP service address (e.g. '"+defaultAddr+"')")

// Usage prints a usage of this command.
func Usage() {
	fmt.Fprintln(os.Stderr, "command: factors <regexp_pattern>")
	fmt.Fprintln(os.Stderr, "server:  factors -http="+defaultAddr)
	flag.PrintDefaults()
}

// Run is the bootstrap function of a command.
func Run() error {
	flag.Usage = Usage
	flag.Parse()
	if len(os.Args) != 2 && *httpAddr == "" {
		Usage()
		os.Exit(1)
	}
	if *httpAddr != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/_demo", http.StatusFound)
		})
		http.HandleFunc("/_demo", demoHandler)
		log.Fatal(http.ListenAndServe(*httpAddr, nil))
	}
	re, err := syntax.Parse(os.Args[1], syntax.Perl)
	if err != nil {
		return err
	}
	a := factors.NewAnalyzer()
	f := a.Factor(re)

	fmt.Printf("Exact: %s\n", f.Exact)
	fmt.Printf("Prefix: %s\n", f.Prefix)
	fmt.Printf("Suffix: %s\n", f.Suffix)
	fmt.Printf("Fragment: %s\n", f.Fragment)

	return nil
}
