package executablegenerator

import (
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"fmt"
	"reflect"
	"strings"
)

func GetFunctionCallBody(node []parser.SecondClass, name string) string {
	parent_function_arguments := parser.NamesAndScopes.FUNCTIONS[name].(parser.FunctionNode).Arguments

	var return_string string

	for idx, i := range node {
		var a ExpressionReturn_t
		if len(parent_function_arguments.TypeArray) == 1 {
			return_string += GetTypeLLVMEncoding(parent_function_arguments.TypeArray[0]) + " "
			a = GenerateExpression(i, parent_function_arguments.TypeArray[0])
		} else {
			return_string += GetTypeLLVMEncoding(parent_function_arguments.TypeArray[idx]) + " "
			a = GenerateExpression(i, parent_function_arguments.TypeArray[idx])
		}
		return_string += a.return_value

		if idx < len(node)-1 {
			return_string += ", "
		}
	}

	return return_string
}

// get the correct typecast opcode for given type
func GetTypeCastLLVMEncoding(given_type, to_type string) string {
	switch given_type {
	case "double":
		{
			return "fptosi"
		}
	}

	// given type is smaller than to type sext otherwise trunc
	if utils.TYPE_SIZES[given_type] < utils.TYPE_SIZES[to_type] {
		return "sext"
	} else {
		return "trunc"
	}
}

func GetComparePrefix(a string) string {
	if a != "f" {
		a = "i"
	}
	return a
}

// Get the corresponding instruction for a given operator
func GetOperatorLLVMEncoding(token string, Type utils.Type) string {
	var prefix string
	switch Type.BaseType {
	case utils.T_FLOAT:
		{
			prefix = "f"
		}
	}

	switch token {
	case "+":
		{
			return prefix + "add"
		}
	case "-":
		{
			return prefix + "sub"
		}
	case "*":
		{
			return prefix + "mul"
		}
	case "/":
		{
			return prefix + "div"
		}
	case "==":
		{
			return GetComparePrefix(prefix) + "cmp eq"
		}
	case "!=":
		{
			return GetComparePrefix(prefix) + "cmp ne"
		}
	case ">=":
		{
			return GetComparePrefix(prefix) + "cmp sge"
		}
	case "<=":
		{
			return GetComparePrefix(prefix) + "cmp sle"
		}
	case ">":
		{
			return GetComparePrefix(prefix) + "cmp sgt"
		}
	case "<":
		{
			return GetComparePrefix(prefix) + "cmp slt"
		}
	case "&&":
		{
			return "and"
		}
	case "||":
		{
			return "or"
		}
	default:
		{
			// error. Not an operator
		}
	}

	return ""
}

func GetDefaultValuesForTypes(Type utils.Type) string {
	switch Type.BaseType {
	case utils.T_FLOAT:
		{
			return "0.0"
		}
	default:
		{
			return "0"
		}
	}
}

func GetTypeLLVMEncoding(Type utils.Type) string {
	switch Type.BaseType {
	case utils.T_INT:
		{
			return "i32"
		}
	case utils.T_INT64:
		{
			return "i64"
		}
	case utils.T_INT16:
		{
			return "i16"
		}
	case utils.T_INT8:
		{
			return "i8"
		}
	case utils.T_UINT64:
		{
			return "i64"
		}
	case utils.T_UINT32:
		{
			return "i32"
		}
	case utils.T_UINT16:
		{
			return "i16"
		}
	case utils.T_UINT8:
		{
			return "i8"
		}
	case utils.T_CHAR:
		{
			return "i8"
		}
	case utils.T_FLOAT:
		{
			return "double"
		}
	case utils.T_STRING:
		{
			return "i8*"
		}
	case utils.T_COCK:
		{
			return "i64"
		}
	case utils.T_BOOL:
		{
			return "i1"
		}
	case utils.T_VOID:
		{
			return "void"
		}
	default:
		{
			return Type.TypeName
		}
	}
}

func GetFunctionArgumentLLVMSignature(node parser.FunctionAssignArguments) {
	var return_string string
	var return_string2 string
	var offset int = len(node.Arguments)

	EMULATED_REGISTER_COUNTER = offset
	for idx, i := range node.Arguments {

		switch reflect.TypeOf(i).Name() {
		case "AssignmentNode":
			{
				// <type><pointers> %<count>
				new_node := i.(parser.AssignmentNode)
				
				return_string += fmt.Sprintf(
					"%s %%%d",
					GetTypeLLVMEncoding(new_node.Type),
					idx,
				)

				return_string2 += fmt.Sprintf(
					"\t%%%s = alloca %s, align 4\n",
					new_node.Variable_name,
					GetTypeLLVMEncoding(new_node.Type),
				)

				return_string2 += fmt.Sprintf(
					"\tstore %s %%%d, %s* %%%s, align 4\n",
					GetTypeLLVMEncoding(new_node.Type),
					idx,
					GetTypeLLVMEncoding(new_node.Type),
					new_node.Variable_name,
				)
				VARIABLES = append(VARIABLES, new_node.Variable_name)
			}
		case "ArrayAssignmentNode":
			{
				// <type><pointers> %<count>
				new_node := i.(parser.ArrayAssignmentNode)

				pointers_needed := strings.Repeat("*", int(node.TypeArray[idx].Dimension))

				return_string += fmt.Sprintf(
					"%s%s %%%d",
					GetTypeLLVMEncoding(new_node.Type),
					pointers_needed,
					idx,
				)
			}
		default:
			{
				// error
			}
		}

		if idx < len(node.Arguments)-1 {
			return_string += ", "
		}
	}

	CodeAdd(return_string+") {\n"+return_string2, false)
	EMULATED_REGISTER_COUNTER++
}
