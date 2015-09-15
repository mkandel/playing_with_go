package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/headzoo/surf"
	"github.com/moovweb/gokogiri"
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
	doc, err := gokogiri.ParseHtml([]byte(body))
	check(err)

	div := doc.NodeById("middleContainer")
	//spew.Printf("Div: %+v\n", div)
	//content := string(div.Content())
	//fmt.Printf("Div: %+v\n", div.Content())
	//fmt.Println(content)
	//fmt.Println(div.String())

	if err := os.MkdirAll(*dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	/*
		child := doc.Root()
		//child := div.Root().FirstChild().Name()
		//fmt.Printf("Header: %#v\n", child)
		spew.Printf("Child: %+v\n", child)
		next := child.NextSibling()
		spew.Printf("Next: %+v\n", next)
		//fmt.Printf("Header: %v\n", child.String())
	*/

	if *debug {
		fmt.Printf("Mech: %#v\n", mech)
		fmt.Printf("Body: %v\n", string(body))
		fmt.Printf("Dom: %#v\n", mech.Dom())
		fmt.Printf("Div: %+v\n", div)
	}

	//Rows of interest have the form:
	// <a href="/1577/" title="2015-9-14">Advent</a><br>
}

func check(err error) {
	if err != nil {
		os.Exit(1)
	}
}
