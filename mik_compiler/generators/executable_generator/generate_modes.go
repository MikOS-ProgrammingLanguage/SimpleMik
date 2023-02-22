package executablegenerator

import (
	"MicNewDawn/errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GenerateBinary(output_name string) {
	e := os.WriteFile("ir.ll", []byte(string(DEFAULT_CONFIG)+*CODE+"\ndefine i32 @main() {\n"+*MAIN_FUNCTION_CODE+"\n\tret i32 0\n}"), 0667)
	if e != nil {
		panic(e)
	}

	cmd := exec.Command("clang", "ir.ll", "-o", output_name)
	message, exec_error := cmd.CombinedOutput()
	if exec_error != nil {
		s := string(message)
		s2 := strings.Split(s, "\n")
		var res []string = s2
		if strings.HasPrefix(s2[0], "warning: overriding the module") {
			res = s2[2:]
		}

		var ret string
		for _, val := range res {
			ret += "\t" + val + "\n"
		}

		errors.ThrowError("An error occurred when compiling LL representation", false)
		fmt.Println(ret)
		os.Exit(1)
	}

	remove_error := os.Remove("ir.ll")
	if remove_error != nil {
		panic(remove_error.Error())
	}
}

func GenerateLLVM(output_name string) {
	e := os.WriteFile(output_name, []byte(string(DEFAULT_CONFIG)+*CODE+"\ndefine i32 @main() {\n"+*MAIN_FUNCTION_CODE+"\n\tret i32 0\n}"), 0677)
	if e != nil {
		panic(e)
	}
}

func GenerateObject(output_name string) {
	e := os.WriteFile("ir.ll", []byte(string(DEFAULT_CONFIG)+*CODE+"\ndefine i32 @main() {\n"+*MAIN_FUNCTION_CODE+"\n\tret i32 0\n}"), 0667)
	if e != nil {
		panic(e)
	}

	exec_error := exec.Command("clang", "ir.ll", "-c", "-o", output_name).Run()
	if exec_error != nil {
		panic(exec_error.Error())
	}

	remove_error := os.Remove("ir.ll")
	if remove_error != nil {
		panic(remove_error.Error())
	}
}
