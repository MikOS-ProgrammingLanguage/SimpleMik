package mip_functions

import (
	"MicNewDawn/errors"
	"MicNewDawn/utils"
	"fmt"
	"os"
	"strings"
)

func RemovePackage(names []string) {
	for _, name := range names {
		// Find home directory and cd into it
		home, home_directory_error := os.UserHomeDir()
		CheckError(home_directory_error)

		CheckError(os.Chdir(home))
		// Cd into .mik directory
		CheckError(os.Chdir(".mik"))

		// Get all names and origins
		NAMES, _ := GetNamesAndOriginsOfInstalledPackages()
		// Check if package is installed
		if !utils.StringInArray(name, NAMES) {
			errors.PackageNotFoundError(name)
			continue
		}

		// Remove the package
		CheckError(os.RemoveAll(fmt.Sprintf("./mip/pkgs/%s", name)))

		// Update mip_pkg_log.log
		NEW_PACKAGE_LOG_TEXT := ""
		a, err := os.ReadFile("mip_pkg_log.log")
		CheckError(err)
		package_log_text := strings.Split(string(a), "\n")

		for _, text_segment := range package_log_text {
			if strings.HasPrefix(text_segment, name) {
			} else {
				NEW_PACKAGE_LOG_TEXT += text_segment
			}
			NEW_PACKAGE_LOG_TEXT += "\n"
		}

		CheckError(os.WriteFile("mip_pkg_log.log", []byte(NEW_PACKAGE_LOG_TEXT[:len(NEW_PACKAGE_LOG_TEXT)-1]), 0777))
	}
}
