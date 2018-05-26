package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

var rgx = regexp.MustCompile("'(.*?)'")
var drupalDir = fs.WorkDir() + "/docroot/core/lib/Drupal.php"
var validDrupalRange = ">=8.5.x"
var checkFailMsg = "\n⚠️  There were some problems detected during the system checks. This won't stop you from running any Tokaido commands, but they may not behave as you were expecting.\n\n"

// Check ...
func Check() {
	if _, err := os.Stat(drupalDir); os.IsNotExist(err) {
		fmt.Println("  ✘  A Drupal installation was not found")
		fmt.Printf(checkFailMsg)
		return
	}

	fmt.Println("  ✓  A Drupal installation was found")

	checkDrupalVersion()
}

// checkDrupalVersion ...
func checkDrupalVersion() {
	drupalVersion := getDrupalVersion()
	if drupalVersion == "" {
		fmt.Printf("  ✘  Tokaido was unable to determine the Drupal version, it should be %s\n", validDrupalRange)
		fmt.Printf(checkFailMsg)
		return
	}

	drupalSemver, _ := semver.Parse(drupalVersion)
	validRange, _ := semver.ParseRange(validDrupalRange)

	if validRange(drupalSemver) {
		fmt.Printf("  ✓  The Drupal version (%s) is supported by Tokaido\n\n", drupalVersion)
	} else {
		fmt.Printf("  ✘  The Drupal version (%s) is not supported by Tokaido. Drupal version %s has been tested and is working with Tokaido\n", drupalVersion, validDrupalRange)
		fmt.Printf(checkFailMsg)
	}
}

// getDrupalVersion ...
func getDrupalVersion() string {
	drupalVersionString := versionStr()
	if drupalVersionString == "" {
		return ""
	}

	drupalVersion := rgx.FindStringSubmatch(drupalVersionString)[0]

	return strings.Replace(drupalVersion, "'", "", 2)
}

// versionStr ...
func versionStr() string {
	f, err := os.Open(drupalDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "const VERSION = '") {
			return scanner.Text()
		}
	}

	return ""
}