package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestPreFlightChecks(t *testing.T) {
	got := preFlightChecks("/tmp")
	if got == nil {
		t.Error("preFlight check should error, got nil")
	}
}

func TestGetCwd(t *testing.T) {
	cwd, _ := os.Getwd()
	if got := getCwd(); got != cwd {
		t.Errorf("got: %v, expected %v", got, cwd)
	}
}

func TestLinesFromReader(t *testing.T) {
	type TestCase struct{
		testString string
		expected []string
	}
	tt := TestCase {
		 "---\ntitle: \"A Title\"\n",
		[]string{
			"---",
			"title: \"A Title\"",
		},
	}

	t.Run("test file read", func(t *testing.T){
		got, err := LinesFromReader(strings.NewReader(tt.testString))
		if err != nil {
			t.Errorf("LinesFromReader error: %v", err)
		}
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("got %v, expected %v", got, tt.expected)
		}
	})
}

func TestGetFileList(t *testing.T) {
	cwd, _ := os.Getwd()

	type args struct {
		files int
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid dir",
			args{
				3,
				cwd,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d, err := ioutil.TempDir(tt.args.dir, "")
			if err != nil {
				t.Errorf("error creating temp directory: %v", err)
			}
			t.Cleanup(func() {
				CleanDir(t, d)
			})

			want := TempFiles(t, d, tt.args.files)

			got, err := getFileList(d)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(want) != len(got) {
				t.Errorf("unequal lengths got: %v want: %v", len(got), len(want))
			}
			if !EqualSlices(t, got, want) {
				t.Log("got")
				for _, i := range got {
					t.Log(i)
				}
				t.Log("want")
				for _, i := range want {
					t.Log(i)
				}
				t.Errorf("slices not equal")
			}
		})
	}
}

func TestGetTitle(t *testing.T) {
	tests := []struct{
		name string
		args []string
		want string
		wantErr bool
	}{
		{
			name: "quick test",
			args: []string{
				"# A Title",
				"some content",
			},
			want: "A Title",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTitle(tt.args)
			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if got != tt.want{
				t.Errorf("got: \"%v\" expected: \"%v\"", got, tt.want)
			}
		})
	}

}

func TempFiles(t *testing.T,dir string, number int) []string {

	t.Helper()
	var files []string
	for i := 0; i < number; i++ {
		f, _ := ioutil.TempFile(dir, "*.md")
		files = append(files, f.Name())
	}
	return files
}

func CleanDir(t *testing.T, name string) {
	t.Helper()
	if err := os.RemoveAll(name); err != nil {
		t.Errorf("error removing dir: %v", err)
	}
}

func EqualSlices(t *testing.T, a, b []string) bool {
	t.Helper()
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !InSlice(t, b, v){
			return false
		}
	}
	return true
}

func InSlice(t *testing.T, slice []string, element string) bool {
	t.Helper()
	for _, v := range slice {
		if element == v {
			return true
		}
	}
	return false
}

func Test_titleFormatter(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal header",
			args: args{
				title: "Easy Header",
			},
			want: "Easy Header",
		},
		{
			name: "add whitespace",
			args: args{
				title: " whitespace ",
			},
			want: "Whitespace",
		},
		{
			name: "link",
			args: args{
				title: "[A Link](https://www.link.com)",
			},
			want: "A Link",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := titleFormatter(tt.args.title); got != tt.want {
				t.Errorf("titleFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}