package ronn

import (
	"os"
	"path/filepath"
)

func example(name string) string {
	return filepath.Join(examples, name)
}

var examples string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	examples = filepath.Join(wd, "examples")
}
