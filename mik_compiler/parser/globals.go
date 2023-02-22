package parser

import (
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
)

// All global names and scopes
type NamesAndScopes_t struct {
	sections                 []string
	currently_used_sections  map[string]string // sections and their aliases. e.g.: use test as t -> map[t]test
	Function_names           []string
	Variable_names           []string
	Struct_names             []string
	VARIABLES                map[string]FirstClass
	FUNCTIONS                map[string]FirstClass
	STRUCTS                  map[string]StructNode
	GLOBALS                  map[string]FirstClass
	VARIABLE_REFERENCE_COUNT map[string]int
	PURE_SOURCE_CODE         []string
}

var NamesAndScopes NamesAndScopes_t = NamesAndScopes_t{
	sections:                 []string{},
	Struct_names:             []string{},
	currently_used_sections:  make(map[string]string),
	VARIABLES:                make(map[string]FirstClass),
	FUNCTIONS:                make(map[string]FirstClass),
	STRUCTS:                  make(map[string]StructNode),
	GLOBALS:                  make(map[string]FirstClass),
	VARIABLE_REFERENCE_COUNT: make(map[string]int),
}

// All global flags
type Flags_t struct {
	is_out_of_tokens        bool
	in_function             bool
	in_function_definition  bool
	compare_types           bool
	expect_return_statement bool
	in_boolean_expression   bool
	return_hit              bool
}

var Flags Flags_t = Flags_t{
	is_out_of_tokens:        false,
	in_function:             false,
	in_function_definition:  false,
	compare_types:           false,
	expect_return_statement: false,
	in_boolean_expression:   false,
	return_hit:              false,
}

// Control variables
type ControlVariables_t struct {
	position                         int
	current_line                     int
	current_section                  string
	most_significant_type            utils.Type
	token_causing_error              string
	current_token                    lexer.Token
	current_line_string              string
	global_has                       int
	current_expected_type            utils.Type
	expected_function_return_type    utils.Type
}

var ControlVariables ControlVariables_t = ControlVariables_t{
	position:                         -1,
	current_line:                     0,
	current_section:                  "",
	most_significant_type:            utils.Type{BaseType: utils.T_INVALID, Dimension: 0, AdditionalType: "", TypeName: ""},
	token_causing_error:              "",
	current_token:                    lexer.Token{},
	current_line_string:              "",
	global_has:                       0,
	current_expected_type:            utils.Type{BaseType: utils.T_INVALID, Dimension: 0, AdditionalType: "", TypeName: ""},
	expected_function_return_type:    utils.Type{BaseType: utils.T_INVALID, Dimension: 0, AdditionalType: "", TypeName: ""},
}

