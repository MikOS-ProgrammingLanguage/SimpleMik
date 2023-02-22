package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
	"fmt"
	"reflect"
)

func Reassignment() []FirstClass {
	prev_pos := ControlVariables.position
	prev_tok := ControlVariables.current_token
	variable_name, prefix := GetFullName()

	if utils.StringInArray(variable_name.Value, NamesAndScopes.Function_names) {
		referenced_function := NamesAndScopes.FUNCTIONS[variable_name.Value]
		ScopeValidityCheck(referenced_function, prefix)
		ControlVariables.current_token = prev_tok
		ControlVariables.position = prev_pos

		reference := referenced_function.(FunctionNode)
		if reference.Return_type.BaseType == utils.T_VOID || utils.StringInArray("__noret__", reference.Variable_behavior_descriptors) {
			if reference.Return_type.BaseType == utils.T_VOID {
				ControlVariables.current_expected_type = utils.VoidTypeConstr()
			}
			return []FirstClass{FirstClassFunctionCall{Node: Factor(true)}}
		}
	}
	ScopeValidityCheck(NamesAndScopes.VARIABLES[variable_name.Value], prefix)

	if !ItemInFirstClassMap(variable_name.Value, NamesAndScopes.VARIABLES) {
		errors.VariableDoesNotExist(variable_name.Value, GetLocation())
	}

	if reflect.TypeOf(NamesAndScopes.VARIABLES[variable_name.Value]).Name() == "AssignmentNode" {
		variable := NamesAndScopes.VARIABLES[variable_name.Value].(AssignmentNode)
		ControlVariables.current_expected_type = NamesAndScopes.VARIABLES[variable_name.Value].(AssignmentNode).Type
		_type := ControlVariables.current_expected_type

		if !(ControlVariables.current_token.Token_type == lexer.TT_ASSIGNMENT) {
			errors.ExpectedToken("=", GetLocation())
		}
		ParserAdvance()

		reassignment_body := Expression()

		// check for bounding
		var ret []FirstClass 
		if variable.Bounds.BKeep || variable.Bounds.BRoll {
			var UpperBoundExceeded FirstClass
			var LowerBoundExceeded FirstClass
			DefaultAssign := ReAssignmentNode{
				Type:               _type,
				Is_global:          ItemInFirstClassMap(variable_name.Value, NamesAndScopes.GLOBALS),
				Variable_name:      variable_name.Value,
				Body:               reassignment_body,
				Section:            ControlVariables.current_section,
			}

			if variable.Bounds.BKeep {
				UpperBoundExceeded = ReAssignmentNode{
					Type: _type,
					Is_global: ItemInFirstClassMap(variable_name.Value, NamesAndScopes.GLOBALS),
					Variable_name: variable_name.Value,
					Body: variable.Bounds.Upper,
					Section: ControlVariables.current_section,
				}

				LowerBoundExceeded = ReAssignmentNode{
					Type: _type,
					Is_global: ItemInFirstClassMap(variable_name.Value, NamesAndScopes.GLOBALS),
					Variable_name: variable_name.Value,
					Body: variable.Bounds.Lower,
					Section: ControlVariables.current_section,
				}
			} else {
				UpperBoundExceeded = ReAssignmentNode{
					Type: _type,
					Is_global: ItemInFirstClassMap(variable.Variable_name, NamesAndScopes.GLOBALS),
					Variable_name: variable_name.Value,
					Body: variable.Bounds.Lower,
				}

				LowerBoundExceeded = ReAssignmentNode{
					Type: _type,
					Is_global: ItemInFirstClassMap(variable.Variable_name, NamesAndScopes.GLOBALS),
					Variable_name: variable_name.Value,
					Body: variable.Bounds.Upper,
				}
			}

			ret = append(ret, IfNode{
				Type: _type,
				Boolean_statement: BinaryOperationNode{
					Left_branch: reassignment_body,
					Operator_token: ">",
					Right_branch: variable.Bounds.Upper,
				},

				CodeBlock: []FirstClass{UpperBoundExceeded},

				Else: true,
				Else_body: []FirstClass{IfNode{
					Type: _type,
					Boolean_statement: BinaryOperationNode{
						Left_branch: reassignment_body,
						Operator_token: "<",
						Right_branch: variable.Bounds.Lower,
					},

					CodeBlock: []FirstClass{LowerBoundExceeded},

					Else: true,
					Else_body: []FirstClass{DefaultAssign},
				}},
			})

		} else {
			return []FirstClass{ReAssignmentNode{
				Type:               _type,
				Is_global:          ItemInFirstClassMap(variable_name.Value, NamesAndScopes.GLOBALS),
				Variable_name:      variable_name.Value,
				Body:               reassignment_body,
				Section:            ControlVariables.current_section,
			}}
		}
		return ret
	} else {
		var array_indices []SecondClass
		var dimensions int
		var is_indexed bool = false

		// Check for slice
		for ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACKET && ControlVariables.current_token.Token_type != lexer.TT_END_OF_FILE {
			ParserAdvance()
			dimensions++
			ControlVariables.current_expected_type = utils.IntTypeConstr()
			Flags.compare_types = true
			array_indices = append(array_indices, Expression())
			is_indexed = true

			if !(ControlVariables.current_token.Token_type == lexer.TT_RIGHT_BRACKET) {
				errors.ClosingBracketExpected(GetLocation())
			}
			ParserAdvance()
		}

		ControlVariables.current_expected_type = NamesAndScopes.VARIABLES[variable_name.Value].(ArrayAssignmentNode).Type
		ControlVariables.current_expected_type.Dimension -= int8(dimensions)
		_type := ControlVariables.current_expected_type

		if !(ControlVariables.current_token.Token_type == lexer.TT_ASSIGNMENT) {
			errors.ExpectedToken("=", GetLocation())
		}
		ParserAdvance()

		Flags.compare_types = true
		reassignment_body := Expression()
		Flags.compare_types = false

		if is_indexed {
			return []FirstClass{ArraySliceReassignmentNode{
				Type:               _type,
				Is_global:          ItemInFirstClassMap(variable_name.Value, NamesAndScopes.GLOBALS),
				Variable_name:      variable_name.Value,
				Body:               reassignment_body,
				Indices:            array_indices,
				Lengths:            NamesAndScopes.VARIABLES[variable_name.Value].(ArrayAssignmentNode).Array_length,
				Dimensions: 		NamesAndScopes.VARIABLES[variable_name.Value].(ArrayAssignmentNode).Dimensions,
				Section:            ControlVariables.current_section,
			}}
		} else {
			errors.IsNotAssignableError(fmt.Sprintf("%s[]", _type), GetLocation())
		}
	}

	return nil
}
