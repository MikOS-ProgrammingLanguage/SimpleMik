package parser

import (
	"MicNewDawn/mik_compiler/lexer"
)

// updates the current token to the next token
func ParserAdvance() {
	ControlVariables.current_section = ControlVariables.current_token.Section_name
	ControlVariables.position++
	if ControlVariables.position < len(tokens) {
		ControlVariables.current_token = tokens[ControlVariables.position]
		if ControlVariables.current_line != ControlVariables.current_token.Line {
			ControlVariables.current_line = ControlVariables.current_token.Line
			ControlVariables.current_line_string = NamesAndScopes.PURE_SOURCE_CODE[ControlVariables.current_line-1]
		}
	} else {
		ControlVariables.current_token = lexer.Token{}
		Flags.is_out_of_tokens = true
	}
}
