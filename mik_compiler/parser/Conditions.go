package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
)

func If() FirstClass {
	ParserAdvance()

	ControlVariables.current_expected_type = utils.InvalidTypeConstr()
	Flags.compare_types = false
	boolean_expression := Expression()
	goal_type := ControlVariables.most_significant_type

	if !utils.StringInArray(boolean_expression.(BinaryOperationNode).Operator_token, utils.OPERATORS) {
		errors.TypeMismatchError("bool", ControlVariables.most_significant_type.TypeName, GetLocation())
	}

	Flags.compare_types = true

	if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACE {
		var else_exists bool = false
		var else_body []FirstClass

		previous_variable_names := NamesAndScopes.Variable_names
		previous_variables := NamesAndScopes.VARIABLES
		ParserAdvance()
		body := ParseBody(lexer.TT_RIGHT_BRACE)
		ParserAdvance()

		NamesAndScopes.Variable_names = previous_variable_names
		NamesAndScopes.VARIABLES = previous_variables

		if ControlVariables.current_token.Value == "else" {
			else_exists = true
			ParserAdvance()
			if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACE {
				ParserAdvance()
				else_body = ParseBody(lexer.TT_RIGHT_BRACE)
				ParserAdvance()
				NamesAndScopes.Variable_names = previous_variable_names
				NamesAndScopes.VARIABLES = previous_variables
			} else {
				errors.OpeningBraceExpected(GetLocation())
			}
		} else if ControlVariables.current_token.Value == "elif" {
			else_exists = true
			else_body = []FirstClass{If()}
		}

		return IfNode{
			Boolean_statement: boolean_expression,
			CodeBlock:         body,
			Type:              goal_type,
			Else:              else_exists,
			Else_body:         else_body,
		}
	} else {
		errors.OpeningBraceExpected(GetLocation())
	}
	return nil
}

func Else() FirstClass {
	return nil
}
