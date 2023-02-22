package mip_functions

import (
	"MicNewDawn/errors"
	"fmt"
	"os"
	"strings"
)

// Checks if error occurred and throws it
func CheckError(e error) {
	if e != nil {
		//errors.ThrowError(e.Error(), true)
		panic(e)
	}
}

// Analyzes the command further and decides how to initialize appropriately
func Init(args []string) {
	if len(args) < 2 {
		errors.ToFewArgumentError()
	}

	switch args[0] {
	case "mik":
		{
			init_normal_mik_project(args[1])
		}
	case "milk":
		{
			init_milk_package_project(args[1])
		}
	default:
		{
			errors.UnknownArgumentError(args[0])
		}
	}
}

func base_init(path string) {
	/*
		- .gitignore
		- milk.pkg
		- .mik_pkgs/
		- .mik_pkgs/mip_pkg_log.log
		- src/
	*/

	CheckError(os.MkdirAll(path+".mik_pkgs", 0755))
	CheckError(os.WriteFile(path+".mik_pkgs/mip_pkg_log.log", []byte(""), 0755))
	CheckError(os.MkdirAll(path+"src", 0755))

	CheckError(os.WriteFile(path+".gitignore", []byte("/.mik_pkgs/\n/src/\n"), 0755))
	CheckError(os.WriteFile(path+"milk.pkg", []byte(fmt.Sprintf("package-name: %s\nignore: none\ndepends: none", path[:len(path)-1])), 0755))
}

func init_normal_mik_project(path string) {
	/*
		- Makefile
		- main.mik
	*/

	// add '/' if not suffix of path
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	base_init(path)
	CheckError(os.WriteFile(path+"Makefile", []byte("CC = mic -i\n\nall:\n\tCC main.mik\n"), 0755))
	CheckError(os.WriteFile(path+"main.mik", []byte(""), 0755))

	a, err := os.ReadFile(path + "milk.pkg")
	file_contents := string(a)
	CheckError(err)
	CheckError(os.WriteFile(path+"milk.pkg", []byte(file_contents+"\nentry: main.mik"), 0755))
}

func init_milk_package_project(path string) {
	/*
		- main.milk
		- milk.pkg (update)
	*/

	// add '/' if not suffix of path
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	base_init(path)
	CheckError(os.WriteFile(path+"main.milk", []byte(""), 0755))

	a, err := os.ReadFile(path + "milk.pkg")
	file_contents := string(a)
	CheckError(err)
	CheckError(os.WriteFile(path+"milk.pkg", []byte(file_contents+"\nentry: main.milk"), 0755))
}
