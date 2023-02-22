package mip_functions

import (
	"fmt"
	"os"
	"strings"
)

// Lists all installed packages
func ListPackages() {
	// List all packages installed on the system and all packages installed in a project if found

	// On System
	SYSTEM_NAMES, SYSTEM_ORIGINS := GetNamesAndOriginsOfInstalledPackages()

	fmt.Printf("SYSTEM (%d)\n│\n", len(SYSTEM_NAMES))

	if len(SYSTEM_NAMES) == 0 {
		fmt.Println("└── NONE")
	} else {
		for index := 0; index < len(SYSTEM_NAMES)-1; index++ {
			fmt.Printf("├── %s : %s\n", SYSTEM_NAMES[index], SYSTEM_ORIGINS[index])
		}
		fmt.Printf("└── %s : %s\n", SYSTEM_NAMES[len(SYSTEM_NAMES)-1], SYSTEM_ORIGINS[len(SYSTEM_ORIGINS)-1])
	}

	// In project if found
	// if package file exists
	var PACKAGE_NAME string

	path, err := os.Getwd()
	CheckError(err)

	file_contents, read_error := os.ReadFile(path + "/milk.pkg")
	if read_error != nil {
		if read_error.Error() == fmt.Sprintf("open %s: no such file or directory", path+"/milk.pkg") {
			return
		}
	}

	// remove whitespaces
	file_contents = []byte(strings.ReplaceAll(string(file_contents), "\t", ""))
	file_contents = []byte(strings.ReplaceAll(string(file_contents), " ", ""))

	// split the whole file content into statements
	statements := strings.Split(string(file_contents), "\n")

	for _, statement := range statements {
		if strings.HasPrefix(statement, "package-name:") {
			a := strings.SplitAfterN(statement, ":", 2)
			_, PACKAGE_NAME = a[0], a[1]
		}
	}

	LOCAL_PACKAGE_NAMES, LOCAL_PACKAGE_ORIGINS := GetNamesAndOriginsOfLocallyInstalledPackages()

	fmt.Printf("\nPROJECT: %s (%d)\n│\n", PACKAGE_NAME, len(LOCAL_PACKAGE_NAMES))

	if len(LOCAL_PACKAGE_NAMES) == 0 {
		fmt.Println("└── NONE")
	} else {
		for index := 0; index < len(LOCAL_PACKAGE_NAMES)-1; index++ {
			fmt.Printf("├── %s : %s\n", LOCAL_PACKAGE_NAMES[index], LOCAL_PACKAGE_ORIGINS[index])
		}
		fmt.Printf("└── %s : %s\n", LOCAL_PACKAGE_NAMES[len(LOCAL_PACKAGE_NAMES)-1], LOCAL_PACKAGE_ORIGINS[len(LOCAL_PACKAGE_ORIGINS)-1])
	}
}
