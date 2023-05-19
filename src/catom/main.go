package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	catom "catom/feed"
)

func dyingMsg(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}

func main() {
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		os.Exit(0)
	}

	cfg, err := catom.NewConfig(env_flag)
	if err != nil {
		dyingMsg("Config open error: %s", err)
	}

	now := time.Now()
	gEntries, err := catom.Feed4git(cfg.Dotgit, cfg.Urlroot, cfg.Root, now, cfg.Diff)
	if err != nil {
		dyingMsg("Failed create data for Feed: %s", err)
	}

	feed, err := cfg.NewAtom(now, gEntries)
	if err != nil {
		dyingMsg("Create Feed error: %s\n", err)
	}

	perm := "0644"
	outpath := cfg.Outpath + "/" + cfg.Outfile
	err = catom.WriteFeed(feed, outpath, perm)
	if err != nil {
		dyingMsg("Write Feed error: %s\n", err)
	}
}
