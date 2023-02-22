package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
)

func Assignment() FirstClass {
	var expression SecondClass = nil

	Section := ControlVariables.current_token.Section_name

	VariableBehaviorDescriptors := VariableBehaviorDescriptorValidityCheck(GetVariableBehaviorDescriptors())
	if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "mikf" {
		if Flags.in_function {
			errors.CantNestFunctions(GetLocation())
		}
		return Function(VariableBehaviorDescriptors)
	}

	TypeDescriptors := TypeDescriptorValidityCheck(GetTypeDescriptors())
	Type_ := GetType()
	ControlVariables.current_expected_type = Type_
	ControlVariables.current_section = Section
	Flags.compare_types = true

	Name := CheckNameTaken(GetName())
	is_global := utils.StringInArray("global", VariableBehaviorDescriptors)
	bound := Bound{
		BKeep: false,
		BRoll: false,
		Upper: nil,
		Lower: nil,
	}

	// check for bounding var
	if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACE {
		// check if variable is of type int or char or float
		if !(utils.TYPE_HIERARCHY[Type_.BaseType] >= utils.TYPE_HIERARCHY[utils.T_CHAR] && utils.TYPE_HIERARCHY[Type_.BaseType] <= utils.TYPE_HIERARCHY[utils.T_FLOAT]) {
			errors.ThrowError("Boundaries not defined for values other than int or float.", true)
		}

		// lower bound
		ParserAdvance()
		ControlVariables.current_expected_type = Type_
		bound.Lower = Expression()

		// comma separation
		if ControlVariables.current_token.Token_type == lexer.TT_COMMA {
			ParserAdvance()
		} else {
			errors.ExpectedToken(",", GetLocation())
		}

		// upper bound
		bound.Upper = Expression()

		// comma separation
		if ControlVariables.current_token.Token_type == lexer.TT_COMMA {
			ParserAdvance()

			// check for specified roll or keep
			if ControlVariables.current_token.Token_type == lexer.TT_ID {
				if ControlVariables.current_token.Value == "keep" {
					bound.BKeep = true
				} else if ControlVariables.current_token.Value == "roll" {
					bound.BRoll = true
				} else {
					errors.ExpectedToken("roll/keep", GetLocation())
				}
				ParserAdvance()
			} else {
				errors.IDExpectedError(GetLocation(), ControlVariables.current_token.Value)
			}
		} else {
			bound.BKeep = true
		}

		// check for closing brace
		if ControlVariables.current_token.Token_type == lexer.TT_RIGHT_BRACE {
			ParserAdvance()
		} else {
			errors.ClosingBraceExpected(GetLocation())
		}
	}

	if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACKET { // Is array
		
		var dimension_lengths []SecondClass
		for (ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACKET) && (ControlVariables.current_token.Token_type != lexer.TT_END_OF_FILE) {
			Type_.Dimension++
			var array_length_expression SecondClass
			ParserAdvance()

			if !Flags.in_function_definition {
				ControlVariables.current_expected_type.BaseType = utils.T_INT
				Flags.compare_types = true

				array_length_expression = Expression()
				Flags.compare_types = false
			}

			if ControlVariables.current_token.Token_type == lexer.TT_RIGHT_BRACKET {
				ParserAdvance()
			} else {
				errors.ClosingBracketExpected(GetLocation())
			}

			dimension_lengths = append(dimension_lengths, array_length_expression)
		}
		return_array_assignment := ArrayAssignmentNode{
			Variable_behavior_descriptors: VariableBehaviorDescriptors,
			Type_descriptors:              TypeDescriptors,
			Type:                          Type_,
			Is_global:                     is_global,
			Variable_name:                 Name,
			Array_length:                  dimension_lengths,
			Section:                       ControlVariables.current_token.Section_name,
		}

		NamesAndScopes.VARIABLES[Name] = return_array_assignment
		NamesAndScopes.Variable_names = append(NamesAndScopes.Variable_names, return_array_assignment.Variable_name)
		return return_array_assignment
	}

	ControlVariables.current_expected_type = Type_

	if ControlVariables.current_token.Token_type == lexer.TT_ASSIGNMENT {
		ParserAdvance()
		expression = Expression()
	}

	NamesAndScopes.Variable_names = append(NamesAndScopes.Variable_names, Name)
	ReturnNode := AssignmentNode{
		Variable_behavior_descriptors: VariableBehaviorDescriptors,
		Type_descriptors:              TypeDescriptors,
		Type:                          Type_,
		Is_global:                     is_global,
		Variable_name:                 Name,
		Variable_body:                 expression,
		Section:                       Section,
		Bounds: bound,
	}

	NamesAndScopes.VARIABLES[Name] = ReturnNode
	ControlVariables.current_line = ControlVariables.current_token.Line
	ControlVariables.current_section = ControlVariables.current_token.Section_name
	return ReturnNode
}
