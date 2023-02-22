package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
	"fmt"
	"reflect"
	"strings"
)

func GetVariableBehaviorDescriptors() []string {
	var FoundVariableDescriptors []string

	for ControlVariables.current_token.Token_type == lexer.TT_ID && utils.StringInArray(ControlVariables.current_token.Value, utils.VARIABLE_BEHAVIOR_DESCRIPTORS) {
		FoundVariableDescriptors = append(FoundVariableDescriptors, ControlVariables.current_token.Value)
		ParserAdvance()
	}

	return FoundVariableDescriptors
}

func GetTypeDescriptors() []string {
	var FoundTypeDescriptors []string

	for ControlVariables.current_token.Token_type == lexer.TT_ID && utils.StringInArray(ControlVariables.current_token.Value, utils.TYPE_DESCRIPTORS) {
		FoundTypeDescriptors = append(FoundTypeDescriptors, ControlVariables.current_token.Value)
		ParserAdvance()
	}

	return FoundTypeDescriptors
}

func GetType() utils.Type {
	t := utils.Type{BaseType: utils.T_INVALID, Dimension: 0, AdditionalType: "", TypeName: ""}

	if ControlVariables.current_token.Token_type == lexer.TT_ID && utils.StringInArray(ControlVariables.current_token.Value, utils.TYPES) || utils.StringInArray(ControlVariables.current_token.Value, utils.CUSTOM_TYPES) {

		// get type based on the current token
		if utils.StringInArray(ControlVariables.current_token.Value, utils.TYPES) {
			t.BaseType = utils.STRING_TO_TYPE[ControlVariables.current_token.Value]
		} else {
			// custom type like struct
			t.AdditionalType = ControlVariables.current_token.Value
		}

		ParserAdvance()
		return t
	} else {
		// no type found error
		errors.TypeExpectedError(GetLocation())
	}

	return t
}

func GetFullName() (lexer.Token, string) {
	// If valid name
	prefix := ""
	if ControlVariables.current_token.Token_type == lexer.TT_ID && (utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.Variable_names) || ItemInStringMap(ControlVariables.current_token.Value, NamesAndScopes.currently_used_sections) || utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.Struct_names) || utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.Function_names)) {
		if ItemInStringMap(ControlVariables.current_token.Value, NamesAndScopes.currently_used_sections) {
			// if prefixed with used section
			prefix = ControlVariables.current_token.Value
			ParserAdvance()

			if ControlVariables.current_token.Token_type == lexer.TT_DOT {
				ParserAdvance()

				if ControlVariables.current_token.Token_type == lexer.TT_ID && utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.Variable_names) || utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.Struct_names) || utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.Function_names) {
					c := ControlVariables.current_token
					c.Value = GetFullVariableName().Value
					ParserAdvance()
					return c, prefix
				} else {
					// error invalid name in outer scope
					if prefix != "" {
						prefix += "."
					}

					errors.ReferenceError(prefix+ControlVariables.current_token.Value, GetLocation())
				}
			} else {
				// error no variable referenced in section
				errors.ThrowError(fmt.Sprintf("VariableExpectedException. A variable was expected after \"%s\" but not found at %s", prefix, GetLocation()), true)
			}
		}

		c := GetFullVariableName()
		ParserAdvance()
		return c, prefix
	} else {
		// error no variable
		if prefix != "" {
			prefix += "."
		}

		errors.ReferenceError(prefix+ControlVariables.current_token.Value, GetLocation())
	}

	return lexer.Token{}, ""
}

func GetName() string {
	c := ControlVariables.current_token                                                                                                                                                                                                                                                // for simplicity
	if c.Token_type == lexer.TT_ID && !(utils.StringInArray(c.Value, NamesAndScopes.Variable_names) || utils.StringInArray(c.Value, NamesAndScopes.Function_names) || utils.StringInArray(c.Value, NamesAndScopes.Struct_names) || utils.StringInArray(c.Value, utils.CUSTOM_TYPES)) { // Is not taken and an ID

		ParserAdvance()
		return c.Value
	} else {
		if c.Token_type != lexer.TT_ID {
			// ID was expected but %s was found
			errors.IDExpectedError(GetLocation(), c.Value)
		} else {
			// name Already taken
			errors.NameAlreadyTakenError(GetLocation(), c.Value)
		}
	}

	return ""
}

// returns spaces and a pipe
func GetLineOffset(line int) string {
	return_string := errors.BLUE + "  "
	length := len(fmt.Sprint(line))

	for ; length != 0; length-- {
		return_string += " "
	}

	return return_string + "| " + errors.COLOR_RESET
}

// returns the code snippet, that is causing an error
func GetErrorCodeSnippet() string {
	if ControlVariables.current_token.Value == "" || utils.StringInArray(ControlVariables.current_token.Value, []string{")", "]", "}"}) {
		ControlVariables.current_token.Value = tokens[ControlVariables.position-1].Value
	}
	if ControlVariables.token_causing_error != "" {
		ControlVariables.current_token.Value = ControlVariables.token_causing_error
	}

	text_at_location := utils.RAW_TEXT[ControlVariables.current_token.File][ControlVariables.current_token.Line]

	return_string := text_at_location + "\n" + GetLineOffset(ControlVariables.current_line) + "\t"
	error_index := strings.Index(return_string, ControlVariables.current_token.Value)
	
	for i := 0; i <= len(text_at_location); i++ {
		return_string += errors.RED
		if i >= error_index && i < error_index+len(ControlVariables.current_token.Value) {
			return_string += "^"
		} else {
			return_string += " "
		}
	}

	return_string = errors.COLOR_RESET + fmt.Sprintf("%s %d |\t%s", errors.BLUE, ControlVariables.current_line+1, errors.COLOR_RESET) + return_string + errors.COLOR_RESET
	return return_string
}

// returns a formatted string with a location
func GetLocation() string {
	return fmt.Sprintf("File: \"%s\" Section: %s; Code:\n%s\n%s", ControlVariables.current_token.File, ControlVariables.current_section, GetLineOffset(ControlVariables.current_line), GetErrorCodeSnippet())
}

func GetTypeOfToken(value int8, num string) utils.Type {
	switch value {
	case lexer.TT_INT:
		{
			return utils.IntTypeConstr()
		}
	case lexer.TT_FLOAT:
		{
			return utils.FloatTypeConstr()
		}
	case lexer.TT_STRING:
		{
			return utils.StringTypeConstr()
		}
	case lexer.TT_CHARACTER:
		{
			return utils.CharTypeConstr()
		}
	}

	return utils.InvalidTypeConstr()
}

// returns the type another type is built on to compare
func GetBaseType(i int8) int8 {
	switch i {
	case utils.T_COCK:
		{
			return utils.T_INT
		}
	default:
		{
			return i
		}
	}
}

// returns a function definition node
func GetFunctionArguments() FunctionAssignArguments {
	var Arguments []FirstClass
	var Types []utils.Type

	if ControlVariables.current_token.Token_type == lexer.TT_RIGHT_PARENTHESIS {
		ParserAdvance()
		return FunctionAssignArguments{
			Arguments:    nil,
			TypeArray:    nil,
		}
	}

	line_cpy := ControlVariables.current_line - 1
	for ControlVariables.current_token.Token_type != lexer.TT_RIGHT_PARENTHESIS {
		if ControlVariables.current_line-1 != line_cpy {
			break
		}
		Arguments = append(Arguments, Assignment())

		switch reflect.TypeOf(Arguments[len(Arguments)-1]).Name() {
		case "AssignmentNode":
			{
				temp := Arguments[len(Arguments)-1].(AssignmentNode)

				if temp.Variable_body != nil {
					errors.NoInitializedExpressionExpected(GetLocation())
				}
				Types = append(Types, temp.Type)
			}
		case "ArrayAssignmentNode":
			{
				temp := Arguments[len(Arguments)-1].(ArrayAssignmentNode)
				Types = append(Types, temp.Type)
			}
		}

		if ControlVariables.current_token.Token_type == lexer.TT_COMMA {
			if GetToken(1).Token_type != lexer.TT_RIGHT_PARENTHESIS {
				ParserAdvance()
			} else {
				errors.FunctionArgumentExpectedError(GetLocation())
			}
		}
	}
	if ControlVariables.current_token.Token_type != lexer.TT_RIGHT_PARENTHESIS {
		ControlVariables.current_line = line_cpy - 1
		ControlVariables.current_line_string = NamesAndScopes.PURE_SOURCE_CODE[line_cpy-1]
		ControlVariables.current_token = tokens[ControlVariables.position-1]
		errors.ClosingParenthesisExpected(GetLocation())
	}
	ParserAdvance()

	return FunctionAssignArguments{
		Arguments:    Arguments,
		TypeArray:    Types,
	}
}

// Returns a token that is somewhere in the array of tokens
func GetToken(number int) lexer.Token {
	position_now := ControlVariables.position

	if position_now+number < len(tokens) {
		return tokens[position_now+number]
	} else {
		Flags.is_out_of_tokens = true
		return lexer.Token{}
	}
}

// Returns the full name of a variable. for example struct accesses
func GetFullVariableName() lexer.Token {
	name_node := ControlVariables.current_token

	if (GetToken(1).Token_type == lexer.TT_DOT || GetToken(1).Token_type == lexer.TT_ARROW) && GetToken(2).Token_type == lexer.TT_ID && (!utils.StringInArray(GetToken(2).Value, utils.TYPES) && !utils.StringInArray(GetToken(2).Value, utils.INSTRUCTIONS) && !utils.StringInArray(GetToken(2).Value, utils.CUSTOM_TYPES)) {
		// iterates through the tokens for as long as it's a contiguous name with dots like -> hello.type
		for (GetToken(1).Token_type == lexer.TT_DOT || GetToken(1).Token_type == lexer.TT_ARROW) && GetToken(2).Token_type == lexer.TT_ID && (!utils.StringInArray(GetToken(2).Value, utils.TYPES) && !utils.StringInArray(GetToken(2).Value, utils.INSTRUCTIONS) && !utils.StringInArray(GetToken(2).Value, utils.CUSTOM_TYPES)) {
			ParserAdvance()
			name_node.Value += ControlVariables.current_token.Value
			ParserAdvance()
			name_node.Value += ControlVariables.current_token.Value
		}

		if !utils.StringInArray(name_node.Value, NamesAndScopes.Variable_names) || !utils.StringInArray(name_node.Value, NamesAndScopes.Struct_names) || !utils.StringInArray(name_node.Value, NamesAndScopes.Function_names) {
			errors.ReferenceError(name_node.Value, GetLocation())
		}
	}
	return name_node
}

// returns a typecast node
func GetAutoConvert(t utils.Type, value SecondClass) (TypeCastNode, bool) {
	if t.AdditionalType == "" {
		if utils.TYPE_HIERARCHY[t.BaseType] == utils.TYPE_HIERARCHY[ControlVariables.current_expected_type.BaseType] {
			return TypeCastNode{}, false
		} else {
			n := TypeCastNode{
				Expression: value,
				From: t,
				To: ControlVariables.current_expected_type,
			}
			return n, true
		}
	} else {
		return TypeCastNode{}, false
	}
}