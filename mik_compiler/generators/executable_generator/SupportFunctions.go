package executablegenerator

import (
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"fmt"
)

func CodeAdd(a string, b bool) {
	if In_function || b {
		*CODE += a
	} else {
		*MAIN_FUNCTION_CODE += a
	}
}

func GenerateArrayDimensionString(lens []parser.SecondClass, node_type utils.Type) string {
	var dimensions string
	var cnt int = 0
	for _, i := range lens {
		length_expression_return := GenerateExpression(i, node_type)
		var length string
		if length_expression_return.store_expression_return {
			length = length_expression_return.return_value
		} else {
			length = fmt.Sprintf("%%%d", EMULATED_REGISTER_COUNTER-1)
		}

		dimensions += fmt.Sprintf(
			"[%s x ",
			length,
		)
		cnt++
	}
	dimensions += fmt.Sprint(GetTypeLLVMEncoding(node_type))

	for i := 0; i < cnt; i++ {
		dimensions += "]"
	}

	return dimensions
}