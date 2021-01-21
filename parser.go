package main

import (
	"bufio"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	link = `\[(.+)\]\(.+\)`
	header = `^#(.+)$`
)

type frontMatter struct {
	Title string
}

func (f *frontMatter) MarshalYAML() ([]byte, error) {
	b := []byte("---\n")
	c, err := yaml.Marshal(f)
	c = append(c, "\n---\n"...)
	b = append(b, c...)
	return b, err
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

func alreadyPresent(lines []string, str string) bool {
	if len(lines) < 2 {
		return false
	}
	if strings.Contains(lines[0], "---") || strings.Contains(lines[1], "---") {
		return true
	}
	return false
}

func (f *frontMatter) addToFile(path string) error {
	b, err := f.MarshalYAML()
	if err != nil {
		return err
	}
	str := string(b)

	lines, err := getFileLines(path)
	if err != nil {
		return err
	}
	h := regexp.MustCompile(header)

	fileContent := ""
	for i, line := range lines {
		// if it's the first line in file and the frontMatter isn't already there, add it
		if i == 0 && !alreadyPresent(lines, str) {
			fileContent += str
		}
		// if the line is a main header and contains the title string, remove it
		if h.MatchString(line) && strings.Contains(line, f.Title) {
			line = ""
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

func getExt(file string) string {
	x := strings.Split(file, ".")
	return x[len(x)-1]
}

func getFileList(cwd string) ([]string, error) {
	var files []string
	w := func(p string, info os.FileInfo, err error) error {
		name := info.Name()
		// switch for exclusions not related to names
		switch {
		case info.IsDir():
			return nil
		case strings.Contains(p, ".git"):
			return nil
		case getExt(name) != "md":
			return nil
		}
		// switch for excluding named files/dirs
		switch name {
		case ".git", "hugo", ".DS_Store":
			return nil
		}

		files = append(files, p)
		return nil
	}

	err := filepath.Walk(cwd, w)

	return files, err
}

func getTitle(lines []string) (string, int, error) {
	y := len(lines)/2 + 1 // adding plus one helps to avoid problems if file is only 1 line

	var title string
	var index int
	var err error

	r := regexp.MustCompile(header)

	for x := 0; x < y; x++ {
		l := lines[x]
		// get regex match groups for the line at index x
		s := r.FindStringSubmatch(l)
		if len(s) < 1 {
			continue
		}
		if s[1] != "" {
			title = titleFormatter(title)
			index = x
			return title, index, err
		}
	}

	return title, index, err
}

func titleFormatter(title string) string {
	l := regexp.MustCompile(link)
	t := strings.TrimSpace(title)
	if l.MatchString(t) {
		// title is a link
		g := l.FindStringSubmatch(t)
		return strings.Title(g[1])
	}
	return strings.Title(t)
}
