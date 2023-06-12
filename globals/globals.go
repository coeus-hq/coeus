package globals

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var Secret = []byte("secret")

var DBNAME = "coeus-sample"

// Hack: Change CWD to project root for both running main and
// testing. Import "coeus/globals" in *_test.go files to set the
// path correctly.
func init() {
	cwd, _ := os.Getwd()
	fmt.Println("CWD: " + cwd)
	if !IsBinary() {
		changeRoot()
	}
	fmt.Println("Executable file: " + os.Args[0])
}

func IsBinary() bool {
	if strings.Contains(os.Args[0], "coeus-bin") {
		return true
	} else {
		return false
	}
}

func changeRoot() {
	// See if a filename ROOT exists. If not we are in a package folder.
	// Chdir one level up.
	_, err := os.Stat("./ROOT")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Changing dir to ROOT")
		_ = os.Chdir("../")
		cwd, _ := os.Getwd()
		fmt.Println("CWD: " + cwd)
	}
}
