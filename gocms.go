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

	nodeMap := make(map[string]*html.Node)
	outputMap := make(map[string]*html.Node)

	// finds all div with class attributes and adds them to nodeMap
	inputTree.Walk(func(n *html.Node) {
		if n.Data == "div" {
			for _, b := range n.Attr {
				if b.Key == "id" {
					nodeMap[b.Val] = n
				}
			}
		}
	})

	// finds all divs with matching class inside template and replaces them with nodes in nodeMap
	templateTree.Walk(func(n *html.Node) {
		if n.Data == "div" {
			for _, b := range n.Attr {
				if b.Key == "id" {
					if node, ok := nodeMap[b.Val]; ok {
						outputMap[b.Val] = node
					}
				}
			}
		}
	})

	for key, node := range outputMap {
		templateTransform.Apply(transform.Replace(h5.CloneNode(node)), fmt.Sprintf("#%s", key))
	}

	finalOutput := templateTransform.String()
	finalBytes := []byte(finalOutput)

	err = ioutil.WriteFile(*outputFile, finalBytes, 0644)
	if err != nil {
		panic(err)
	}
}
