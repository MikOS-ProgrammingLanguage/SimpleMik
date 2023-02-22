package mip_functions

import "fmt"

// Prints a string with help as well as usefull information
func PrintHelp() {
	fmt.Println("MIP Help")

	fmt.Println("\n\tArguments and usage:")
	fmt.Println("\t\tinstall:\n\t\t\tIs used to install packages. Usage: \"mip install <url> <url> ...\"")

	fmt.Println("\n\t\tupdate:\n\t\t\tIs used to update all installed packages. Usage: \"mip update\"")

	fmt.Println("\n\t\tremove:\n\t\t\tIs used to remove packages. Usage: \"mip remove <package-name> <package-name> ...\"")

	fmt.Println("\n\t\tinit:\n\t\t\tIs used to initialize a project with a given path. Usage: \"mip init (mik|milk) <path>\"")

	fmt.Println("\n\t\tmilk:\n\t\t\tInterprets the milk.pkg file if found. Usage: \"mip milk\"")

	fmt.Println("\n\t\tlist:\n\t\t\tLists all packages installed on the system and in the current project. Usage: \"mip list\"")

	fmt.Println("\n\t\tadd:\n\t\t\tInstalles packages in the current project. Usage: \"mip add <url> <url> ...\"")
}
