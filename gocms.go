package main

import (
	// "code.google.com/p/go-html-transform/css/selector"
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	var templateFile = flag.String("template", "layout.html", "Template file")
	var inputFile = flag.String("input", "/dev/stdin", "Input file")
	var outputFile = flag.String("output", "/dev/stdout", "Output file")
	flag.Parse()

	templateBytes, err := ioutil.ReadFile(*templateFile)
	if err != nil {
		panic(err)
	}

	inputBytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		panic(err)
	}

	inputReader := bytes.NewReader(inputBytes)
	inputTree, _ := h5.New(inputReader)

	templateReader := bytes.NewReader(templateBytes)
	templateTree, _ := h5.New(templateReader)
	templateTransform := transform.New(templateTree)

	fmt.Printf("%#v\n", templateTransform)

	inputTree.Walk(func(n *html.Node) {
		if n.Data == "div" {
			// templateTransform.Apply(Replace(Text("my new text"), "a")
			fmt.Printf("%#v\n", n)
		}
	})

	err = ioutil.WriteFile(*outputFile, inputBytes, 0644)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(*outputFile, templateBytes, 0644)
	if err != nil {
		panic(err)
	}
}
