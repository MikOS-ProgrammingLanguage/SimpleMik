package mip_functions

import (
	"MicNewDawn/errors"
	"fmt"
	"os"
	"strings"
)

// Milks a package and returns a string with the name of the package
func Milk(path string, local bool) (string, string) {
	file_contents, read_error := os.ReadFile(path + "milk.pkg")
	CheckError(read_error)

	// remove whitespaces
	file_contents = []byte(strings.ReplaceAll(string(file_contents), "\t", ""))
	file_contents = []byte(strings.ReplaceAll(string(file_contents), " ", ""))

	// split the whole file content into statements
	statements := strings.Split(string(file_contents), "\n")

	var PACKAGE_NAME string
	var ENTRY_FILE string
	for _, statement := range statements {
		if statement == "" {
			continue
		}

		a := strings.SplitN(statement, ":", 2)
		keyword := a[0]
		value := a[1]

		switch keyword {
		case "package-name":
			{
				PACKAGE_NAME = value
			}
		case "ignore":
			{
				MilkIgnore(path, value)
			}
		case "depends":
			{
				if value == "none" {
					continue
				}
				if !local {
					InstallPackage([]string{value})
				} else {
					Add([]string{value})
				}
			}
		case "entry":
			// is the main file of the program which would essentially be compiled
			{
				ENTRY_FILE = value
			}
		default:
			{
				errors.UnknownArgumentError(keyword)
			}
		}
	}

	if PACKAGE_NAME == "" {
		errors.ThrowError(fmt.Sprintf("No package name defined in %smilk.pkg", path), true)
	}
	return PACKAGE_NAME, ENTRY_FILE
}
