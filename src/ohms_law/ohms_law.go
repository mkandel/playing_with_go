package main

import (
	"flag"
	"fmt"
)

func main() {
	var ohms = flag.Float64("r", 0, "ohms")
	var watts = flag.Float64("w", 0, "watts")
	var amps = flag.Float64("a", 0, "amps")
	var volts = flag.Float64("v", 0, "volts")
	flag.Parse()
	fmt.Printf("O: %f\tW: %f\tA: %f\tV: %f\n", *ohms, *watts, *amps, *volts)
}
