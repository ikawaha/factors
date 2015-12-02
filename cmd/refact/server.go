package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp/syntax"
	"strings"
	"time"

	"github.com/ikawaha/regexp/factor"
)

const (
	cmdTimeout = 30 * time.Second
)

var (
	demoTemplate = template.Must(template.New("top").Parse(demoHTML))
)

func demoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		svg, cmdErr string
		f           factor.Factor
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
			goto L_END
		}
		f = factor.DebugParse(w0, re)
		w0.Close()

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
L_END:
	d := struct {
		Regexp                      string
		Exact, Prefix, Suffix, Fact string
		CmdErr                      string
		GraphSvg                    template.HTML
	}{
		Regexp:   input,
		Exact:    strings.Join(f.Exact, ","),
		Prefix:   strings.Join(f.Prefix, ","),
		Suffix:   strings.Join(f.Suffix, ","),
		Fact:     strings.Join(f.Fact, ","),
		CmdErr:   cmdErr,
		GraphSvg: template.HTML(svg),
	}
	if e := demoTemplate.Execute(w, d); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

var demoHTML = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <style type="text/css">
      body {
        text-align: center;
      }
      div#center{
        width: 800px;
        margin: 0 auto;
        text-align: left;
      }
      .tbl{
        width: 100%;
        border-collapse: separate;
      }
      .tbl th{
        width: 20%;
        padding: 6px;
        text-align: left;
        vertical-align: top;
        color: #333;
        background-color: #eee;
        border: 1px solid #b9b9b9;
      }
      .tbl td{
        padding: 6px;
        background-color: #fff;
        border: 1px solid #b9b9b9;
      }
      .frm {
        min-height: 10px;
        padding: 0 10px 0;
        margin-bottom: 20px;
        background-color: #f5f5f5;
        border: 1px solid #e3e3e3;
        -webkit-border-radius: 4px;
        -moz-border-radius: 4px;
        border-radius: 4px;
        -webkit-box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
        -moz-box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
        box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
      }
      .txar {
         border:10px;
         padding:10px;
         font-size:1.1em;
         font-family:Arial, sans-serif;
         border:solid 1px #ccc;
         margin:0;
         width:80%;
         -webkit-border-radius: 3px;
         -moz-border-radius: 3px;
         border-radius: 3px;
         -moz-box-shadow: inset 0 0 4px rgba(0,0,0,0.2);
         -webkit-box-shadow: inset 0 0 4px rgba(0, 0, 0, 0.2);
         box-shadow: inner 0 0 4px rgba(0, 0, 0, 0.2);
      }
      .btn {
        background: -moz-linear-gradient(top,#FFF 0%,#EEE);
        background: -webkit-gradient(linear, left top, left bottom, from(#FFF), to(#EEE));
        border: 1px solid #DDD;
        border-radius: 3px;
        color:#111;
        width: 100px;
        height: 45px;
        padding: 5px 0;
        margin: 0;
      }
      #box {
        width:100%;
        margin:10px;
        auto;
      }
  </style>
  <meta charset="UTF-8">
  <title>Regexp factors demo</title>
  <!-- for IE6-8 support of HTML elements -->
  <!--[if lt IE 9]>
  <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
  <![endif]-->
  <body>
  <div id="center">
  <h1>Regexp factors</h1>
    <p>Computation of a set of necessary factors. </p>
    <p><pre><code>BestFactor(v)
1.    If v = [ | ](l, r) OR v = [・](l, r) Then
2.        (all_l, pref_l, suff_l, fact_l) ← BestFactor(l)
3.        (all_r, pref_r, suff_r, fact_r) ← BestFactor(r)
4.    End of If
5.    If v = [ | ](l, r) Then
6.        Return (all_l ∪ all_r, pref_l ∪ pref_r, suff_l ∪ suff_r, fact_l ∪ fact_r)
7.    Else If v = [・](l, r) Then
8.        Return (all_l・all_r, best(pref_l, all_l・pref_r), best(suff_r, suff_l・all_r), best(fact_l, fact_r, suff_l・pref_r))
9.    Elese If v = [ * ](v) Then Return (θ, θ, θ, θ)
10.   Else If v = α, α ∈ Σ∪{ε} Then Return ({α}, {α}, {α}, {α})
11.   End of if</code></pre>
    <p>
    <p>See. <a href="http://amzn.to/1NJFHRf">Gonzalo Navarro, Mathieu Raffinot, Flexible Pattern Matching in Strings, pp.130-131, 2002</a></p>
  <form class="frm" action="/_demo" method="POST">
    <div id="box">
    <input class="txar" name="re" placeholder="Enter regexp and click analyze." value="{{.Regexp}}" />
    <input class="btn" type="submit" value="analyze"/>
    </div>
  </form>
  {{if .CmdErr}}
    Error: {{.CmdErr}}
  {{else if .Regexp}}
  <table class="tbl">
    <tbody>
      <tr>
        <th>Exact</th>
        <td>{{.Exact}}</td>
      </tr>
      <tr>
        <th>Prefix</th>
        <td>{{.Prefix}}</td>
      </tr>
      <tr>
        <th>Suffix</th>
        <td>{{.Suffix}}</td>
      </tr>
      <tr>
        <th>Fact</th>
        <td>{{.Fact}}</td>
      </tr>
      </tr>
    </tbody>
  </table>
  <div id="graph">
    {{if .GraphSvg}}
    <h3><p>Parse Tree<p></h3>
      {{.GraphSvg}}
    {{end}}
    </div>
  </div>
  {{end}}
  </body>
</html>
`
