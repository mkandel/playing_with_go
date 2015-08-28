package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
)

const Version = "0.0.1"
const Url = "http://xkcd.com/archive"

func main() {
	spew.Config.Indent = "\t"
	spew.Printf("Version: %v\nUrl: %v\n", Version, Url)

	resp, err := http.Get(Url)
	check(err)

	body, err := ioutil.ReadAll(bufio.NewReader(resp.Body))
	defer resp.Body.Close()
	if err != nil {
		spew.Printf("Response: %#+v\n", resp)
	}
	check(err)

	fmt.Printf("Body: %v\n", string(body))
	fmt.Print("Exiting\n")
}

func check(err error) {
	if err != nil {
		os.Exit(1)
	}
}
