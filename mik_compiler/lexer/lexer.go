package lexer

import (
	"MicNewDawn/errors"
	"MicNewDawn/utils"
	"fmt"
	"strconv"
	"strings"
)

// is used to track the current file
var current_file string

// is used to track is lexer will be in an assembly function soon
var expect_assembly_function bool = false

// is used to track if lexer is in an assembly function
var is_in_assembly_function bool = false

// is used to access tokens in a string
var lexer_index int = 0

// is the text that should be lexed
var text *string

// indicates wether the end of the text is reached or not
var is_end_of_file bool = false

// is the position. More information in token.go in struct Position
var position Position
var sections []string

func Lex(file_name *string, file_text *string, ignore_sections []string) ([]Token, []string) {
	var Tokens []Token

	text = file_text
	position = Position{sections: []SectionPosition{}, previous_section_index: 0, current_section: SectionPosition{section_name: *file_name, current_line: 1}}

	for !is_end_of_file {
		if utils.StringInArray(position.current_section.section_name, ignore_sections) && (*text)[lexer_index] != '@' {
			LexerAdvance()
			continue
		}

		// create a special string in case of an assembly function body
		if is_in_assembly_function {
			Tokens = append(Tokens, MakeStringToken())
			continue
		}

		switch (*text)[lexer_index] {
		// Skip whitespaces and new lines
		case ' ':
			{
				LexerAdvance()
			}
		case '\t':
			{
				LexerAdvance()
			}
		case '\n':
			{
				position.current_section.current_line++
				LexerAdvance()
			}

		// Hexadecimal and binary numbers
		case '0':
			{
				LexerAdvance()

				if (*text)[lexer_index] == 'b' && ((*text)[lexer_index+1] == '0' || (*text)[lexer_index+1] == '1') {
					LexerAdvance()
					Tokens = append(Tokens, MakeNumberToken("01"))
					integer_value, _ := strconv.ParseInt(Tokens[len(Tokens)-1].Value, 2, 0)
					Tokens[len(Tokens)-1].Value = fmt.Sprint(integer_value)
				} else if (*text)[lexer_index] == 'x' && strings.Contains(LEGAL_HEXADECIMAL_CHARACTERS, string((*text)[lexer_index+1])) {
					LexerAdvance()
					Tokens = append(Tokens, MakeNumberToken(LEGAL_HEXADECIMAL_CHARACTERS))
					integer_value, _ := strconv.ParseInt(Tokens[len(Tokens)-1].Value, 16, 0)
					Tokens[len(Tokens)-1].Value = fmt.Sprint(integer_value)
				} else {
					lexer_index--
					Tokens = append(Tokens, MakeNumberToken(LEGAL_CHARACTERS_IN_NUMBERS))
				}
			}

		// Other operators
		case '<':
			{
				LexerAdvance()

				// if token is <= or << or just <
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_LESS_EQUAL, Value: "<="})
					LexerAdvance()
				} else if (*text)[lexer_index] == '<' { // if shift left
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_BIT_SHIT_LEFT, Value: "<<"})
					LexerAdvance()
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_LESS_THAN, Value: "<"})
				}
			}
		case '>':
			{
				LexerAdvance()

				// if token is >= or >> or just >
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_GREATER_EQUAL, Value: ">="})
					LexerAdvance()
				} else if (*text)[lexer_index] == '<' { // if shift left
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_BIT_SHIT_RIGHT, Value: ">>"})
					LexerAdvance()
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_GREATER_THAN, Value: ">"})
				}
			}
		case '=':
			{
				LexerAdvance()

				// if token is == or just =
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_EQUALS, Value: "=="})
					LexerAdvance()
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_ASSIGNMENT, Value: "="})
				}
			}
		case '!':
			{
				LexerAdvance()

				// if token is != or just !
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_NOT_EQUALS, Value: "!="})
					LexerAdvance()
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_BANG, Value: "!"})
				}
			}
		case '|':
			{
				LexerAdvance()

				// if token is ||
				if (*text)[lexer_index] == '|' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_OR, Value: "||"})
					LexerAdvance()
				} else {
					errors.UnexpectedTokenError(GetLocation(), rune((*text)[lexer_index]))
				}
			}
		case '+':
			{
				LexerAdvance()

				// if token is += or just +
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_PLUS_EQUALS, Value: "+="})
					LexerAdvance()
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_PLUS, Value: "+"})
				}
			}
		case '-':
			{
				LexerAdvance()

				// if token is -= or -> or just -
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_MINUS_EQUALS, Value: "-="})
					LexerAdvance()
				} else if (*text)[lexer_index] == '>' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_ARROW, Value: "->"})
					LexerAdvance()
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_MINUS, Value: "-"})
				}
			}
		case '*':
			{
				LexerAdvance()

				// if token is *= or just *
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_TIMES_EQUALS, Value: "*="})
					LexerAdvance()
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_ASTERISK, Value: "*"})
				}
			}
		case '/':
			{
				LexerAdvance()

				// if token is /= or // or /* or just /
				if (*text)[lexer_index] == '=' {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_DIVIDED_EQUALS, Value: "/*"})
					LexerAdvance()
				} else if (*text)[lexer_index] == '/' {
					for (*text)[lexer_index] != '\n' && !is_end_of_file {
						LexerAdvance()
					}
				} else if (*text)[lexer_index] == '*' {
					LexerAdvance()

					// skip everything until '*/' which closes a multiline comment
					for !is_end_of_file {
						if (*text)[lexer_index] == '*' {
							LexerAdvance()
							if (*text)[lexer_index] == '/' {
								LexerAdvance()
								break
							}
						} else if (*text)[lexer_index] == '\n' {
							position.current_section.current_line++
						}
						LexerAdvance()
					}

					if is_end_of_file {
						// comment was never closed :(
						errors.ThrowError(fmt.Sprintf("A multiline comment was started but never closed (%s)", GetLocation()), true)
					}
				} else {
					Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_SLASH, Value: "/"})
				}
			}
		case '%':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_MODULO, Value: "%"})
			}
		case ',':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_COMMA, Value: ","})
			}
		case ';':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_SEMICOLON, Value: ";"})
			}
		case '.':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_DOT, Value: "."})
			}

		// Parentheses
		case '(':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_LEFT_PARENTHESIS, Value: "("})
			}
		case ')':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_RIGHT_PARENTHESIS, Value: ")"})
			}
		case '[':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_LEFT_BRACKET, Value: "["})
			}
		case ']':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_RIGHT_BRACKET, Value: "]"})
			}
		case '{':
			{
				// if an assembly function is expected, it's clear now that the lexer is in the body after this
				if expect_assembly_function {
					is_in_assembly_function = true
					expect_assembly_function = false
				}

				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_LEFT_BRACE, Value: "{"})
			}
		case '}':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_RIGHT_BRACE, Value: "}"})
			}

		// @directives and xor
		case '^':
			{
				LexerAdvance()
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_BIT_XOR, Value: "^"})
			}
		case '@':
			{
				LexerAdvance()
				directive_name := MakeIdToken(Tokens).Value

				switch directive_name {
				case "section":
					{
						if (*text)[lexer_index] == '(' {
							LexerAdvance()

							if (*text)[lexer_index] == '"' {
								section_name := MakeStringToken().Value

								if (*text)[lexer_index] == ')' {
									sec_len := position.current_section.current_line
									position.sections = append(position.sections, position.current_section)
									position.previous_section_index++
									sections = append(sections, section_name)
									position.current_section = SectionPosition{section_name: section_name, current_line: sec_len}

									LexerAdvance()
								} else {
									errors.ClosingParenthesisExpected(GetLocation())
								}
							} else {
								errors.ThrowError(fmt.Sprintf("StringExpectedError A string (\"\") was expected at %s after '@section(' but not found", GetLocation()), true)
							}
						} else {
							errors.OpeningParenthesisExpected(GetLocation())
						}
					}
				case "secend":
					{
						if position.previous_section_index > 0 {
							position.current_section = position.sections[position.previous_section_index-1]
							position.previous_section_index--
						}
					}
				case "file":
					if (*text)[lexer_index] == '(' {
						LexerAdvance()

						if (*text)[lexer_index] == '"' {
							file_name := MakeStringToken().Value

							if (*text)[lexer_index] == ')' {
								current_file = file_name
								position.current_section.current_line = -1
								utils.RAW_TEXT[file_name] = map[int]string{1:""}
								LexerAdvance()
							} else {
								errors.ClosingParenthesisExpected(GetLocation())
							}
						} else {
							errors.ThrowError(fmt.Sprintf("StringExpectedError A string (\"\") was expected at %s after '@section(' but not found", GetLocation()), true)
						}
					} else {
						errors.OpeningParenthesisExpected(GetLocation())
					}
				case "O":
					if (*text)[lexer_index] == '(' {
						LexerAdvance()

						if (*text)[lexer_index] == '"' {
							O_of := MakeStringToken().Value

							if (*text)[lexer_index] == ')' {
								Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_BIG_O, Value: O_of})
								LexerAdvance()
							} else {
								errors.ClosingParenthesisExpected(GetLocation())
							}
						} else {
							errors.ThrowError(fmt.Sprintf("StringExpectedError A string (\"\") was expected at %s after '@section(' but not found", GetLocation()), true)
						}
					} else {
						errors.OpeningParenthesisExpected(GetLocation())
					}
				default:
					{
						Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, File: current_file, Line: position.current_section.current_line, Token_type: TT_ATTRIBUTE, Value: "@"+directive_name})
					}
				}
			}

		// Strings and character ("", '')
		case '"':
			{
				Tokens = append(Tokens, MakeStringToken())
			}
		case '\'':
			{
				Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, Line: position.current_section.current_line, Token_type: TT_CHARACTER, Value: string((*text)[lexer_index+1])})
				LexerAdvance()
				LexerAdvance()
				LexerAdvance()
			}

		// Integers/Floats, and ID's
		default:
			{
				if strings.Contains(LEGAL_CHARACTERS_IN_ID, string((*text)[lexer_index])) {
					// bitwise logic
					if (*text)[lexer_index] == 'b' {
						LexerAdvance()

						if (*text)[lexer_index] == '&' {
							LexerAdvance()
							Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, Line: position.current_section.current_line, Token_type: TT_BIT_AND, Value: "b&"})
						} else if (*text)[lexer_index] == '|' {
							LexerAdvance()
							Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, Line: position.current_section.current_line, Token_type: TT_BIT_OR, Value: "b|"})
						} else if (*text)[lexer_index] == '!' {
							LexerAdvance()
							Tokens = append(Tokens, Token{Section_name: position.current_section.section_name, Line: position.current_section.current_line, Token_type: TT_BIT_NOT, Value: "b!"})
						} else {
							lexer_index--
							Tokens = append(Tokens, MakeIdToken(Tokens))
						}
					} else {
						Tokens = append(Tokens, MakeIdToken(Tokens))
					}
				} else if strings.Contains(LEGAL_CHARACTERS_IN_NUMBERS, string((*text)[lexer_index])) {
					Tokens = append(Tokens, MakeNumberToken(LEGAL_CHARACTERS_IN_NUMBERS))
				} else {
					errors.ThrowError(fmt.Sprintf("IllegalTokenError '%c' was not expected here: %s", (*text)[lexer_index], GetLocation()), true)
				}
			}
		}
	}
	return Tokens, sections
}
