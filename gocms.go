package main

import (
	//"code.google.com/p/go-html-transform/html/transform"
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	var templateFile = flag.String("template", "layout.html", "Template file")
	var inputFile = flag.String("input", "/dev/stdin", "Input file")
	flag.Parse()
	// var outputFile = flag.String("output", "/dev/stdout", "Output file")
	templateBytes, err := ioutil.ReadFile(*templateFile)
	if err != nil {
		panic(err)
	}

	inputBytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		panic(err)
	}

	templateDocStr := string(templateBytes)
	inputDocStr := string(inputBytes)

	fmt.Println(inputDocStr)
	fmt.Println("**************")
	fmt.Println(templateDocStr)

}
