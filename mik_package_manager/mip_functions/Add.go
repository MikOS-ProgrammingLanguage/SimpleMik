package mip_functions

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/preprocessor"
	"MicNewDawn/utils"
	"fmt"
	"os"
	"os/exec"

	cp "github.com/otiai10/copy"
)

// Installs a package to the local packages
func Add(urls []string) {
	// check if in project
	path, err := os.Getwd()
	mik_pkgs_path := path + "/.mik_pkgs"
	CheckError(err)

	_, read_error := os.ReadFile(path + "/milk.pkg")
	if read_error != nil {
		if read_error.Error() == fmt.Sprintf("open %s: no such file or directory", path+"/milk.pkg") {
			errors.NotAPackageError("add")
		}
	}

	for _, url := range urls {
		errors.ThrowInfo(fmt.Sprintf("   TRY -> Downloading URL: \"%s\"", url), false)

		// Clear temporary git directory
		CheckError(os.RemoveAll(mik_pkgs_path + "/git"))
		CheckError(os.MkdirAll(mik_pkgs_path+"/git/", 0777))

		// Clone the repository
		CheckError(exec.Command("git", "clone", url, mik_pkgs_path+"/git").Run())

		// Fetch package
		package_name, entry := Milk(mik_pkgs_path+"/git/", true)

		// Check if package Already exists
		names, _ := GetNamesAndOriginsOfLocallyInstalledPackages()
		if utils.StringInArray(package_name, names) {
			CheckError(os.RemoveAll(mik_pkgs_path + "/git"))
			errors.PackageAlreadyExistsError(package_name)
		}

		// Make package structure
		CheckError(os.MkdirAll(fmt.Sprintf(mik_pkgs_path+"/pkgs/%s", package_name), 0777))
		CheckError(cp.Copy(mik_pkgs_path+"/git/", fmt.Sprintf(mik_pkgs_path+"/pkgs/%s", package_name), cp.Options{AddPermission: 0777}))

		// Preprocess and make a main_package_name.milk
		a, err := os.ReadFile(fmt.Sprintf(mik_pkgs_path+"/pkgs/%s/%s", package_name, entry))
		text := string(a)
		CheckError(err)

		new_path := fmt.Sprintf(mik_pkgs_path+"/pkgs/%s/%s", package_name, entry)
		preprocessed_text := preprocessor.Preprocess(&text, &new_path)
		CheckError(os.WriteFile(fmt.Sprintf(mik_pkgs_path+"/pkgs/%s/main_%s.milk", package_name, package_name), []byte(*preprocessed_text), 0777))

		// Update log file
		b, err2 := os.ReadFile(mik_pkgs_path + "/mip_pkg_log.log")
		file_contents := string(b)
		CheckError(err2)
		CheckError(os.WriteFile(mik_pkgs_path+"/mip_pkg_log.log", []byte(fmt.Sprintf("%s%s:::%s\n", file_contents, package_name, url)), 0777))

		// Clear git folder
		CheckError(os.RemoveAll(mik_pkgs_path + "/git"))

		errors.ThrowSuccess(fmt.Sprintf("Downloaded PKG from URL: \"%s\"\n", url), false)
	}
}
