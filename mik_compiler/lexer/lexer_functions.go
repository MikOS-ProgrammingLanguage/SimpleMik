package lexer

import (
	"MicNewDawn/errors"
	"MicNewDawn/utils"
	"fmt"
	"strings"
)

// returns a formated string with the location of a token
func GetLocation() string {
	return fmt.Sprintf("File: %s Section: %s. Line: %d", current_file, position.current_section.section_name, position.current_section.current_line)
}

// increments the lexer_index and checks for the end of file
func LexerAdvance() {
	if lexer_index+1 >= len(*text) {
		is_end_of_file = true
	} else {
		lexer_index++
		if (*text)[lexer_index] != '\n' {
			utils.RAW_TEXT[current_file][position.current_section.current_line] += string((*text)[lexer_index])
		}
	}
}

// iterates over a string that starts with " and ands with ". Also does assembly function body lexing
func MakeStringToken() Token {
	var return_string string

	if is_in_assembly_function {
		cnt := 0

		// loops over the assembly functions' body
		for !is_end_of_file && (*text)[lexer_index] != '}' {
			current_character := (*text)[lexer_index]

			// check for a new line and if more than one character has been appended
			if current_character == '\n' && cnt > 0 {
				return_string += ";"
				position.current_section.current_line++
			} else if current_character == '\n' && cnt == 0 {
				position.current_section.current_line++
			} else if current_character == '\t' {
				// ignore tabs.
				continue
			} else {
				return_string += string(current_character)
			}
			LexerAdvance()
			cnt++
		}
	} else {
		LexerAdvance() // skip the double quote

		for !is_end_of_file && (*text)[lexer_index] != '"' {
			return_string += string((*text)[lexer_index]) // append character
			LexerAdvance()
		}
		LexerAdvance() // skip the double quote
	}

	is_in_assembly_function = false
	return Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_STRING, Value: return_string}
}

// iterates over a id that starts with a alphabetic character or a '_' and ends with a ' ' or any character not in the alphabet + _
func MakeIdToken(Tokens []Token) Token {
	var id_string string

	ALLOWED_CHAR_STRING := LEGAL_CHARACTERS_IN_ID + LEGAL_CHARACTERS_IN_NUMBERS
	// add a dot to the allowed characters whenever "use" is the previous token
	if len(Tokens) > 0 && (Tokens[len(Tokens)-1].Token_type == TT_ID && Tokens[len(Tokens)-1].Value == "use") {
		ALLOWED_CHAR_STRING += "."
	}

	// loops over the tokens as long as the token is in the alphabet plus numbers plus _
	for !is_end_of_file && strings.Contains(ALLOWED_CHAR_STRING, string((*text)[lexer_index])) {
		id_string += string((*text)[lexer_index])
		LexerAdvance()
	}

	// If the id is mikas a flag is set to true to determine if the lexer is in a assembly function
	if id_string == "mikas" {
		expect_assembly_function = true
	}
	return Token{current_file, position.current_section.section_name, position.current_section.current_line, TT_ID, id_string}
}

// iterates over a number like 10 or 10.5 and returns an appropriate token
func MakeNumberToken(LEGAL_CHARACTERS string) Token {
	var number_string string
	decimal_point_count := 0

	for !is_end_of_file && strings.Contains(LEGAL_CHARACTERS+".", string((*text)[lexer_index])) {
		if (*text)[lexer_index] == '.' {
			// more than one decimal point
			if decimal_point_count == 1 {
				errors.UnexpectedTokenError(GetLocation(), rune((*text)[lexer_index]))
			}
			decimal_point_count++
			number_string += "."
		} else {
			number_string += string((*text)[lexer_index])
		}
		LexerAdvance()
	}

	// returns a float if there was a decimal point
	if decimal_point_count == 1 {
		return Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_FLOAT, Value: number_string}
	} else {
		return Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_INT, Value: number_string}
	}
}
