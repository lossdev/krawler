// package krawl implements many methods to successfully traverse a given root webpage up to a depth d (which
// is supplied by the user). MIME types can also be extracted from each page.
package krawl

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bobesa/go-domain-util/domainutil"
	"golang.org/x/net/html"
)

var (
	c            *http.Client
	fileHandler  *os.File
	visitedLinks = make(map[string]bool)
	tld          string
	subdomain    string
)

func Init(insecure bool, timeout int, outputFile string) error {
	// set up HTTP client
	if insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c = &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: tr,
		}
	} else {
		c = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}

	// set up output handler (stderr or output file provided)
	var err error
	if len(outputFile) != 0 {
		fileHandler, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
	} else {
		fileHandler = os.Stderr
	}
	return nil
}

func ParseUrl(link string) (string, error) {
	if !strings.Contains(link, "http://") && !strings.Contains(link, "https://") {
		link = "http://" + link
	}

	// remember the user supplied link's domain and subdomain
	// when new links are found, compare their domain and subdomain
	// to the ones on file - if they match, then the link will be crawled
	// else, it will be discarded to avoid infinite and/or re-crawling
	tld = domainutil.Domain(link)
	subdomain = domainutil.Subdomain(link)

	var err error
	baseUrl, err := url.ParseRequestURI(link)
	if err != nil {
		return "", err
	}
	return baseUrl.String(), nil
}

func CheckRootLink(link string) error {
	resp, err := c.Head(link)
	if err != nil {
		return err
	}
	ctype := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	if ctype != "text/html" {
		e := fmt.Sprintf("Supplied link is not an html document (expected text/html, got %s)", ctype)
		return errors.New(e)
	}
	return nil
}

func Krawl(pageUrl string, parentLink string, desiredDepth int, currentDepth int, mimeScrape bool) {
	// standardize URL format
	// curPageUrlStruct will be used to resolve relative paths found in links on this page
	// against the current page. For example, if we crawl "http://foo.com/bar/baz" and find a
	// link to "../bat", then that relative link to "../bat" will be transformed into
	// "http://foo.com/bar/bat"

	curPageUrlStruct, err := url.Parse(pageUrl)
	if err != nil {
		log.Printf("Error parsing current page URL %s: %s\n", pageUrl, err)
		return
	}
	pageUrl = curPageUrlStruct.String()

	if _, exists := visitedLinks[pageUrl]; exists {
		return
	} else {
		visitedLinks[pageUrl] = true
	}

	// perform HEAD request on current link
	resp, err := c.Head(pageUrl)
	if err != nil {
		log.Printf("Error performing HEAD request on %s: %s\n", pageUrl, err)
		return
	}
	contentType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]

	fmt.Fprintf(fileHandler, "--------------------------\nURL: %s\n", pageUrl)
	fmt.Fprintf(fileHandler, "Parent Link: %s\n", parentLink)
	fmt.Fprintf(fileHandler, "Content-Type: %s\n", contentType)
	fmt.Fprintf(fileHandler, "Depth: %d\n", currentDepth)
	fmt.Fprintf(fileHandler, "--------------------------\n")

	if contentType != "text/html" {
		return
	} else {
		// add a slash to path (if applicable) to ensure relative links can resolve
		if string(pageUrl[len(pageUrl)-1]) != "/" {
			pageUrl += "/"
		}
		curPageUrlStruct, err = url.Parse(pageUrl)
		if err != nil {
			log.Printf("Error parsing current page URL: %s: %s\n", pageUrl, err)
			return
		}
	}

	doc, reader, err := getGoqueryDoc(pageUrl)
	if err != nil {
		log.Printf("Error getting goquery doc for %s: %s\n", pageUrl, err)
		return
	}
	domTokens := html.NewTokenizer(reader)
	previousToken := domTokens.Token()

tokenLoop:
	for {
		tt := domTokens.Next()
		switch {
		case tt == html.ErrorToken:
			break tokenLoop // End of the document,  done
		case tt == html.StartTagToken:
			previousToken = domTokens.Token()
		case tt == html.TextToken:
			// ignore <script>, <noscript>, and <style> text elements
			ignoreTags := []string{"script", "noscript", "style"}
			m := make(map[string]bool)
			for i := 0; i < len(ignoreTags); i++ {
				m[ignoreTags[i]] = true
			}
			// if previousToken.Data is one of the tags we don't want, ignore it
			if _, ok := m[previousToken.Data]; ok {
				continue
			}
			content := strings.TrimSpace(html.UnescapeString(string(domTokens.Text())))
			if len(content) > 0 {
				fmt.Fprintf(fileHandler, "%s\n", content)
			}
		}
	}

	// if current depth is equal to desired depth, return and don't find new links
	if currentDepth == desiredDepth {
		fmt.Fprintf(fileHandler, "--------------------------\n\n")
		return
	}

	var links []string
	doc.Find("a").Each(func(index int, html *goquery.Selection) {
		foundLink, _ := html.Attr("href")
		// remove 'www.' from link (if it is there), so that domain/sd parsing doesn't
		// pick it up as a sd
		foundLink = strings.Replace(foundLink, "www.", "", 1)
		// resolve any link found against base link - will convert relative to absolute paths
		foundUrlStruct, err := url.Parse(foundLink)
		if err != nil {
			log.Printf("Error parsing found URL %s: %s\n", foundLink, err)
			return
		}
		resolvedLink := curPageUrlStruct.ResolveReference(foundUrlStruct).String()

		foundLinkTld := domainutil.Domain(resolvedLink)
		foundLinkSD := domainutil.Subdomain(resolvedLink)

		if (foundLinkTld == tld) && (foundLinkSD == subdomain) && (resolvedLink != pageUrl) {
			// valid link to follow
			fmt.Fprintf(fileHandler, "Found link: %s\n", resolvedLink)
			links = append(links, resolvedLink)
		}
	})
	fmt.Fprintf(fileHandler, "--------------------------\n\n")

	for i := 0; i < len(links); i++ {
		Krawl(links[i], pageUrl, desiredDepth, currentDepth+1, mimeScrape)
	}
}

func getGoqueryDoc(link string) (*goquery.Document, io.Reader, error) {
	resp, err := c.Get(link)
	if err != nil {
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	resp.Body.Close()

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	reader := bytes.NewBuffer(body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return doc, reader, nil
}
