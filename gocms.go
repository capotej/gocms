package main

import (
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"io/ioutil"
)

var (
	templateFile = flag.String("template", "layout.html", "Template file")
	inputFile    = flag.String("input", "/dev/stdin", "Input file")
	outputFile   = flag.String("output", "/dev/stdout", "Output file")
)

func init() {
	flag.Parse()
}

func bytesToTree(inputBytes []byte) *h5.Tree {
	inputReader := bytes.NewReader(inputBytes)
	inputTree, err := h5.New(inputReader)
	if err != nil {
		panic(err)
	}
	return inputTree
}

func bytesFromFile(filename *string) []byte {
	templateBytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
	return templateBytes
}

func treeFromFile(filename *string) *h5.Tree {
	return bytesToTree(bytesFromFile(filename))
}

func findDivId(n *html.Node) bool {
	var result = false
	if n.Data == "div" {
		for _, b := range n.Attr {
			if b.Key == "id" {
				result = true
				return result
			}
		}
	}
	return result
}

func writeResult(templateTransform *transform.Transformer, filename *string) {
	finalOutput := templateTransform.String()
	finalBytes := []byte(finalOutput)

	err := ioutil.WriteFile(*filename, finalBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func ProcessTemplateWithInput(inputFile *string, templateFile *string, outputFile *string) {

	inputTree := treeFromFile(inputFile)
	templateTree := treeFromFile(templateFile)

	templateTransform := transform.New(templateTree)

	inputTree.Walk(func(n *html.Node) {
		if findDivId(n) {
			templateTransform.Apply(transform.Replace(h5.CloneNode(n)), fmt.Sprintf("#%s", n.Attr[0].Val))
		}
	})

	writeResult(templateTransform, outputFile)
}

func main() {
	ProcessTemplateWithInput(inputFile, templateFile, outputFile)
}
