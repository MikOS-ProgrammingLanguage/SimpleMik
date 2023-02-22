package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
)

func Function(Variable_behavior_descriptors []string) FirstClass {
	Flags.in_function = true
	ParserAdvance()
	var function_name string
	variable_behavior_descriptors := GetVariableBehaviorDescriptors()
	Variable_behavior_descriptors = append(Variable_behavior_descriptors, variable_behavior_descriptors...)
	if ControlVariables.current_token.Token_type == lexer.TT_ID {
		function_name = CheckNameTaken(ControlVariables.current_token.Value)
		NamesAndScopes.Function_names = append(NamesAndScopes.Function_names, function_name)
		ParserAdvance()

		if ControlVariables.current_token.Token_type == lexer.TT_LEFT_PARENTHESIS {
			ParserAdvance()
			Flags.in_function_definition = true

			// change variable scopes. NamesAndScopes.FUNCTIONS can stay
			old_VARIABLES := NamesAndScopes.VARIABLES
			old_Variable_names := NamesAndScopes.Variable_names
			NamesAndScopes.VARIABLES = make(map[string]FirstClass)

			arguments := GetFunctionArguments()

			var return_type utils.Type
			if ControlVariables.current_token.Token_type == lexer.TT_ARROW {
				ParserAdvance()

				if ControlVariables.current_token.Token_type == lexer.TT_ID && (utils.StringInArray(ControlVariables.current_token.Value, utils.TYPES) || utils.StringInArray(ControlVariables.current_token.Value, utils.CUSTOM_TYPES)) {
					return_type = GetType()
					Flags.expect_return_statement = true
				} else {
					// error. Unexpected return type
					errors.UnexpectedReturnType(GetLocation(), ControlVariables.current_token.Value)
				}
			} else {
				return_type = utils.VoidTypeConstr()
				Flags.expect_return_statement = false
			}
			Flags.in_function_definition = false

			if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACE {
				ParserAdvance()
				// function

				ControlVariables.expected_function_return_type = return_type
				body := ParseBody(lexer.TT_RIGHT_BRACE)
				ParserAdvance()

				if !Flags.return_hit && Flags.expect_return_statement {
					errors.ReturnStatementExpected(GetLocation())
				}

				Flags.expect_return_statement = false
				NamesAndScopes.Variable_names = old_Variable_names
				Flags.in_function = false
				Flags.return_hit = false
				return_node := FunctionNode{
					Declared:                      false,
					Function_name:                 function_name,
					Arguments:                     arguments,
					Return_type:                   return_type,
					CodeBlock:                     body,
					VARIABLES:                     NamesAndScopes.VARIABLES,
					Section:                       ControlVariables.current_section,
					Variable_behavior_descriptors: Variable_behavior_descriptors,
				}
				NamesAndScopes.FUNCTIONS[function_name] = return_node
				NamesAndScopes.Function_names = append(NamesAndScopes.Function_names, function_name)
				NamesAndScopes.VARIABLES = old_VARIABLES
				return return_node
			} else {
				// function declaration
				NamesAndScopes.VARIABLES = old_VARIABLES
				NamesAndScopes.Variable_names = old_Variable_names
				Flags.in_function = false
				return_node := FunctionNode{
					Declared:                      true,
					Function_name:                 function_name,
					Arguments:                     arguments,
					Return_type:                   return_type,
					CodeBlock:                     nil,
					Section:                       ControlVariables.current_section,
					Variable_behavior_descriptors: Variable_behavior_descriptors,
				}
				NamesAndScopes.FUNCTIONS[function_name] = return_node
				NamesAndScopes.Function_names = append(NamesAndScopes.Function_names, function_name)
				return return_node
			}
		}
	} else {
		// error
		errors.IDExpectedError(GetLocation(), ControlVariables.current_token.Value)
	}
	return nil
}
