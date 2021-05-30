# krawler
A feature-full command line program to crawl arbitrary sites on the Internet with ease

## Installation

### From Source

```bash
git clone https://github.com/lossdev/krawler; cd krawler
make
make install
```

### Compiled Binary

Compiled binaries are available for the following platforms:

* macOS
  * [amd64](https://github.com/lossdev/krawler/releases/download/v1.0.0/krawler-macos-amd64)
  * [arm64](https://github.com/lossdev/krawler/releases/download/v1.0.0/krawler-macos-arm64)
* Linux
  * [amd64](https://github.com/lossdev/krawler/releases/download/v1.0.0/krawler-linux-amd64)
  * [i386](https://github.com/lossdev/krawler/releases/download/v1.0.0/krawler-linux-i386)
  * [arm64](https://github.com/lossdev/krawler/releases/download/v1.0.0/krawler-linux-arm64)
* Windows
  * [amd64](https://github.com/lossdev/krawler/releases/download/v1.0.0/krawler-windows-amd64)
  * [i386](https://github.com/lossdev/krawler/releases/download/v1.0.0/krawler-windows-i386)

## Usage

Executing `krawler`, `krawler -h`, `krawler -help`, or `krawler --help` will prompt you with the available options.
```bash
[lossdev@foo ~] krawler
Flag required but not provided: -u, -url, --url [url]
Usage of krawler:
  -d, -depth, --depth [number]
      Depth level to crawl to. 0 = full crawl, 1 = root level ( / ) crawl
  -f, -format, --format (json | yaml)
      Output data in JSON or YAML format
  -i, -insecure, --insecure
      Crawl sites with problematic SSL certificates (self-signed, unknown, expired, etc)
  -m, -mime, --mime
      Include MIME type information with page scrape
  -o, -output, --output [file]
      Output data to a file
  -t, -timeout, --timeout
      Amount of time (in seconds) to wait for a link to resolve and return content before failing and moving on
  -u, -url, --url [url] (required)
      URL to crawl
  -h, -help, --help
      Display valid options
[lossdev@foo ~]
```

### Required Options
The only required option is the seed URL parameter (`-u`, `-url`, or `--url`). All others are optional.

### Other Options
* `-d, -depth, --depth [number]`

The depth parameter for how far `krawler` will go relative to the seed URL before stopping. A `0` argument will instruct krawler to crawl as far deep as the site
has links, while an argument of `1` will only crawl the URL given and nothing else. `1` is the default option.

* `-f, -format, --format (json | yaml)`

An argument of `json` or `yaml` will print out the content data as the respective format. If an argument to this option is missing or invalid, it will print in 
the default format.

* `-i, -insecure, --insecure`

Some sites may have a problematic SSL certificate - for example, it could be self signed or expired. This option will bypass any errors that may occur due to these
issues and crawl the site anyway. Use this option if it fits your threat model and at your discretion.

* `-m, -mime, --mime`

Print out the MIME type of a found link. New links will not be followed or crawled if their MIME type is not `text/html`, however, it may be useful to log 
instances of other binary data on a webpage (such as images, pdfs, Word documents, etc)

* `-o, -output, --output [file]`

Use this option to write to a specified file instead of stderr. If the filename provided conflicts with an existing file, `krawler` will explicitly ask for a 
confirmation to overwrite the existing file before continuing.

* `-t, -timeout, --timeout [seconds]`

Each link crawled will wait for 3 seconds before failing and continuing to another link by default. It may be useful to set this option to a higher number on a
slower internet connection or when crawling websites with large content blocks returned, or to a lower number for the impatient (who knows, the option's there).

* `-h, -help, --help`

Will print the help message, which contains short information about how to use `krawler` and each of its options.
