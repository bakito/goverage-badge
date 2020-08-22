package mod

import (
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"strings"
)

// Get the repo of the current mod file
func Repo() (string, error) {
	return repoOf("./go.mod")
}

func repoOf(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	f, _ := modfile.Parse("", b, nil)
	if err != nil {
		return "", err
	}

	path := f.Module.Mod.Path
	if strings.HasPrefix(path, "github.com/") {
		return strings.Replace(path, "github.com/", "", 1), nil
	}
	return "", nil
}
