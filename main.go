package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/timwmillard/slang/urbandictionary"
)

func main() {

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Printf("%s [OPTIONS] <word>\n", os.Args[0])
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	all := flag.Bool("all", false, "display all definitions")
	auth := flag.Bool("auth", false, "configure API Key")
	example := flag.Bool("example", false, "show example")
	flag.Parse()

	if *auth {
		promptAuth()
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	apiKey, err := readAuth()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading config:", err)
		os.Exit(1)
	}
	if apiKey == "" {
		promptAuth()
	}

	word := flag.Arg(0)

	urbanDict := urbandictionary.NewClient(apiKey)
	definitions, err := urbanDict.Define(word)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Urban Dictionary error:", err)
		os.Exit(1)
	}

	if len(definitions) == 0 {
		fmt.Println("No definitions")
		os.Exit(1)
	}

	if *all {
		printAll(definitions)
	} else {
		printRand(definitions, *example)
	}

}
func printRand(defs []urbandictionary.Definition, showExample bool) {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(defs))
	fmt.Println(defs[index].Definition)
	if showExample {
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println(defs[index].Example)
	}
}
func printAll(defs []urbandictionary.Definition) {
	for _, def := range defs {
		fmt.Println(def.Definition)
		fmt.Println()
	}
}
