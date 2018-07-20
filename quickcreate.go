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
	fn := os.Args[1]
	flag.Parse()
	switch {
	case fn == "tsconfig":
		tsconfig()
	case fn == "index":
		index()
	default:
		log.Fatal("Please provide one of tsconfig,index")
	}
	if fn != "" {
		fmt.Println(fn)
		return
	}
}

func index() {
	folders := readFolders(*libs)
	for _, folder := range folders {
		if !strings.Contains(folder.Name(), ".") {
			// move the file
			folderName := folder.Name()
			oldpath := path.Join(*libs, folderName, "index.ts")
			newpath := path.Join(*libs, folderName, "src", "index.ts")
			fmt.Println(oldpath, newpath)
			err := os.Rename(oldpath, newpath)
			if err != nil {
				fmt.Println("An error occurred")
				log.Fatal(err)
			}

			// update the contents
			fmt.Println("Updating index.ts contents")
			lines := readSource(newpath)
			replacer := strings.NewReplacer("/src", "")
			for i, line := range lines {
				if strings.Contains(line, "/src") {
					lines[i] = replacer.Replace(line)
				}
			}
			output := strings.Join(lines, "\n")
			ioutil.WriteFile(newpath, []byte(output), 0644)
		}
	}
}

func tsconfig() {
	lines := readSource(*source)
	folders := readFolders(*libs)
	modified := []string{}

	for _, folder := range folders {
		if !strings.Contains(folder.Name(), ".") {
			folderName := folder.Name()
			modified = append(modified, folderName)
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

func readSource(path string) []string {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(f), "\n")
}

func readFolders(path string) (folders []os.FileInfo) {
	folders, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return folders
}
