package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
	"fmt"
)

/*
Use:
	Syntax: use <section-name> as <name>

	Lets you use NamesAndScopes.VARIABLES that are local to another section

	!!! The variable currently_used_setions should be resetted after each block !!!
*/
func EvaluateUse() {
	ParserAdvance()

	var Section string
	var Alias string

	if ControlVariables.current_token.Token_type == lexer.TT_ID && utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.sections) {
		Section = ControlVariables.current_token.Value
		ParserAdvance()

		if ControlVariables.current_token.Token_type == lexer.TT_ID && ControlVariables.current_token.Value == "as" {
			ParserAdvance()

			if ControlVariables.current_token.Token_type == lexer.TT_ID {
				Alias = ControlVariables.current_token.Value
				ParserAdvance()

				NamesAndScopes.currently_used_sections[Alias] = Section
				return
			} else {
				// error no alias specified
				errors.ThrowError(fmt.Sprintf("NoAliasSpecifiedException. After \"as\" in \"use\" an alias was expected but not found at %s", GetLocation()), true)
			}
		} else {
			errors.KeywordExpectedError("as", GetLocation())
		}
	} else if ControlVariables.current_token.Token_type == lexer.TT_BANG {
		// remove variable if in scope2
		ParserAdvance()

		// if the specified section is used
		if ControlVariables.current_token.Token_type == lexer.TT_ID && ItemInStringMap(ControlVariables.current_token.Value, NamesAndScopes.currently_used_sections) {
			delete(NamesAndScopes.currently_used_sections, ControlVariables.current_token.Value)
			ParserAdvance()
		} else {
			errors.CantRemoveSectionThatIsNotInScope(ControlVariables.current_token.Value, GetLocation())
		}
	} else {
		if !utils.StringInArray(ControlVariables.current_token.Value, NamesAndScopes.sections) {
			// error No section specified
			errors.ThrowError(fmt.Sprintf("SpecifiedSectionDoesNotExistException. The section specified in \"use\" (%s) is not valid at %s", ControlVariables.current_token.Value, GetLocation()), true)
		} else {
			// error. No ID
			errors.IDExpectedError(GetLocation(), ControlVariables.current_token.Value)
		}
	}
}
