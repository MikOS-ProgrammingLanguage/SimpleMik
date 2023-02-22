package utils

import (
	"MicNewDawn/errors"
	"os"
	"os/exec"
)

// Installs a config file into /etc/
func Install_config() {
	// make mip package structure
	home, home_directory_error := os.UserHomeDir()
	if home_directory_error != nil {
		panic(home_directory_error)
	}

	ch_dir_home_error := os.Chdir(home)
	if ch_dir_home_error != nil {
		panic(ch_dir_home_error)
	}

	make_mik_dir := os.Mkdir(".mik", 0777)
	if make_mik_dir != nil {
		panic(make_mik_dir)
	}

	make_mip_package_dirs := os.MkdirAll("./.mik/mip/pkgs", 0777)
	if make_mip_package_dirs != nil {
		panic(make_mip_package_dirs)
	}

	make_mip_package_log_file := os.WriteFile("./.mik/mip_pkg_log.log", []byte(""), 0777)
	if make_mip_package_log_file != nil {
		panic(make_mip_package_log_file)
	}

	errors.ThrowInfo("Setup Mip directory structure... Checking dependencies!", false)

	// Check for dependencies
	clang_call_output, _ := exec.Command("which", "clang").Output()

	if len(clang_call_output) == 0 {
		errors.ThrowWarning("clang required!", false)
	}

	errors.ThrowSuccess("Installed Successfully!", false)
}
