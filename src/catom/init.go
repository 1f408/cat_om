package main

import (
	"flag"
	"fmt"
)

func init() {
	flag.CommandLine.Init("CATOM", flag.ContinueOnError)
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\n=^-^= =^-^= =^-^= =^-^= =^-^= =^-^=\n")
		fmt.Fprintf(o, "=^-^= [CATs]Atom Feed for Git =^-^=\n")
		fmt.Fprintf(o, "=^-^= =^-^= =^-^= =^-^= =^-^= =^-^=\n")
		fmt.Fprintf(o, "\nUsage: \n")
		fmt.Fprintf(o, "\tPlease set catom in cron.\n")
		fmt.Fprintf(o, "\t Ex. 0 * * * * cat_om > /dev/null 2>&1\n")
		fmt.Fprintf(o, "\nOptions: \n")
		flag.PrintDefaults()
	}
	flag.StringVar(&env_flag, "e", "./etc/env.toml", "path of env.toml")
}

var (
	env_flag string
)
