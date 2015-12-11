[![experimental](http://badges.github.io/stability-badges/dist/experimental.svg)](http://github.com/badges/stability-badges)

#Regexp Factors

## Installation

```
$ go get github.com/ikawaha/regexp/...
```

## Usage

```
$ refact
usage: refact regexp
       refact -http=:6060
  -http string
        HTTP service address (e.g. ':6060')
```

### Web App

![demo](https://raw.githubusercontent.com/wiki/ikawaha/regexp/images/regexp_factors_demo.png)

### Programing Example

```
package main

import (
	"fmt"
	"os"
	"regexp/syntax"

	"github.com/ikawaha/regexp/factor"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(2)
	}
	re, err := syntax.Parse(os.Args[1], syntax.Perl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	f := factor.Analyze(re)
	fmt.Println("Exact", f.Exact)
	fmt.Println("Prefix", f.Prefix)
	fmt.Println("Suffix", f.Suffix)
	fmt.Println("Fact", f.Fact)
}
```
