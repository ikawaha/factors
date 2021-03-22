# Regexp Necessary Factors

## Installation

```
$ go install github.com/ikawaha/factors@latest
```

To draw parse trees, demo application uses graphviz . You need graphviz installed.

## Usage

```
$ factors --help
command: factors <regexp_pattern>
server:  factors -http=:6060
  -http string
    	HTTP service address (e.g. ':6060')
```

**Example**

```shellsession
$ factors '((GA|AAA)*)(TA|AG)'
Exact: θ
Prefix: θ
Suffix: {AG, TA}
Fragment: {AG, TA}
```

### Web App

![demo](https://raw.githubusercontent.com/wiki/ikawaha/regexp/images/regexp_factors_demo.png)

---
MIT


