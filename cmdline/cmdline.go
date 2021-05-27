// package cmdline provides the command line option description tags and a helper function
// that overrides the default printing behavior of flag.Usage()
package cmdline

import (
	"fmt"
	"os"
)

const (
	DepthDescription    = "Depth level to crawl to. 0 = full crawl, 1 = root level ( / ) crawl"
	InsecureDescription = "Crawl sites with problematic SSL certificates (self-signed, unknown, expired, etc)"
	MimeDescription     = "Include MIME type information with page scrape"
	OutputDescription   = "Output data to a file"
	TimeoutDescription  = "Amount of time (in seconds) to wait for a link to resolve and return content before failing and moving on"
	UrlDescription      = "URL to crawl"
	HelpDescription     = "Display valid options"
)

// GetHelpString builds a string that is printed to the console when the help option is provided,
// or when an incorrect option is supplied
func GetHelpString() string {
	str := fmt.Sprintf("Usage of %s:\n", os.Args[0])
	str = fmt.Sprintf("%s  -d, -depth, --depth [number]\n", str)
	str = fmt.Sprintf("%s      %s\n", str, DepthDescription)
	str = fmt.Sprintf("%s  -i, -insecure, --insecure\n", str)
	str = fmt.Sprintf("%s      %s\n", str, InsecureDescription)
	str = fmt.Sprintf("%s  -m, -mime, --mime\n", str)
	str = fmt.Sprintf("%s      %s\n", str, MimeDescription)
	str = fmt.Sprintf("%s  -o, -output, --output [file]\n", str)
	str = fmt.Sprintf("%s      %s\n", str, OutputDescription)
	str = fmt.Sprintf("%s  -t, -timeout, --timeout\n", str)
	str = fmt.Sprintf("%s      %s\n", str, TimeoutDescription)
	str = fmt.Sprintf("%s  -u, -url, --url [url] (required)\n", str)
	str = fmt.Sprintf("%s      %s\n", str, UrlDescription)
	str = fmt.Sprintf("%s  -h, -help, --help\n", str)
	str = fmt.Sprintf("%s      %s\n", str, HelpDescription)
	return str
}
