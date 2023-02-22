package executablegenerator

import (
	"MicNewDawn/mik_compiler/parser"
	"fmt"
)

func GenerateIf(node parser.FirstClass) {
	new_node := node.(parser.IfNode)
	GenerateExpression(new_node.Boolean_statement, new_node.Type)
	if_return_value := EMULATED_REGISTER_COUNTER - 1

	var if_block_position int
	var else_block_position int
	var normal_code_position int

	prev_code := *CODE
	prev_main_code := *MAIN_FUNCTION_CODE
	prev_emulated_register_counter := EMULATED_REGISTER_COUNTER
	prev_emulated_register_counter_string := EMULATED_REGISTER_STRING_COUNTER

	if_block_position = EMULATED_REGISTER_COUNTER
	EMULATED_REGISTER_COUNTER++
	LLVMGenerateIR(new_node.CodeBlock)

	if new_node.Else {
		else_block_position = EMULATED_REGISTER_COUNTER
		EMULATED_REGISTER_COUNTER++
		LLVMGenerateIR(new_node.Else_body)
	}
	normal_code_position = EMULATED_REGISTER_COUNTER
	if else_block_position == 0 {
		else_block_position = normal_code_position
	}

	*CODE = prev_code
	*MAIN_FUNCTION_CODE = prev_main_code
	EMULATED_REGISTER_COUNTER = prev_emulated_register_counter
	EMULATED_REGISTER_STRING_COUNTER = prev_emulated_register_counter_string

	CodeAdd(fmt.Sprintf("\tbr i1 %%%d, label %%%d, label %%%d\n", if_return_value, if_block_position, else_block_position), false)

	CodeAdd(fmt.Sprintf("%d:\n", EMULATED_REGISTER_COUNTER), false)
	EMULATED_REGISTER_COUNTER++
	LLVMGenerateIR(new_node.CodeBlock)
	CodeAdd(fmt.Sprintf("\tbr label %%%d\n", normal_code_position), false)

	if new_node.Else {
		CodeAdd(fmt.Sprintf("\n%d:\n", EMULATED_REGISTER_COUNTER), false)
		EMULATED_REGISTER_COUNTER++
		LLVMGenerateIR(new_node.Else_body)
		CodeAdd(fmt.Sprintf("\tbr label %%%d\n", EMULATED_REGISTER_COUNTER), false)
	}

	CodeAdd(fmt.Sprintf("%d:\n", EMULATED_REGISTER_COUNTER), false)
	EMULATED_REGISTER_COUNTER++
}
