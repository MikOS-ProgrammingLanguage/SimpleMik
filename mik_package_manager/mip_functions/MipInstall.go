package mip_functions

import (
	"MicNewDawn/errors"
	"fmt"
	"os"
	"os/exec"
)

func InstallPackage(urls []string) {
	// Find home directory and cd into it
	home, home_directory_error := os.UserHomeDir()
	CheckError(home_directory_error)

	CheckError(os.Chdir(home))
	// Cd into .mik directory
	CheckError(os.Chdir(".mik"))

	for _, url := range urls {
		errors.ThrowInfo(fmt.Sprintf("   TRY -> Downloading URL: \"%s\"", url), false)

		// Clear temporary git directory
		CheckError(os.RemoveAll("./mip/git"))
		CheckError(os.MkdirAll("./mip/git/", 0777))

		// Clone the repository
		CheckError(exec.Command("git", "clone", url, "./mip/git").Run())

		// Fetch package
		FetchPackage("./mip/git/", url)

		// Clear git folder
		CheckError(os.RemoveAll("./mip/git"))

		errors.ThrowSuccess(fmt.Sprintf("Downloaded PKG from URL: \"%s\"\n", url), false)
	}
}
