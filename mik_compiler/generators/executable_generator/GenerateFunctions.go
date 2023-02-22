package executablegenerator

import (
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"fmt"
	"strings"
)

func GenerateFunctionDeclaration(node parser.FunctionNode) {
	// declare <type><pointers> @<function_name>(<arg>, <arg>)
	var argument_text string

	for idx, value := range node.Arguments.TypeArray {
		argument_text += GetTypeLLVMEncoding(value)
		argument_text += strings.Repeat("*", int(value.Dimension))

		if idx < len(node.Arguments.TypeArray)-1 {
			argument_text += ", "
		}
	}

	CodeAdd(fmt.Sprintf(
		"declare %s @%s(%s)\n",
		GetTypeLLVMEncoding(node.Return_type),
		node.Function_name,
		argument_text,
	), true)
}

func GenerateFunctionDefinition(node parser.FunctionNode) {
	old_counter1 := EMULATED_REGISTER_COUNTER
	old_counter2 := EMULATED_REGISTER_STRING_COUNTER
	EMULATED_REGISTER_COUNTER = 0
	EMULATED_REGISTER_STRING_COUNTER = 1

	In_function = true
	// define <type><pointers> @<function_name>(<arg>, <arg>) {<code_body>}

	CodeAdd(fmt.Sprintf(
		"define %s @%s(",
		GetTypeLLVMEncoding(node.Return_type),
		node.Function_name,
	), true)
	GetFunctionArgumentLLVMSignature(node.Arguments)

	old_VARIABLES := parser.NamesAndScopes.VARIABLES
	parser.NamesAndScopes.VARIABLES = node.VARIABLES
	LLVMGenerateIR(node.CodeBlock)

	if node.Return_type.BaseType == utils.T_VOID {
		CodeAdd("\n\tret void", false)
	}

	parser.NamesAndScopes.VARIABLES = old_VARIABLES
	EMULATED_REGISTER_COUNTER = old_counter1
	EMULATED_REGISTER_STRING_COUNTER = old_counter2

	CodeAdd("\n}\n", false)
	In_function = false
}

func GenerateFunctions(node parser.FirstClass) {
	new_node := node.(parser.FunctionNode)
	// check for unused and remove if unused flag exists
	if (remove_unused && !utils.StringInArray("volatile", new_node.Variable_behavior_descriptors)) && parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Function_name] < 1 {
		return
	} else if utils.StringInArray("volatile", new_node.Variable_behavior_descriptors) {
		parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Function_name]++
	}

	if new_node.Declared {
		GenerateFunctionDeclaration(new_node)
	} else {
		GenerateFunctionDefinition(new_node)
	}
}
