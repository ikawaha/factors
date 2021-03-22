package cmd

import (
	"bytes"
	_ "embed" //nolint:golint
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp/syntax"
	"strings"
	"time"

	"github.com/ikawaha/factors/factors"
)

const cmdTimeout = 30 * time.Second

var (
	//go:embed demo.html
	demoHTML     string
	demoTemplate = template.Must(template.New("top").Parse(demoHTML))
)

func demoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		svg, cmdErr string
		f           factors.Factor
	)
	if _, e := exec.LookPath("dot"); e != nil {
		log.Print("graphviz is not in your future\n")
		return
	}
	var buf bytes.Buffer
	cmd := exec.Command("dot", "-Tsvg")
	r0, w0 := io.Pipe()
	cmd.Stdin = r0
	cmd.Stdout = &buf
	if err := cmd.Start(); err != nil {
		cmdErr = err.Error()
		log.Printf("process error, %v", err)
	}
	input := r.FormValue("re")
	log.Println("input:", input)
	if len(input) != 0 {
		re, err := syntax.Parse(input, syntax.Perl)
		log.Println("re:", re)
		if err != nil {
			cmdErr = err.Error()
			log.Printf("input error, %v", err)
			goto END
		}
		a := factors.NewAnalyzer()
		f = a.DebugParse(w0, re)
		_ = w0.Close()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()
		select {
		case <-time.After(cmdTimeout):
			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill: ", err)
			}
			cmdErr = "Time out"
			<-done
		case err := <-done:
			if err != nil {
				cmdErr = "Error"
				log.Printf("process done with error, %v", err)
			}
		}
		svg = buf.String()
		if pos := strings.Index(svg, "<svg"); pos > 0 {
			svg = svg[pos:]
		}
	}
END:
	d := struct {
		Regexp                          string
		Exact, Prefix, Suffix, Fragment string
		CmdErr                          string
		GraphSvg                        template.HTML
	}{
		Regexp:   input,
		Exact:    strings.Join(f.Exact.Items(), ","),
		Prefix:   strings.Join(f.Prefix.Items(), ","),
		Suffix:   strings.Join(f.Suffix.Items(), ","),
		Fragment: strings.Join(f.Fragment.Items(), ","),
		CmdErr:   cmdErr,
		GraphSvg: template.HTML(svg),
	}
	if e := demoTemplate.Execute(w, d); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
