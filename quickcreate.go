package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var source = flag.String("source", "./tsconfig.lib.json", "Source tsconfig json file")
var destination = flag.String("destination", "./tsconfig.lib.json", "Destination tsconfig json file")
var libs = flag.String("libs", "./libs", "the folder to iterate")

func main() {
	flag.Parse()
	fmt.Println("Creating files")

	f := readSource(*source)
	folders := readFolders(*libs)
	modified := []string{}

	for _, folder := range folders {
		if !strings.Contains(folder.Name(), ".") {
			folderName := folder.Name()
			modified = append(modified, folderName)
			lines := strings.Split(string(f), "\n")
			for i, line := range lines {
				if strings.Contains(line, "outDir") {
					lines[i] = `    "outDir": "../../dist/out-tsc/libs/` + folderName + `",`
				}
			}
			output := strings.Join(lines, "\n")
			newpath := path.Join(*libs, folder.Name(), *destination)
			ioutil.WriteFile(newpath, []byte(output), 0644)
		}
	}

	fmt.Println(modified, "created")
}

func readSource(path string) (f []byte) {
	f, err := ioutil.ReadFile(path)
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
