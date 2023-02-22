package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
	"fmt"
	"reflect"
)

// returns the node to a term a term is a expression with only * and /
func Term() SecondClass {
	operators := []int8{lexer.TT_ASTERISK, lexer.TT_SLASH, lexer.TT_BANG}
	return BinaryOperation(true, false, false, operators)
}

// returns a whole mathematical expression
// Boolean expressions (== or !=) as the order of operations in boolean algebra is as follows -> ((==, !=), &&, ||, !)
func SubExpression() SecondClass {
	operators := []int8{lexer.TT_PLUS, lexer.TT_MINUS, lexer.TT_BIT_SHIT_LEFT, lexer.TT_BIT_AND, lexer.TT_BIT_OR, lexer.TT_BIT_XOR, lexer.TT_EQUALS, lexer.TT_NOT_EQUALS, lexer.TT_LESS_EQUAL, lexer.TT_LESS_THAN, lexer.TT_GREATER_EQUAL, lexer.TT_GREATER_THAN}
	return BinaryOperation(false, true, false, operators)
}

func Expression() SecondClass {
	operators := []int8{lexer.TT_AND, lexer.TT_OR}
	return BinaryOperation(false, false, true, operators)
}

// returns a factor. A factor can be a literal or another expression in parentheses
func Factor(is_called_as_first_class bool) SecondClass {
	token := ControlVariables.current_token

	if ControlVariables.current_token.Line != ControlVariables.current_line {
		return UniversalNone{}
	}

	var not bool = false
	var minus_prefix bool = false
	var bitwise_not bool = false

	// returns either a function call variable name or array slice
	if token.Token_type == lexer.TT_MINUS {
		minus_prefix = true
		ParserAdvance()
	}
	if token.Token_type == lexer.TT_BIT_NOT {
		bitwise_not = true
		ParserAdvance()
	}

	if token.Token_type == lexer.TT_BANG && Flags.in_boolean_expression {
		not = true
		ParserAdvance()
	}

	if ControlVariables.current_token.Token_type == lexer.TT_LEFT_PARENTHESIS {
		ParserAdvance()
		expression := Expression()

		if ControlVariables.current_token.Token_type != lexer.TT_RIGHT_PARENTHESIS {
			errors.ClosingParenthesisExpected(GetLocation())
		}

		ParserAdvance()
		return FactorExpression{Binary_operation_node: expression, not: not, minus_prefix: minus_prefix, bitwise_not: bitwise_not}
	}

	token = ControlVariables.current_token
	prefix := ""
	if token.Token_type == lexer.TT_ID {
		if utils.StringInArray(token.Value, utils.RESERVED_KEYWORD_CONSTANTS) {
			TypeCompare(utils.RESERVED_KEYWORD_CONSTANTS_TYPES[token.Value])
			ParserAdvance()

			return DirectNode{Type: token.Token_type, Value: utils.RESERVED_KEYWORD_CONSTANTS_VALUES[token.Value], Is_negated: minus_prefix, Bitwise_not: bitwise_not}
		}

		if utils.StringInArray(ControlVariables.current_token.Value, utils.TYPES) || utils.StringInArray(ControlVariables.current_token.Value, utils.CUSTOM_TYPES) {
			token = ControlVariables.current_token
		} else {
			token, prefix = GetFullName()
		}
		//ParserAdvance()

		// if current token is a variable
		if ItemInFirstClassMap(token.Value, NamesAndScopes.VARIABLES) || utils.StringInArray(token.Value, NamesAndScopes.Variable_names) {
			NamesAndScopes.VARIABLE_REFERENCE_COUNT[token.Value]++

			// if the variable is an array
			if reflect.TypeOf(NamesAndScopes.VARIABLES[token.Value]).Name() == "ArrayAssignmentNode" {
				type_to_compare := NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode)
				ControlVariables.most_significant_type.BaseType = GetBaseType(type_to_compare.Type.BaseType)
				TypeCompare(ControlVariables.most_significant_type)

				// TODO: check if accessed. If not and the type of the variable assigned to is a array with the same dimension, copy array
				if GetToken(1).Token_type != lexer.TT_LEFT_BRACKET && GetToken(1).Line != ControlVariables.current_line { // If no slice is accessed pointer comparison is executed
					//return VariableNameNode{Name: token.Value, Pointers: pointers, Get_address: reference, Not: not, Is_negated: minus_prefix, Bitwise_not: bitwise_not, Is_global: NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode).Is_global}
				}
			}

			// returns a slice if found
			if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACKET {
				var indices []SecondClass
				var dims int = 0

				for ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACKET && ControlVariables.current_token.Token_type != lexer.TT_END_OF_FILE {
					original_expected_type := ControlVariables.current_expected_type
					ControlVariables.current_expected_type.BaseType = utils.IntTypeConstr().BaseType
					original_compare_types := Flags.compare_types
					Flags.compare_types = true
					ParserAdvance()

					indices = append(indices, Expression())

					if ControlVariables.current_token.Token_type == lexer.TT_RIGHT_BRACKET {
						ParserAdvance()
						// return array slice
					} else {
						// error
						errors.ClosingBracketExpected(GetLocation())
					}

					ControlVariables.current_expected_type = original_expected_type
					dims++
					Flags.compare_types = original_compare_types
				}
				t := NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode).Type
				t.Dimension -= int8(dims)
				TypeCompare(t)

				return ArraySliceNode{
					Name:            token.Value,
					Positions:       indices,
					Lengths:         NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode).Array_length,
					Not:             not,
					Is_negated:      minus_prefix,
					Bitwise_not:     bitwise_not,
					Is_global:       NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode).Is_global,
				}
			}

			if reflect.TypeOf(NamesAndScopes.VARIABLES[token.Value]).Name() == "AssignmentNode" {
				compare_type := NamesAndScopes.VARIABLES[token.Value].(AssignmentNode)
				ControlVariables.most_significant_type.BaseType = GetBaseType(compare_type.Type.BaseType)
				s := ControlVariables.current_section
				ScopeValidityCheck(compare_type, prefix)
				TypeCompare(utils.Type{BaseType: GetBaseType(compare_type.Type.BaseType), Dimension: compare_type.Type.Dimension, AdditionalType: compare_type.Type.AdditionalType, TypeName: compare_type.Type.TypeName})
				ControlVariables.current_section = s

				node := VariableNameNode{
					Name:                      token.Value,
					Type:                      NamesAndScopes.VARIABLES[token.Value].(AssignmentNode).Type,
					Not:                       not,
					Is_negated:                minus_prefix,
					Bitwise_not:               bitwise_not,
					Is_global:                 NamesAndScopes.VARIABLES[token.Value].(AssignmentNode).Is_global}
				temp, set := GetAutoConvert(compare_type.Type, node)
				if (set) {
					return temp
				} else {
					return node
				}
			} else {
				// Basically assign the whole array without accessing by index
				compare_type := NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode)
				ControlVariables.most_significant_type.BaseType = GetBaseType(compare_type.Type.BaseType)
				s := ControlVariables.current_section
				ScopeValidityCheck(compare_type, prefix)
				TypeCompare(utils.Type{BaseType: GetBaseType(compare_type.Type.BaseType), Dimension: compare_type.Type.Dimension, AdditionalType: compare_type.Type.AdditionalType, TypeName: compare_type.Type.TypeName})
				ControlVariables.current_section = s

				return WholeArrayAsLiteral{
					Name:            token.Value,
					Lengths:          NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode).Array_length,
					Dimensions: NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode).Dimensions,
					Not:             not,
					Is_negated:      minus_prefix,
					Bitwise_not:     bitwise_not,
					Is_global:       NamesAndScopes.VARIABLES[token.Value].(ArrayAssignmentNode).Is_global,
				}
			}
		} else if ItemInFirstClassMap(token.Value, NamesAndScopes.FUNCTIONS) || utils.StringInArray(token.Value, NamesAndScopes.Function_names) {
			NamesAndScopes.VARIABLE_REFERENCE_COUNT[token.Value]++
			called_function_name := token.Value

			if ItemInFirstClassMap(token.Value, NamesAndScopes.FUNCTIONS) {
				t := NamesAndScopes.FUNCTIONS[called_function_name].(FunctionNode).Return_type
				t.BaseType = GetBaseType(t.BaseType)
				ControlVariables.most_significant_type = t

				// if void or noret and is called as first class
				if !is_called_as_first_class && (ControlVariables.most_significant_type.BaseType == utils.T_VOID || utils.StringInArray("__noret__", NamesAndScopes.FUNCTIONS[called_function_name].(FunctionNode).Variable_behavior_descriptors)) {
					TypeCompare(ControlVariables.most_significant_type)
				}

				if ControlVariables.current_token.Token_type == lexer.TT_LEFT_PARENTHESIS {
					ParserAdvance()
					prev_expected_type := ControlVariables.current_expected_type
					ControlVariables.current_expected_type = t
					var function_call_arguments []SecondClass
					ScopeValidityCheck(NamesAndScopes.FUNCTIONS[called_function_name], prefix)

					if ItemInFirstClassMap(token.Value, NamesAndScopes.FUNCTIONS) {
						function_call_arguments = ParseFunctionCallArguments(lexer.TT_RIGHT_PARENTHESIS, ArgumentLength(NamesAndScopes.FUNCTIONS[called_function_name]), true, called_function_name)
					} else {
						function_call_arguments = ParseFunctionCallArguments(lexer.TT_RIGHT_PARENTHESIS, 0, false, called_function_name)
					}

					ParserAdvance()

					ControlVariables.current_expected_type = prev_expected_type

					node := FunctionCallNode{Function_name: called_function_name, Called_function_arguments: function_call_arguments, Is_negated: minus_prefix, Bitwise_not: bitwise_not}
					temp, set := GetAutoConvert(t, node)

					if set {
						return temp
					} else {
						return node
					}
				} else {
					errors.OpeningParenthesisExpected(GetLocation())
				}
			}
		} else if utils.StringInArray(token.Value, utils.TYPES) || utils.StringInArray(token.Value, utils.CUSTOM_TYPES) {
			// probable typecast
			// looks like: <type>(<expr>)
			type_name := GetType()

			if ControlVariables.current_token.Token_type == lexer.TT_LEFT_PARENTHESIS {
				// set stuff like expected pointers, type, etc...
				previous_expected_type := ControlVariables.current_expected_type
				ControlVariables.current_expected_type = utils.Type{BaseType: utils.T_INVALID, Dimension: 0, AdditionalType: "", TypeName: "any"}

				ParserAdvance()
				expression := Expression()

				if ControlVariables.current_token.Token_type != lexer.TT_RIGHT_PARENTHESIS {
					errors.ClosingParenthesisExpected(GetLocation())
				}

				ParserAdvance()
				factor_expression := FactorExpression{Binary_operation_node: expression, not: not, minus_prefix: minus_prefix, bitwise_not: bitwise_not}

				ControlVariables.current_expected_type = previous_expected_type

				var return_node SecondClass
				if (ControlVariables.most_significant_type.BaseType == GetBaseType(type_name.BaseType)) || (ControlVariables.most_significant_type.BaseType == utils.T_INT && GetBaseType(type_name.BaseType) == utils.T_FLOAT) {
					// alter return node to just be a factor expression, as this typecast is redundant and not needed
					return_node = factor_expression
					t := ControlVariables.most_significant_type
					t.BaseType = GetBaseType(type_name.BaseType)
					ControlVariables.most_significant_type = t
				} else {

					return_node = TypeCastNode{
						Expression:  factor_expression,
						From:        ControlVariables.most_significant_type,
						To:          utils.Type{BaseType: GetBaseType(type_name.BaseType), Dimension: type_name.Dimension, AdditionalType: type_name.AdditionalType, TypeName: type_name.TypeName},
					}
					ControlVariables.most_significant_type = return_node.(TypeCastNode).To
				}
				return return_node
			} else {
				errors.OpeningParenthesisExpected(GetLocation())
			}
		} else {
			errors.ReferenceError(prefix+token.Value, GetLocation())
		}
	}

	if GetTypeOfToken(token.Token_type, token.Value).BaseType == utils.T_CHAR {
		token.Value = fmt.Sprint(int(rune(token.Value[0])))
		token.Token_type = lexer.TT_INT
	}
	TypeCompare(GetTypeOfToken(token.Token_type, token.Value))
	node := DirectNode{Type: token.Token_type, Value: token.Value, Is_negated: minus_prefix, Bitwise_not: bitwise_not}
	temp, set := GetAutoConvert(GetTypeOfToken(token.Token_type, token.Value), node)
	ParserAdvance()

	if (set) {
		return temp
	} else {
		return node
	}
}

// Returns the node of a binary operation
func BinaryOperation(factor, term, sub_expression bool, operators []int8) SecondClass {
	var Left_branch SecondClass
	var current_line_in_expression int = ControlVariables.current_token.Line

	if sub_expression {
		Left_branch = SubExpression()
	} else if term {
		Left_branch = Term()
	} else if factor {
		Left_branch = Factor(false)
	}

	for utils.Int8InArray(ControlVariables.current_token.Token_type, operators) && !Flags.is_out_of_tokens && ControlVariables.current_line == current_line_in_expression {
		Operator_token := ControlVariables.current_token
		ParserAdvance()

		if utils.Int8InArray(ControlVariables.current_token.Token_type, []int8{lexer.TT_AND, lexer.TT_OR, lexer.TT_BANG, lexer.TT_EQUALS, lexer.TT_NOT_EQUALS, lexer.TT_LESS_THAN, lexer.TT_LESS_EQUAL, lexer.TT_GREATER_THAN, lexer.TT_GREATER_EQUAL}) {
			ControlVariables.token_causing_error = ControlVariables.current_token.Value
			errors.TypeMismatchError(ControlVariables.current_expected_type.TypeName, "bool", GetLocation())
		}

		var Right_branch SecondClass
		if sub_expression {
			Right_branch = SubExpression()
		} else if term {
			Right_branch = Term()
		} else if factor {
			Right_branch = Factor(false)
		}
		Left_branch = BinaryOperationNode{Left_branch: Left_branch, Operator_token: Operator_token.Value, Right_branch: Right_branch}
	}

	return Left_branch
}
