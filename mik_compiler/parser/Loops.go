package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
)

func While() FirstClass {
	ParserAdvance()

	ControlVariables.current_expected_type.BaseType = utils.InvalidTypeConstr().BaseType
	Flags.compare_types = false
	boolean_expression := Expression()
	goal_type := ControlVariables.most_significant_type

	if !utils.StringInArray(boolean_expression.(BinaryOperationNode).Operator_token, utils.OPERATORS) {
		errors.TypeMismatchError("bool", ControlVariables.most_significant_type.TypeName, GetLocation())
	}

	if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACE {
		previous_variable_names := NamesAndScopes.Variable_names
		previous_variables := NamesAndScopes.VARIABLES
		ParserAdvance()
		body := ParseBody(lexer.TT_RIGHT_BRACE)
		ParserAdvance()

		NamesAndScopes.Variable_names = previous_variable_names
		NamesAndScopes.VARIABLES = previous_variables

		return WhileNode{
			Boolean_statement: boolean_expression,
			CodeBlock:         body,
			Type:              goal_type,
		}
	} else {
		errors.OpeningBraceExpected(GetLocation())
	}

	return nil
}
