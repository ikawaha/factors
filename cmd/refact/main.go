package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp/syntax"

	"github.com/ikawaha/regexp/factor"
)

const (
	defaultAddr = ":6060"
)

var (
	httpAddr = flag.String("http", "", "HTTP service address (e.g. '"+defaultAddr+"')")
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"usage: refact regexp\n"+
			"	refact -http="+defaultAddr+"\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) != 2 && *httpAddr == "" {
		usage()
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
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	f := factor.Analyze(re)

	printFactor("Exact", f.Exact)
	printFactor("Prefix", f.Prefix)
	printFactor("Suffix", f.Suffix)
	printFactor("Fact", f.Fact)
}

func printFactor(label string, s []string) {
	if s == nil {
		fmt.Printf("%s: <undef>\n", label)
		return
	}
	fmt.Printf("%s: %v\n", label, s)
}
