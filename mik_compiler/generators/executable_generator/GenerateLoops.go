package executablegenerator

import (
	"MicNewDawn/mik_compiler/parser"
	"fmt"
)

func GenerateWhile(node parser.FirstClass) {
	var loop_body_block_position int
	var loop_head_block_position int
	var normal_code_position int

	new_node := node.(parser.WhileNode)
	CodeAdd(fmt.Sprintf("\tbr label %%%d\n\n%d:\n", EMULATED_REGISTER_COUNTER, EMULATED_REGISTER_COUNTER), false) // branch to the loop head
	loop_head_block_position = EMULATED_REGISTER_COUNTER
	EMULATED_REGISTER_COUNTER++

	// evaluate the loop head code
	GenerateExpression(new_node.Boolean_statement, new_node.Type)
	loop_body_block_position = EMULATED_REGISTER_COUNTER

	// get the corresponding labels
	prev_code := *CODE
	prev_main_code := *MAIN_FUNCTION_CODE
	prev_emulated_register_counter := EMULATED_REGISTER_COUNTER
	prev_emulated_register_counter_string := EMULATED_REGISTER_STRING_COUNTER

	LLVMGenerateIR(new_node.CodeBlock)
	normal_code_position = EMULATED_REGISTER_COUNTER + 1

	*CODE = prev_code
	*MAIN_FUNCTION_CODE = prev_main_code
	EMULATED_REGISTER_COUNTER = prev_emulated_register_counter
	EMULATED_REGISTER_STRING_COUNTER = prev_emulated_register_counter_string

	// generate the conditional branch with the corresponding labels
	CodeAdd(fmt.Sprintf("\tbr i1 %%%d, label %%%d, label %%%d\n", EMULATED_REGISTER_COUNTER-1, loop_body_block_position, normal_code_position), false)

	// generate loop body label
	CodeAdd(fmt.Sprintf("%d:\n", EMULATED_REGISTER_COUNTER), false)
	EMULATED_REGISTER_COUNTER++
	LLVMGenerateIR(new_node.CodeBlock)
	CodeAdd(fmt.Sprintf("\tbr label %%%d\n", loop_head_block_position), false)

	// generate normal code label
	CodeAdd(fmt.Sprintf("%d:\n", EMULATED_REGISTER_COUNTER), false)
	EMULATED_REGISTER_COUNTER++
}
