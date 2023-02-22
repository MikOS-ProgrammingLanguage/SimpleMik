package mip_functions

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/preprocessor"
	"MicNewDawn/utils"
	"fmt"
	"os"
	"strings"

	cp "github.com/otiai10/copy"
)

// Deletes a certain file from project for cleanup
func MilkIgnore(path, name string) {
	if name == "none" {
		return
	}

	CheckError(os.RemoveAll(path + name))
}

// Returns package names and origins
func GetNamesAndOriginsOfInstalledPackages() ([]string, []string) {
	ORIGINAL_PATH, path_err := os.Getwd()
	CheckError(path_err)

	var PACKAGE_NAMES []string
	var PACKAGE_ORIGINS []string

	// Find home directory and cd into it
	home, home_directory_error := os.UserHomeDir()
	CheckError(home_directory_error)

	CheckError(os.Chdir(home + "/.mik"))

	// Read log file
	a, err := os.ReadFile("./mip_pkg_log.log")
	text := string(a)
	CheckError(err)

	// Check if file has contents at all

	statements := strings.Split(text, "\n")
	if len(statements)-1 == 1 {
		x := strings.SplitAfterN(statements[0], ":::", 2)
		if len(x) == 1 {
			name := []string{""}
			origin := []string{""}
			CheckError(os.Chdir(ORIGINAL_PATH))
			return name, origin
		}
		name := []string{x[0][:len(x[0])-3]}
		origin := []string{x[1]}
		CheckError(os.Chdir(ORIGINAL_PATH))
		return name, origin
	}

	for _, statement := range statements {
		if statement == "" {
			continue
		}
		x := strings.SplitAfterN(statement, ":::", 2)
		PACKAGE_NAMES = append(PACKAGE_NAMES, x[0][:len(x[0])-3])
		PACKAGE_ORIGINS = append(PACKAGE_ORIGINS, x[1])
	}

	CheckError(os.Chdir(ORIGINAL_PATH))
	return PACKAGE_NAMES, PACKAGE_ORIGINS
}

// Returns package names and origins of local packages
func GetNamesAndOriginsOfLocallyInstalledPackages() ([]string, []string) {
	PATH, path_err := os.Getwd()
	CheckError(path_err)

	var PACKAGE_NAMES []string
	var PACKAGE_ORIGINS []string

	// Read log file
	a, err := os.ReadFile(PATH + "/.mik_pkgs/mip_pkg_log.log")
	text := string(a)
	CheckError(err)

	// Check if file has contents at all

	statements := strings.Split(text, "\n")
	if len(statements)-1 == 1 {
		x := strings.SplitAfterN(statements[0], ":::", 2)
		if len(x) == 1 {
			name := []string{""}
			origin := []string{""}
			return name, origin
		}
		name := []string{x[0][:len(x[0])-3]}
		origin := []string{x[1]}
		return name, origin
	}

	for _, statement := range statements {
		if statement == "" {
			continue
		}
		x := strings.SplitAfterN(statement, ":::", 2)
		PACKAGE_NAMES = append(PACKAGE_NAMES, x[0][:len(x[0])-3])
		PACKAGE_ORIGINS = append(PACKAGE_ORIGINS, x[1])
	}

	return PACKAGE_NAMES, PACKAGE_ORIGINS
}

func FetchPackage(path, origin string) {
	// Find home directory and cd into it
	home, home_directory_error := os.UserHomeDir()
	CheckError(home_directory_error)

	CheckError(os.Chdir(home + "/.mik"))

	package_name, entry := Milk(path, false)

	// Check if package Already exists
	names, _ := GetNamesAndOriginsOfInstalledPackages()
	if utils.StringInArray(package_name, names) {
		errors.PackageAlreadyExistsError(package_name)
	}

	// Make package structure
	CheckError(os.MkdirAll(fmt.Sprintf("./mip/pkgs/%s", package_name), 0777))
	CheckError(cp.Copy("./mip/git/", fmt.Sprintf("./mip/pkgs/%s", package_name), cp.Options{AddPermission: 0777}))

	// Preprocess and make a main_package_name.milk
	a, err := os.ReadFile(fmt.Sprintf("./mip/pkgs/%s/%s", package_name, entry))
	text := string(a)
	CheckError(err)

	new_path := fmt.Sprintf("./mip/pkgs/%s/%s", package_name, entry)
	preprocessed_text := preprocessor.Preprocess(&text, &new_path)
	CheckError(os.WriteFile(fmt.Sprintf("./mip/pkgs/%s/main_%s.milk", package_name, package_name), []byte(*preprocessed_text), 0777))

	// Update log file
	b, err2 := os.ReadFile("mip_pkg_log.log")
	file_contents := string(b)
	CheckError(err2)
	CheckError(os.WriteFile("mip_pkg_log.log", []byte(fmt.Sprintf("%s%s:::%s\n", file_contents, package_name, origin)), 0777))
}
