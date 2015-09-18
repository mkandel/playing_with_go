package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/davecgh/go-spew/spew"
	"github.com/headzoo/surf"
)

const Version = "0.0.1"
const Url = "http://xkcd.com/archive"

var (
	debug = flag.Bool("debug", false, "enable debug output")
	dir   = flag.String("dir", "/tmp/xkcd_images", "Directory to drop XKCD images")
)

type Comic struct {
	id   int
	url  string
	name string
	date string
}

func main() {
	flag.Parse()
	spew.Config.Indent = "    "
	fmt.Printf("Version: %v\nUrl: %v\n", Version, Url)

	mech := surf.NewBrowser()
	err := mech.Open(Url)
	check(err)
	body := mech.Body()
	//doc, err := html.Parse(strings.NewReader(s))
	doc, err := html.Parse(strings.NewReader(body))
	//doc, err := gokogiri.ParseHtml([]byte(body))
	//doc, err := gokogiri.ParseHtml([]byte(body))
	check(err)

	//div := doc.NodeById("middleContainer")
	//fmt.Printf("TypeOf div: %v\n", reflect.TypeOf(div))
	//spew.Printf("Div: %+v\n", div)
	//os.Exit(0)

	//content := string(div.Content())
	//fmt.Printf("Div: %+v\n", div.Content())
	//fmt.Println(content)
	//fmt.Println(div.String())

	if err := os.MkdirAll(*dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	b := bytes.NewBufferString(string(doc.Data))
	//b := bytes.NewBufferString(string(body))
	//b := bytes.NewBufferString(string(div.String()))
	z := html.NewTokenizer(b)

	depth := 0
	for {
		tt := z.Next()

		fmt.Println("======================================================================")
		spew.Printf("z: %+v\n", z)
		fmt.Println("======================================================================")
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
			/*
				case tt == html.StartTagToken:
					t := z.Token()

					isAnchor := t.Data == "a"
					if isAnchor {
						fmt.Println("We found a link!")
						spew.Printf("Data: %+v\n", t)
					}
			*/
		case tt == html.TextToken:
			if depth > 0 {
				// emitBytes should copy the []byte it receives,
				// if it doesn't process it immediately.
				fmt.Println("*********")
				fmt.Println(z.Token())
				fmt.Println("*********")
			}
		case tt == html.StartTagToken, tt == html.EndTagToken:
			tn, _ := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' {
				if tt == html.StartTagToken {
					depth++
				} else {
					depth--
				}
			}
		}
	}

	//Rows of interest have the form:
	// <a href="/1577/" title="2015-9-14">Advent</a><br>
}

func check(err error) {
	if err != nil {
		os.Exit(1)
	}
}
