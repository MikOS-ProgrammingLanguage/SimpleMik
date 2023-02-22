package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
	"reflect"
)

// Parses arguments in a function with a type to end the statement, a count of expected arguments, and if the length is expected at all
func ParseFunctionCallArguments(end_of_call_arguments int8, length_of_expected_arguments int, length_expected bool, name string) []SecondClass {
	var call_arguments []SecondClass
	var iterations int = 0
	var parent_function_assign_body FunctionAssignArguments
	if reflect.TypeOf(NamesAndScopes.FUNCTIONS[name]).Name() == "FunctionNode" {
		parent_function_assign_body = NamesAndScopes.FUNCTIONS[name].(FunctionNode).Arguments
	}

	var cnt int = 0
	var comma_is_set bool = true
	line_cpy := ControlVariables.current_line
	for ControlVariables.current_token.Token_type != end_of_call_arguments && ControlVariables.current_line == line_cpy {
		if length_expected {
			if cnt == length_of_expected_arguments { // exceeded length
				errors.ToManyArgumentsException(GetLocation(), name, length_of_expected_arguments, cnt+1)
			}
		}

		var argument SecondClass
		next_token := GetToken(1)
		the_next_but_one_token := GetToken(2)

		// validates commas
		if next_token.Token_type == lexer.TT_COMMA && (utils.StringInArray(the_next_but_one_token.Value, utils.TYPES) || utils.StringInArray(the_next_but_one_token.Value, utils.CUSTOM_TYPES)) && cnt >= length_of_expected_arguments {
			errors.ArgumentExpectedError(GetLocation())
		}

		if ControlVariables.current_token.Token_type == lexer.TT_COMMA {
			if comma_is_set {
				errors.UnexpectedTokenError(GetLocation(), rune(ControlVariables.current_token.Value[0]))
			}
			comma_is_set = true
			ParserAdvance()
			continue
		} else {
			if !comma_is_set {
				errors.ExpectedToken(",", GetLocation())
			}
			comma_is_set = false
			prev_expected_type := ControlVariables.current_expected_type
			ControlVariables.current_expected_type = parent_function_assign_body.TypeArray[cnt]

			argument = Expression()
			call_arguments = append(call_arguments, argument)

			ControlVariables.current_expected_type = prev_expected_type
			cnt++
		}
		iterations++
	}
	if ControlVariables.current_token.Token_type != end_of_call_arguments {
		ControlVariables.current_line = line_cpy
		ControlVariables.current_line_string = NamesAndScopes.PURE_SOURCE_CODE[line_cpy-1]
		ControlVariables.current_token = tokens[ControlVariables.position-1]
		ControlVariables.current_section = ControlVariables.current_token.Section_name
		errors.ClosingParenthesisExpected(GetLocation())
	}

	if cnt < len(parent_function_assign_body.Arguments) {
		errors.ToFewArgumentsException(GetLocation(), name, len(parent_function_assign_body.Arguments), cnt)
	}
	return call_arguments
}

// Returns how many arguments are expected for a function
func ArgumentLength(node FirstClass) int {
	return len(node.(FunctionNode).Arguments.Arguments)
}
