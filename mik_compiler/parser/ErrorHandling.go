package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/utils"
	"reflect"
)

// Throws error if types are not compatible
func TypeCompare(compare_type utils.Type) {
	// TODO: compare custom types

	// set type to the current type, as all types after that have to be equal even in a compare
	if !Flags.compare_types {
		ControlVariables.current_expected_type = compare_type
	}

	// filter string
	if ControlVariables.current_expected_type.BaseType == utils.T_STRING && compare_type.BaseType != utils.T_STRING {
		errors.TypeMismatchError(ControlVariables.current_expected_type.TypeName, compare_type.TypeName, GetLocation())
	}

	// check if dimensions of array check out
	if compare_type.Dimension != ControlVariables.current_expected_type.Dimension {
		errors.TypeMismatchError(ControlVariables.current_expected_type.TypeName, compare_type.TypeName, GetLocation())
	}

	// check wether both types are some sort of ints. If so, both types compare just fine, as implicit type conversion will take care of truncating or extending
	if compare_type.BaseType >= utils.T_CHAR && compare_type.BaseType <= utils.T_INT64 && ControlVariables.current_expected_type.BaseType >= utils.T_CHAR && ControlVariables.current_expected_type.BaseType <= utils.T_INT64 {
		return
	}

	if utils.TYPE_HIERARCHY[ControlVariables.current_expected_type.BaseType] < utils.TYPE_HIERARCHY[compare_type.BaseType] /*&& Flags.compare_types*/ {
		// booleans are handled differently
		if !Flags.in_boolean_expression {
			errors.TypeMismatchError(ControlVariables.current_expected_type.TypeName, compare_type.TypeName, GetLocation())
		}
	}

	// update the most significant type
	if ControlVariables.most_significant_type.BaseType == utils.T_INVALID {
		ControlVariables.most_significant_type = compare_type
	} else {
		if utils.TYPE_HIERARCHY[compare_type.BaseType] > utils.TYPE_HIERARCHY[ControlVariables.most_significant_type.BaseType] {
			ControlVariables.most_significant_type = compare_type
		}
	}
}

// checks wether a given name is taken and throws an appropriate error
func CheckNameTaken(a string) string {
	if utils.StringInArray(a, NamesAndScopes.Variable_names) ||
		utils.StringInArray(a, NamesAndScopes.Function_names) ||
		utils.StringInArray(a, NamesAndScopes.Struct_names) {
		errors.NameAlreadyTakenError(GetLocation(), a)
	} else {
		return a
	}

	return a
}

func VariableBehaviorDescriptorValidityCheck(list []string) []string {
	for _, value := range list {
		// if conflicts
		if utils.StringInArray(utils.VARIABLE_BEHAVIOR_DESCRIPTORS_CONFLICT_MAP[value], list) || utils.StringInArray(ReverseStringMap(utils.VARIABLE_BEHAVIOR_DESCRIPTORS_CONFLICT_MAP)[value], list) {
			// prints conflicting error. The long second argument is the value the map at index + inverse map at index. One of them will always be "". Hence this will work
			errors.ConflictingTypesError(GetLocation(), value, utils.VARIABLE_BEHAVIOR_DESCRIPTORS_CONFLICT_MAP[value]+ReverseStringMap(utils.VARIABLE_BEHAVIOR_DESCRIPTORS_CONFLICT_MAP)[value])
		}
	}

	return list
}

func TypeDescriptorValidityCheck(list []string) []string {
	for _, value := range list {
		// if conflicts
		if utils.StringInArray(utils.TYPE_DESCRIPTORS_CONFLICT_MAP[value], list) || utils.StringInArray(ReverseStringMap(utils.TYPE_DESCRIPTORS_CONFLICT_MAP)[value], list) {
			// prints conflicting error
			errors.ConflictingTypesError(GetLocation(), value, utils.TYPE_DESCRIPTORS_CONFLICT_MAP[value]+ReverseStringMap(utils.TYPE_DESCRIPTORS_CONFLICT_MAP)[value])
		}
	}

	return list
}

// TODO: needs work
func ScopeValidityCheck(variable FirstClass, prefix string) {
	switch reflect.TypeOf(variable).Name() {
	case "AssignmentNode":
		{
			assignment := variable.(AssignmentNode)

			if utils.StringInArray("local", assignment.Variable_behavior_descriptors) && (assignment.Section != ControlVariables.current_section && !ItemInStringMap(prefix, NamesAndScopes.currently_used_sections)) { // is in another scope
				errors.LocalScopeError(GetLocation(), assignment.Variable_name, assignment.Section, ControlVariables.current_section)
			}
		}
	case "ReAssignmentNode":
		{
			reassignment := variable.(ReAssignmentNode)
			assignment := NamesAndScopes.VARIABLES[reassignment.Variable_name].(AssignmentNode)

			if utils.StringInArray("local", assignment.Variable_behavior_descriptors) && (assignment.Section != ControlVariables.current_section && !ItemInStringMap(prefix, NamesAndScopes.currently_used_sections)) { // is in another scope
				errors.LocalScopeError(GetLocation(), assignment.Variable_name, assignment.Section, ControlVariables.current_section)
			}
		}
	case "FunctionNode":
		{
			function := variable.(FunctionNode)

			if utils.StringInArray("local", function.Variable_behavior_descriptors) && (function.Section != ControlVariables.current_section && !ItemInStringMap(prefix, NamesAndScopes.currently_used_sections)) {
				ControlVariables.token_causing_error = function.Function_name
				errors.LocalScopeError(GetLocation(), function.Function_name, function.Section, ControlVariables.current_section)
			}
		}
	}
}