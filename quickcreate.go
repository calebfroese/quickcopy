package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var source = flag.String("source", "./tsconfig.json", "Source tsconfig json file")
var libs = flag.String("libs", "./libs", "the folder to iterate")

func main() {
	flag.Parse()
	fmt.Println("Creating files")

	f := readSource(*source)
	libs := readFolders(*libs)

	for _, f := range libs {
		if !strings.Contains(f.Name(), ".") {

			fmt.Println(f.Name())
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "outDir") {
			fmt.Println(line)
		}
	}
}

func readSource(path string) (f *os.File) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func readFolders(path string) (folders []os.FileInfo) {
	folders, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return folders
}
