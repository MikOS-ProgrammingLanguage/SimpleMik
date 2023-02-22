package executablegenerator

import (
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"fmt"
)

func GenerateReassignment(node parser.FirstClass) {
	new_node := node.(parser.ReAssignmentNode)
	// check for unused and remove if unused flag exists
	if remove_unused && parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Variable_name] < 1 {
		return
	}

	node_type := new_node.Type

	var VariableNamePrefix string = "%"

	if new_node.Is_global {
		VariableNamePrefix = "@"
	}


	// generate expression
	expression_return := GenerateExpression(new_node.Body, node_type)

	if node_type.BaseType == utils.T_STRING && expression_return.return_value_is_a_literal {
		CodeAdd(fmt.Sprintf("\tstore i8* %s, i8** %%%s, align 4\n", expression_return.return_value, new_node.Variable_name), false)
	} else if !expression_return.store_expression_return {
		value_to_store_register_count := EMULATED_REGISTER_COUNTER - 1

		// Store the value
		// Equivalent to a simple llvm-ir store: store type %4, type* %name, align 4
		CodeAdd(fmt.Sprintf(
			"\tstore %s %%%d, %s* %s%s, align 4\n",
			GetTypeLLVMEncoding(node_type),
			value_to_store_register_count,
			GetTypeLLVMEncoding(node_type),
			VariableNamePrefix,
			new_node.Variable_name), false)
	} else {
		CodeAdd(fmt.Sprintf(
			"\tstore %s %s, %s* %s%s, align 4\n",
			GetTypeLLVMEncoding(node_type),
			expression_return.return_value,
			GetTypeLLVMEncoding(node_type),
			VariableNamePrefix,
			new_node.Variable_name), false)
	}
}

// buggy as fuck
func GenerateArraySliceReassignment(node parser.FirstClass) {
	new_node := node.(parser.ArraySliceReassignmentNode)
	// check for unused and remove if unused flag exists
	if remove_unused && parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Variable_name] < 1 {
		return
	}

	node_type := new_node.Type
	var VariableNamePrefix string = "%"

	if new_node.Is_global {
		VariableNamePrefix = "@"
	}
	var_name := new_node.Variable_name

	// Evaluate array length expression

	for idx, i := range new_node.Indices {

		dimensions := GenerateArrayDimensionString(new_node.Lengths[:len(new_node.Lengths)-idx], node_type)
		// rework to dereference multiple dimensions
		index_expression_return := GenerateExpression(i, node_type)
		var index string
		if index_expression_return.store_expression_return {
			index = index_expression_return.return_value
		} else {
			index = fmt.Sprintf("%%%d", EMULATED_REGISTER_COUNTER-1)
		}

		// %1 = get element ptr inbounds. (load the array at index)
		CodeAdd(fmt.Sprintf(
			"\t%%%d = getelementptr inbounds %s, %s* %s%s, i64 0, i64 %s\n",
			EMULATED_REGISTER_COUNTER,
			dimensions,
			dimensions,
			VariableNamePrefix,
			var_name,
			index,
		), false)
		EMULATED_REGISTER_COUNTER++
		VariableNamePrefix = "%"
		var_name = fmt.Sprint(EMULATED_REGISTER_COUNTER-1)
	}
	name := fmt.Sprint(EMULATED_REGISTER_COUNTER-1)

	// store a value at that index (%1)
	// generate expression
	expression_return := GenerateExpression(new_node.Body, node_type)

	if expression_return.return_value_is_a_literal {
		CodeAdd(fmt.Sprintf(
			"\tstore %s %s, %s* %%%s, align 4\n",
			GetTypeLLVMEncoding(node_type),
			expression_return.return_value,
			GetTypeLLVMEncoding(node_type),
			name), false)
	} else if !expression_return.store_expression_return {
		CodeAdd(fmt.Sprintf(
			"\tstore %s %%%d, %s* %%%s, align 4\n",
			GetTypeLLVMEncoding(node_type),
			EMULATED_REGISTER_COUNTER-1,
			GetTypeLLVMEncoding(node_type),
			name), false)
	}
}

func GenerateArrayAssignment(node parser.FirstClass) {
	in_array = true
	new_node := node.(parser.ArrayAssignmentNode)
	// check for unused and remove if unused flag exists
	if (remove_unused && !utils.StringInArray("volatile", new_node.Variable_behavior_descriptors)) && parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Variable_name] < 1 {
		return
	} else if utils.StringInArray("volatile", new_node.Variable_behavior_descriptors) {
		parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Variable_name]++
	}

	node_type := new_node.Type

	// Allocation
	var VariableNamePrefix string

	if new_node.Is_global {
		VariableNamePrefix = "@"
	} else {
		VariableNamePrefix = "%"

		// Evaluate array length expression
		dimensions := GenerateArrayDimensionString(new_node.Array_length, node_type)

		CodeAdd(fmt.Sprintf(
			"\t%s%s = alloca %s, align 4\n",
			VariableNamePrefix,
			new_node.Variable_name,
			dimensions,
		), false)
	}

	VARIABLES = append(VARIABLES, new_node.Variable_name)
	in_array = false
}

func GenerateAssignment(node parser.FirstClass) {
	new_node := node.(parser.AssignmentNode)
	// check for unused and remove if unused flag exists
	if (remove_unused && !utils.StringInArray("volatile", new_node.Variable_behavior_descriptors)) && parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Variable_name] < 1 {
		return
	} else if utils.StringInArray("volatile", new_node.Variable_behavior_descriptors) {
		parser.NamesAndScopes.VARIABLE_REFERENCE_COUNT[new_node.Variable_name]++
	}

	node_type := new_node.Type

	// Allocation
	var VariableNamePrefix string

	if new_node.Is_global {
		VariableNamePrefix = "@"
		DefaultValue := GetDefaultValuesForTypes(node_type)

		// Is a simple global llvm assignment: @name = global <type> <value>(0 by default), align 4
		CodeAdd(fmt.Sprintf(
			"@%s = common global %s %s, align 4\n",
			new_node.Variable_name,
			GetTypeLLVMEncoding(node_type),
			DefaultValue), true)

		// make a reassign to actually assign it to a value
		if new_node.Variable_body != nil {
			reassign_node := parser.ReAssignmentNode{Type: new_node.Type,
				Is_global:          new_node.Is_global,
				Variable_name:      new_node.Variable_name,
				Body:               new_node.Variable_body,
				Section:            new_node.Section}

			VARIABLES = append(VARIABLES, new_node.Variable_name)
			GenerateReassignment(reassign_node)
		} else {
			return
		}
	} else {
		VariableNamePrefix = "%"

		// Is a simple llvm-ir allocation: %name = alloca type, align 4
		CodeAdd(fmt.Sprintf(
			"\t%s%s = alloca %s, align 4\n",
			VariableNamePrefix,
			new_node.Variable_name,
			GetTypeLLVMEncoding(node_type)), false)
		VARIABLES = append(VARIABLES, new_node.Variable_name)

		// Body
		if new_node.Variable_body != nil {
			expression_return := GenerateExpression(new_node.Variable_body, node_type)
		
			if node_type.BaseType == utils.T_STRING && expression_return.return_value_is_a_literal {
				CodeAdd(fmt.Sprintf("\tstore i8* %s, i8** %%%s, align 4\n", expression_return.return_value, new_node.Variable_name), false)
			} else if !expression_return.store_expression_return {
				// Store the value
				// Equivalent to a simple llvm-ir store: store type %4, type* %name, align 4
				CodeAdd(fmt.Sprintf(
					"\tstore %s %%%d, %s* %s%s, align 4\n",
					GetTypeLLVMEncoding(node_type),
					EMULATED_REGISTER_COUNTER-1,
					GetTypeLLVMEncoding(node_type),
					VariableNamePrefix,
					new_node.Variable_name), false)
			} else {
				CodeAdd(fmt.Sprintf(
					"\tstore %s %s, %s* %s%s, align 4\n",
					GetTypeLLVMEncoding(node_type),
					expression_return.return_value,
					GetTypeLLVMEncoding(node_type),
					VariableNamePrefix,
					new_node.Variable_name), false)
			}
		}
	}
}

func GenerateReturn(node parser.FirstClass) {
	new_node := node.(parser.ReturnNode)

	if new_node.Return_expression == nil {
		CodeAdd("\tret void", false)
		return
	}

	expression_return := GenerateExpression(new_node.Return_expression, new_node.Goal_type)
	if expression_return.return_value_is_a_reference { // trim the expression return if it's a reference
		expression_return.return_value = expression_return.return_value[1:]
	}
	CodeAdd(fmt.Sprintf(
		"\tret %s %s",
		GetTypeLLVMEncoding(new_node.Goal_type),
		expression_return.return_value,
	), false)
}
