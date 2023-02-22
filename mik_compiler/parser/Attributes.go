package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
)

func Attributes() []int8 {
	var attrs []int8
	for !Flags.is_out_of_tokens && ControlVariables.current_token.Token_type == lexer.TT_ATTRIBUTE {
		if utils.StringInArray(ControlVariables.current_token.Value, utils.VALID_ATTRIBUTES) {
			attrs = append(attrs, utils.ATTRIBUTES[ControlVariables.current_token.Value])
			ParserAdvance()
		} else {
			errors.InvalidAttributeError(ControlVariables.current_token.Value, GetLocation())
		}
	}

	return attrs
}
