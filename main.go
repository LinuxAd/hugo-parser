package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type frontMatter struct {
	Title string
}

func getFileLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func InsertStringToFile(path, str string, index int) error {
	lines, err := getFileLines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

func preFlightChecks(cwd string) error {
	if !filepath.IsAbs(cwd) {
		return errors.Errorf("filepath should be absolute, got: %v", cwd)
	}
	if !strings.Contains(cwd, "hugo") && !strings.Contains(cwd, "docs") {
		return errors.Errorf("can't find hugo or docs in wd: %v", cwd)
	}
	return nil
}

func getCwd() string {
	d, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	return d
}

func getFileList(cwd string) ([]string, error) {
	var files []string
	err := filepath.Walk(cwd, func(p string, info os.FileInfo, err error) error {
		if info.Name() == "hugo"{
			return nil
		}
		files = append(files, info.Name())
		return nil
	})

	return files,err
}

func getTitle(lines []string) (string, error) {
	y := len(lines) / 2 - 1

	var title string
	var err error

	pat := `^#(.+)$`
	r, err := regexp.Compile(pat)
	if err != nil {
		return title, err
	}

	for x := 0; x < y; x++ {
		l := lines[x]
		s := r.FindStringSubmatch(l)
		if s[1] != "" {
			return s[1], nil
		}
	}

	return title, err
}

func main() {
	cwd := getCwd()
	err := preFlightChecks(cwd)
	if err != nil {
		log.Panic(fmt.Sprint("pre-flight checks error:",err))
	}

	lines, err := getFileLines(cwd + "nft.md")
	if err != nil {
		log.Panic(err)
	}

	title, err := getTitle(lines)
	if err != nil {
		log.Panic(err)
	}

	log.Println(title)


}
