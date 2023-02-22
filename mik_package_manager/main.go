package main

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_package_manager/mip_functions"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		errors.ToFewArgumentError()
	}

	/*
		Arguments:
			install:
				installs a package to the system mik directory via github url

			update:
				updates all packages

			remove:
				removes any packages specified from the system mik directory

			init:
				initializes a new mik project

			milk:
				"milks" a milk.pkg file. Basically installing dependencies removing certain files and so on

			list:
				lists all packages installed

			add:
				add a package to the mik_packages folder in your project

			help:
				print a string with usefull information as well as help
	*/
	switch os.Args[1] {
	case "install":
		{
			mip_functions.InstallPackage(os.Args[2:])
		}
	case "update":
		{

		}
	case "remove":
		{
			mip_functions.RemovePackage(os.Args[2:])
		}
	case "init":
		{
			mip_functions.Init(os.Args[2:])
		}
	case "milk":
		{
			mip_functions.Milk(os.Args[2], false)
		}
	case "list":
		{
			mip_functions.ListPackages()
		}
	case "add":
		{
			mip_functions.Add(os.Args[2:])
		}
	case "help":
		{
			mip_functions.PrintHelp()
		}
	default:
		{
			errors.UnknownArgumentError(os.Args[1])
		}
	}
}
