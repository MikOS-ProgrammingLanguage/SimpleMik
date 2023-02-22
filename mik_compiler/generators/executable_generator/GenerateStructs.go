package executablegenerator

import (
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"fmt"
	"reflect"
)

func GenerateStruct(first_class parser.FirstClass) {
	var packed_start string = ""
	var packed_end string = ""
	var type_string string
	struct_node := first_class.(parser.StructNode)

	// check if packed and change rune accordingly
	if utils.Int8InArray(utils.PACKED, struct_node.Attributes) {
		packed_start = "<"
		packed_end = ">"
	}

	// evaluate type string
	for _, val := range struct_node.VARIABLES {
		// can only be assignment or array assignment
		switch reflect.TypeOf(val).Name() {
		case "AssignmentNode":
			{
				node := val.(parser.AssignmentNode)
				type_string += GetTypeLLVMEncoding(node.Type)
			}
		case "ArrayAssignmentNode":
			{
				// needs fixed size
				//node := val.(parser.ArrayAssignmentNode)
			}
		}

		type_string += ", "
	}
	if len(type_string) == 0 {
		type_string = ", "
	}

	CodeAdd(
		fmt.Sprintf("%%struct.%s = type %s{ %s }%s\n", struct_node.Name, packed_start, type_string[:len(type_string)-2], packed_end),
		true,
	)
}

func GenerateEStruct() {

}
