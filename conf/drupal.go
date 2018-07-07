package conf

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetRootPath ...
func GetRootPath() string {
	wd := fs.WorkDir()
	c := GetConfig().Drupal.Path
	if c != "" {
		sa := []rune(c)
		if string(sa[0]) != "/" {
			c = "/" + c
		}
		return wd + c
	}

	rootPath := scanForCoreDrupal()

	CreateOrReplaceDrupalPath(strings.Replace(rootPath, wd, "", -1))

	return rootPath
}

func scanForCoreDrupal() string {
	wd := fs.WorkDir()
	var dp string
	dc := "/core/lib/Drupal.php"
	err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, dc) == true {
			dp = strings.Replace(path, dc, "", -1)
			return io.EOF
		}
		return nil
	})
	if err != io.EOF {
		log.Fatalf("There was an error when searching for your Drupal installation [%v]\n", err)
	}

	return dp
}

// CoreDrupalFile - Return the core drupal file for the users' installation
func CoreDrupalFile() string {
	return GetRootPath() + "/core/lib/Drupal.php"
}

// GetRootDir - Return the drupal root folder name without workdir
func GetRootDir() string {
	dr := GetRootPath()
	ds := strings.Split(dr, "/")

	return ds[len(ds)-1]
}