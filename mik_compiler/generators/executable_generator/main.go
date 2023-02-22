package executablegenerator

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"os/exec"
	"reflect"
	"strings"
)

var DEFAULT_CONFIG /*[]byte*/ string

var VARIABLES []string

var in_array bool = false

var code_default string = ";built-ins\ndeclare i32 @puts(i8*)\n\n"
var null_pointer_dodge2 string = ""

var EMULATED_REGISTER_COUNTER int = 1
var EMULATED_REGISTER_STRING_COUNTER int = 1
var CODE *string = &code_default
var MAIN_FUNCTION_CODE *string = &null_pointer_dodge2

var In_function bool = false
var remove_unused bool = false

func LLVMGenerateIR(nodes []parser.FirstClass) {
	for _, value := range nodes {
		switch reflect.TypeOf(value).Name() {
		case "AssignmentNode":
			{
				GenerateAssignment(value)
			}
		case "ArrayAssignmentNode":
			{
				GenerateArrayAssignment(value)
			}
		case "ReAssignmentNode":
			{
				GenerateReassignment(value)
			}
		case "ArraySliceReassignmentNode":
			{
				GenerateArraySliceReassignment(value)
			}
		case "FunctionNode":
			{
				GenerateFunctions(value)
			}
		case "ReturnNode":
			{
				GenerateReturn(value)
			}
		case "FirstClassFunctionCall":
			{
				GenerateExpression(value.(parser.FirstClassFunctionCall).Node, utils.VoidTypeConstr())
			}
		case "IfNode":
			{
				GenerateIf(value)
			}
		case "WhileNode":
			{
				GenerateWhile(value)
			}
		case "StructNode":
			{
				GenerateStruct(value)
			}
		}
	}
}

func GenerateExecutableOrIntermediateRepresentation(AST *parser.RootNode, output_name, mode string, ignore_unused, unused_warning bool) {
	remove_unused = ignore_unused
	clang_call_output, clang_error := exec.Command("clang", "--help").Output()
	if clang_error != nil {
		panic(clang_error.Error())
	}

	if !strings.HasSuffix(string(clang_call_output), "linker\n") {
		errors.ThrowError("Clang Toolchain missing! Please install it", true)
	}

	/*var err error
	DEFAULT_CONFIG, err = exec.Command("./get_name").Output()
	if err != nil {
		panic(err)
	}*/
	// DEFAULT_CONFIG = "target datalayout = \"e-m:o-i64:64-i128:128-n32:64-S128\"\ntarget triple = \"arm64-apple-macosx12.0.0\"\n\n" // Change later

	LLVMGenerateIR(AST.Nodes)
	if mode == "bin" {
		GenerateBinary(output_name)
	} else if mode == "llvm" {
		GenerateLLVM(output_name)
	} else {
		GenerateObject(output_name)
	}
}
