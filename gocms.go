package main

import (
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	srcDir  = flag.String("src", "src", "Source directory")
	destDir = flag.String("dest", "output", "Destination Directory")
	server  = flag.Bool("server", false, "Start a preview server")
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

func bytesFromFile(filename string) []byte {
	templateBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return templateBytes
}

func treeFromFile(filename string) *h5.Tree {
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

func writeResultToFile(result string, filename string) {
	finalBytes := []byte(result)

	err := ioutil.WriteFile(filename, finalBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func ProcessTemplateWithInput(inputFile string, templateFile string) string {
	inputTree := treeFromFile(inputFile)
	templateTree := treeFromFile(templateFile)

	templateTransform := transform.New(templateTree)

	inputTree.Walk(func(n *html.Node) {
		if findDivId(n) {
			templateTransform.Apply(transform.Replace(h5.CloneNode(n)), fmt.Sprintf("#%s", n.Attr[0].Val))
		}
	})

	return templateTransform.String()
}

func visitFileFunc(filePath string, f os.FileInfo, err error) error {
	ext := filepath.Ext(filePath)
	inputPath := strings.Replace(filePath, *srcDir, "", 1)
	outputPath := filepath.Join(*destDir, inputPath)
	templateFile := filepath.Join(*srcDir, "layout.html")

	if ext == ".html" && inputPath != "/layout.html" {
		outputDir := filepath.Dir(outputPath)
		err := os.MkdirAll(outputDir, 0777)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s -> %s\n", inputPath, outputPath)
		result := ProcessTemplateWithInput(filePath, templateFile)
		writeResultToFile(result, outputPath)
	}

	return nil
}

func GeneratorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ProcessDirectoryWithTemplate(*srcDir, *destDir)
		h.ServeHTTP(w, r)
	})
}

func ProcessDirectoryWithTemplate(srcDir string, destDir string) {
	err := filepath.Walk(srcDir, visitFileFunc)
	if err != nil {
		panic(err)
	}
}

func main() {
	if *server {
		fileHandler := http.FileServer(http.Dir(*destDir))
		http.ListenAndServe(":8080", GeneratorHandler(fileHandler))
	} else {
		ProcessDirectoryWithTemplate(*srcDir, *destDir)
	}
}
