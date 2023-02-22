package preprocessor

import (
	"MicNewDawn/errors"
	"MicNewDawn/utils"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var Included_files, Included_mip_files []string
var Include_regular_expression *regexp.Regexp = regexp.MustCompile("^((#yoink){1}|(#yoink-pkg){1}){1}( )*(<){1}(.)+(>){1}$")
var Include_normal_regex *regexp.Regexp = regexp.MustCompile("^(#yoink){1}( )*(<){1}(.)+(>){1}$")
var Include_mip_regex *regexp.Regexp = regexp.MustCompile("^(#yoink-pkg){1}( )*(<){1}(.)+(>){1}$")

func Preprocess(text, path *string) *string {
	// create the file directive for the root file
	output_text := fmt.Sprintf("@file(\"%s\")\n", *path)
	*path = strings.ReplaceAll(*path, "\\", "/")
	temporary_path := strings.Split(*path, "/")
	*path = ""

	for index, value := range temporary_path {
		if index < len(temporary_path)-1 {
			*path += string(value) + "/"
		}
	}

	*path = strings.ReplaceAll(*path, "\n", "")

	// get the path to the mip package files
	mip_src_path, err := os.ReadFile("/etc/.mik.conf")
	if err != nil {
		panic(err)
	}
	mip_src_path = []byte(strings.ReplaceAll(string(mip_src_path), "\n", ""))

	// preprocess preprocessor directives
	for _, value := range strings.Split(*text, "\n") {
		// checks if the signature is something like "#yoink <hello.milk>"
		if Include_regular_expression.MatchString(value) {
			name_of_current_include := ""

			// get the name that should be included by splitting like this "#yoink <| file_name |>"
			name_of_current_include = strings.Split(strings.Split(value, "<")[1], ">")[0]
			// ignore if Already included
			if utils.StringInArray(name_of_current_include, Included_files) {
				continue
			}

			// validate what type of include is analyzed
			if Include_normal_regex.MatchString(value) {
				// specify a file directive for the lexer to 
				output_text += fmt.Sprintf("@file(\"%s\")\n", name_of_current_include)

				included_file_contents, err := os.ReadFile(*path + name_of_current_include)
				if err != nil {
					errors.YoinkError(name_of_current_include)
				}

				output_text += string(included_file_contents)
				Included_files = append(Included_files, name_of_current_include)
			} else {
				// add a section lexer directive and append append the contents of the included file
				output_text += fmt.Sprintf("@file(\"%s\")\n", name_of_current_include)

				included_file_contents, err := os.ReadFile(string(mip_src_path) + "mik-src/pkg/" + name_of_current_include + "/main_" + name_of_current_include + ".milk")
				if err != nil {
					errors.YoinkPackageError(name_of_current_include)
				}

				output_text += string(included_file_contents)
				Included_mip_files = append(Included_mip_files, name_of_current_include)
			}
		} else if strings.HasPrefix(value, "#") {
			// Illegal preprocessor directive
			errors.ThrowError(fmt.Sprintf("IllegalPreprocessorDirective. Found illegal preprocessor directive: %s", value), true)
		} else {
			output_text += value + "\n"
		}
	}
	return &output_text
}
