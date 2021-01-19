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
	type args struct {
		files int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"current dir",
			args{
				3,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testDir := tt.name+"-testdir"

			TempDir(t,testDir)
			defer CleanDir(t,testDir)
			want := TempFiles(t, testDir, tt.args.files)

			got, err := getFileList(testDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(want) != len(got) {
				t.Errorf("unequal lengths got = %v, want %v", got, want)
			}
			if !EqualSlices(t, got, want) {
				t.Errorf("getFileList() got = %v, want %v", got, want)
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
			expect := "Test File"
			got, err := getTitle(tt.args)
			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if got != expect {
				t.Errorf("got: \"%v\" expected: \"%v\"", got, expect)
			}
		})
	}

}


func TempDir(t *testing.T, name string) {
	t.Helper()
	if err := os.Mkdir(name, 0755); err != nil {
		t.Errorf("error creating test dir: %v", err)
	}
}

func CleanDir(t *testing.T, name string) {
	t.Helper()
	if err := os.RemoveAll(name); err != nil {
		t.Errorf("error removing dir: %v", err)
	}
}

func TempFiles(t *testing.T,dir string, number int) []string {
	t.Helper()

	var files []string
	patt := "test-*"

	for i := 0; i < number; i++ {
		f, err := ioutil.TempFile(dir, patt)
		if err != nil {
			t.Errorf("error creating tempfile: %v. args: %v %v", err, dir, patt)
		}
		_, err = f.WriteString(t.Name() + " test content")
		if err != nil {
			t.Errorf("error writing test content to file: %v", f.Name())
		}
		_ = f.Close()
		files = append(files, f.Name())
	}
	return files
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