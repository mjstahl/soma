package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"github.com/socialmachines/soma/file"
)

var GetUsage = `Usage:
    soma get <library url(s)>
    
    Retrieves a one or more library files located
    at the specified URLs.
    
    Within a project directory, a reference to the
    library will be stored in the project's 
    "lib/manifest.sm" file.
    
Example (within the project 'Test'):
    $ soma get https://example.com/m/Nil
        retrieve https://example.com/m/Nil
        appended https://example.com/m/Nil => Test/lib/manifest.sm
`

func Get(args []string) {
	if len(args) < 1 {
		fmt.Println("soma get: missing library url(s)")
		fmt.Println(GetUsage)
		os.Exit(1)
	}

	pwd, _ := os.Getwd()
	pd := file.ProjDirFrom(pwd)
	if pd == "" {
		displayGetError("must be called within social machines project", nil)
	}

	for _, url := range args {
		found := brokerHasProjectURL(url)
		if !found {
			fmt.Printf("    NOT FOUND %s\n", url)
			continue
		}

		fmt.Printf("    retrieve %s\n", url)

		writeProjURLToManifest(pd, url)

		name := path.Base(pd)
		fmt.Printf("    appended %s => %s/lib/manifest.sm\n", url, name)
	}
}

func brokerHasProjectURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		displayGetError("error contacting broker", err)
	}

	switch resp.StatusCode {
	case 404:
		return false
	case 200:
		return true
	default:
		return false
	}

	return false
}

func writeProjURLToManifest(dir string, url string) {
	path := path.Join(dir, "lib/manifest.sm")
	m, rerr := ioutil.ReadFile(path)
	if rerr != nil {
		displayGetError("error opening manifest file", rerr)
	}

	var dm *manifest
	if len(m) > 0 {
		uerr := json.Unmarshal(m, &dm)
		if uerr != nil {
			displayGetError("error unmarshaling the manifest", uerr)
		}
	} else {
		dm = &manifest{Libs: []string{}}
	}

	dm.Libs = appendIfMissing(dm.Libs, url)
	b, merr := json.Marshal(dm)
	if merr != nil {
		displayGetError("error marshaling the manifest", merr)
	}

	werr := ioutil.WriteFile(path, b, 0700)
	if werr != nil {
		displayGetError("error writing the manifest", werr)
	}
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}

	return append(slice, i)
}

func displayGetError(msg string, err error) {
	if err != nil {
		fmt.Printf("soma get: %s: %s\n", msg, err)
	} else {
		fmt.Printf("soma get: %s\n", msg)
	}

	os.Exit(1)
}
