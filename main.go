package main

import (
	"flag"
	"fmt"
	"fql/parse"
	"os"
)

const (
	FQL_VERSION = "1.0.0"
)

var(
	h bool
	sql,delimiter string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&sql, "s", "", "send an `SQL` parse data file")
	flag.StringVar(&delimiter, "d", parse.DEFAULT_DELIMITER, "set the column delimiter, the default Spaces")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `FQL Version: fql/%s
Usage: fql [-h] [-s sql] [-d delim]

Options:
`, FQL_VERSION)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
	}
	if len(sql)>0 {
		parse.GetInstance(sql, delimiter).Parse()
	}
}

