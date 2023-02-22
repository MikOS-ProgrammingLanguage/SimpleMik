package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
)

var tokens []lexer.Token
var root_node RootNode = RootNode{}
var knownBigO map[string]string

func ParseBody(break_token int8) []FirstClass {
	var ret []FirstClass
	var attributes []int8
	var attribute_location lexer.Token
	var attributes_set bool
	for !Flags.is_out_of_tokens && ControlVariables.current_token.Token_type != break_token {
		if len(attributes) != 0 {
			attributes_set = true
		}
		if ControlVariables.current_token.Token_type == lexer.TT_ID && (utils.StringInArray(ControlVariables.current_token.Value, utils.TYPES) || utils.StringInArray(ControlVariables.current_token.Value, utils.CUSTOM_TYPES) || utils.StringInArray(ControlVariables.current_token.Value, utils.TYPE_DESCRIPTORS) || utils.StringInArray(ControlVariables.current_token.Value, utils.VARIABLE_BEHAVIOR_DESCRIPTORS)) || ControlVariables.current_token.Value == "use" {
			// main switch
			if utils.StringInArray(ControlVariables.current_token.Value, utils.TYPES) || utils.StringInArray(ControlVariables.current_token.Value, utils.VARIABLE_BEHAVIOR_DESCRIPTORS) || utils.StringInArray(ControlVariables.current_token.Value, utils.TYPE_DESCRIPTORS) || utils.StringInArray(ControlVariables.current_token.Value, utils.CUSTOM_TYPES) {
				ret = append(ret, Assignment())
			} else if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "use" {
				EvaluateUse()
			}
		} else if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "if" {
			// if statement
			ret = append(ret, If())
		} else if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "while" {
			// while loop
			ret = append(ret, While())
		} else if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "mikf" {
			// normal function
			if Flags.in_function {
				errors.CantNestFunctions(GetLocation())
			}
			ret = append(ret, Function(nil))
		} else if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "struct" {
			// struct
			if attributes_set {
				ret = append(ret, Struct(attributes))
				attributes_set = false
				attributes = []int8{}
			} else {
				ret = append(ret, Struct([]int8{}))
			}
		} else if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "return" {
			// return expression
			ret = append(ret, Return())
		} else if ControlVariables.current_token.Token_type == lexer.TT_ID {
			// reassign
			ret = append(ret, Reassignment()...)
		} else if ControlVariables.current_token.Token_type == lexer.TT_ATTRIBUTE {
			// attributes
			attribute_location = ControlVariables.current_token
			attributes = Attributes()
		} else {
			// unexpected error
			errors.ThrowWarning("NOT IMPLEMENTED:::QUITTING", true)
		}

		if (attributes_set || (ControlVariables.position >= len(tokens))) && (len(attributes) != 0) {
			// error. Unused attributes
			var attribute_string string = "("
			var new_map map[int8]string = ReverseStringInt8Map(utils.ATTRIBUTES)
			for _, i := range attributes {
				attribute_string += new_map[i] + ", "
			}
			ControlVariables.current_token = attribute_location
			errors.UnusedAttributesError(len(attributes), attribute_string, GetLocation())
		}
		//fmt.Println(ret[len(ret)-1])
	}

	return ret
}

func Parse(Tokens []lexer.Token, lexed_sections []string, pure_code []string, fuck_you bool) (RootNode, map[string]string) {
	// implement standard built-ins
	NamesAndScopes.FUNCTIONS["puts"] = FunctionNode{
		Declared:      true,
		Function_name: "puts",
		Arguments: FunctionAssignArguments{
			Arguments: []FirstClass{AssignmentNode{
				Variable_behavior_descriptors: nil,
				Type_descriptors:              nil,
				Type:                          utils.StringTypeConstr(),
				Is_global:                     false,
				Variable_name:                 "string",
				Variable_body:                 nil,
			}},
			TypeArray:    []utils.Type{utils.StringTypeConstr()},
		},
		Return_type:                   utils.IntTypeConstr(),
		CodeBlock:                     nil,
		Variable_behavior_descriptors: []string{"__noret__"},
		VARIABLES:                     nil,
		FUNCTIONS:                     nil,
		Structs:                       nil,
	}
	NamesAndScopes.Function_names = append(NamesAndScopes.Function_names, "puts")
	
	errors.FUCK_YOU = fuck_you

	NamesAndScopes.PURE_SOURCE_CODE = pure_code
	ControlVariables.current_line_string = NamesAndScopes.PURE_SOURCE_CODE[0]
	tokens = Tokens
	NamesAndScopes.sections = lexed_sections
	ParserAdvance()
	root_node.Nodes = append(root_node.Nodes, ParseBody(lexer.TT_PLACE_HOLDER)...)
	return root_node, knownBigO
}
