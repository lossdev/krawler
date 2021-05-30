package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"krawler/cmdline"
	"krawler/krawl"
)

func main() {
	// parse flags
	setupFlags(flag.CommandLine)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var depth int
	flag.IntVar(&depth, "d", 3, cmdline.DepthDescription)
	flag.IntVar(&depth, "depth", 3, cmdline.DepthDescription)
	var dataFormat string
	flag.StringVar(&dataFormat, "f", "", cmdline.FormatDescription)
	flag.StringVar(&dataFormat, "format", "", cmdline.FormatDescription)
	var insecure bool
	flag.BoolVar(&insecure, "i", false, cmdline.InsecureDescription)
	flag.BoolVar(&insecure, "insecure", false, cmdline.InsecureDescription)
	var mimeScrape bool
	flag.BoolVar(&mimeScrape, "m", false, cmdline.MimeDescription)
	flag.BoolVar(&mimeScrape, "mime", false, cmdline.MimeDescription)
	var outputFile string
	flag.StringVar(&outputFile, "o", "", cmdline.OutputDescription)
	flag.StringVar(&outputFile, "output", "", cmdline.OutputDescription)
	var timeout int
	flag.IntVar(&timeout, "t", 3, cmdline.TimeoutDescription)
	flag.IntVar(&timeout, "timeout", 3, cmdline.TimeoutDescription)
	var url string
	flag.StringVar(&url, "u", "", cmdline.UrlDescription)
	flag.StringVar(&url, "url", "", cmdline.UrlDescription)
	var help bool
	flag.BoolVar(&help, "help", false, cmdline.HelpDescription)
	flag.BoolVar(&help, "h", false, cmdline.HelpDescription)

	flag.Parse()

	// check options
	if help {
		fmt.Fprintf(os.Stderr, cmdline.GetHelpString())
		exit()
	}

	if len(outputFile) != 0 {
		// determine if the outputFile can be opened. If it exists,
		// ask user to approve overwriting the existing file
		if _, err := os.Stat(outputFile); err == nil {
			var inp string
			fmt.Printf("File \"%s\" exists. Overwrite it? (y/N): ", outputFile)
			fmt.Scanln(&inp)
			inp = strings.ToLower(inp)
			if inp != "y" && inp != "yes" {
				fmt.Printf("File \"%s\" will not be overwritten. Exiting\n", outputFile)
				exit()
			}
		} else if errors.Is(err, os.ErrNotExist) {
			// do nothing
		} else {
			log.Printf("Error: stat file %s: %s\n", outputFile, err)
			exit()
		}
	}

	var format = cmdline.DefaultFormat
	if len(dataFormat) != 0 {
		tmp := strings.ToLower(dataFormat)
		if tmp == "json" {
			format = cmdline.JsonFormat
		} else if tmp == "yaml" {
			format = cmdline.YamlFormat
		}
	}

	err := krawl.Init(insecure, timeout, outputFile, format)
	if err != nil {
		log.Printf("Error initializing krawler: %s\n", err)
		exit()
	}

	// check that a url has been provided
	if url == "" {
		fmt.Fprintf(os.Stderr, "Flag required but not provided: -u, -url, --url [url]\n")
		fmt.Fprintf(os.Stderr, cmdline.GetHelpString())
		exit()
	}

	// parse URL into a unified format
	tld, err := krawl.ParseUrl(url)
	if err != nil {
		log.Printf("Error parsing given URL (%s): %s\n", url, err)
		exit()
	}

	// check that url can be resolved
	err = krawl.CheckRootLink(tld)
	if err != nil {
		log.Printf("Error with given URL (%s): %s\n", url, err)
		exit()
	}

	krawl.Krawl(tld, "--", depth, 1, mimeScrape)
	if format != cmdline.DefaultFormat {
		krawl.FlushJson()
	}
}

func setupFlags(f *flag.FlagSet) {
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, cmdline.GetHelpString())
	}
}

func exit() {
	os.Exit(0)
}
