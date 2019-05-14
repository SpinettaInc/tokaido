package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

// IgnoreDefaults - Append the baseline files that should be ignored in source control
func IgnoreDefaults() {
	p := conf.GetConfig().Drupal.Path
	AppendGitignore([]string{
		"docker-compose.tok.yml",
		".tok/local",
		"private/default/*",
		p + "/sites/*/settings.tok.php",
		p + "/sites/*/files/*",
	})
}

// NewGitignore - Generate a fresh .gitignore file
func NewGitignore() {
	c := conf.GetConfig()
	p := c.Drupal.Path
	gi := []byte(`# START Generated by Tokaido

docker-compose.tok.yml
.tok/local
private/default/*
vendor
` + p + `/sites/*/settings.tok.php
` + p + `/sites/*/files/*

` + GitignoreDrupalDefaults() + `

# END Generated by Tokaido
	`)

	fs.TouchByteArray(filepath.Join(conf.GetProjectPath(), "/.gitignore"), gi)
}

// AppendGitignore ...
func AppendGitignore(ignoreList []string) {
	gi := filepath.Join(conf.GetProjectPath(), "/.gitignore")
	if fs.CheckExists(gi) == false {
		fs.TouchEmpty(gi)
	}

	f, err := os.Open(gi)
	if err != nil {
		fmt.Println("There was an issue finding your .gitignore file", err)
		return
	}

	defer f.Close()

	ignoreMap := make(map[string]bool)
	for _, ignoreFile := range ignoreList {
		ignoreMap[ignoreFile] = true
	}

	var containsTokValues = false
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if ignoreMap[scanner.Text()] {
			ignoreMap[scanner.Text()] = false
		}
		if strings.Contains(scanner.Text(), "# START Generated by Tokaido") {
			containsTokValues = true
		}
		buffer.Write([]byte(scanner.Text() + "\n"))
	}

	var ignoreString = ""
	for key := range ignoreMap {
		if ignoreMap[key] == true {
			ignoreString = ignoreString + key + "\n"
		}
	}

	if ignoreString != "" {
		if containsTokValues == true {
			buffer = utils.BufferInsert(buffer, "# START Generated by Tokaido", ignoreString)

			fs.Replace(gi, buffer.Bytes())
			return
		}

		buffer.Write([]byte(`
# START Generated by Tokaido
` + ignoreString + `
# END Generated by Tokaido
`))

		fs.Replace(gi, buffer.Bytes())
	}
}
