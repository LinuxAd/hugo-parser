package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	cwd := flag.String("d", "cwd", "the directory of files to parse")
	flag.Parse()

	var dir string

	if *cwd == "cwd" || *cwd == "." {
		dir = getCwd()
	} else {
		dir = *cwd
	}

	err := preFlightChecks(dir)
	if err != nil {
		log.Panic(fmt.Sprint("pre-flight checks error:", err))
	}

	files, err := getFileList(dir)
	if err != nil {
		log.Panic(err)
	}

	if err := addFrontMatter(files); err != nil {
		log.Panic(err)
	}

}

func addFrontMatter(files []string) error {
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return err
		}

		log.Println("working on", info.Name())

		lines, err := getFileLines(file)
		if err != nil {
			return err
		}

		title, _, err := getTitle(lines)
		if err != nil {
			return err
		}
		log.Println("title: ", title)

		f := frontMatter{Title: title}

		err = f.addToFile(file)
		if err != nil {
			return err
		}
		log.Println("done")
	}
	return nil
}
