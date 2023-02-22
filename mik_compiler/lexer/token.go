package lexer

// Token Types (TT)
// All token types. A valid program can only be split in those tokens. Otherwise it's invalid
const (
	// Place holder
	TT_PLACE_HOLDER int8 = iota //

	// Bitwise operators
	TT_BIT_AND        int8 = iota // b&
	TT_BIT_OR         int8 = iota // b|
	TT_BIT_XOR        int8 = iota // ^
	TT_BIT_NOT        int8 = iota // b!
	TT_BIT_SHIT_LEFT  int8 = iota // <<
	TT_BIT_SHIT_RIGHT int8 = iota // >>

	// boolean operators
	TT_LESS_THAN     int8 = iota // <
	TT_GREATER_THAN  int8 = iota // >
	TT_LESS_EQUAL    int8 = iota // <=
	TT_GREATER_EQUAL int8 = iota // >=
	TT_AND           int8 = iota // &&
	TT_OR            int8 = iota // ||
	TT_BANG          int8 = iota // !
	TT_NOT_EQUALS    int8 = iota // !=
	TT_EQUALS        int8 = iota // ==

	// Mathematical operators
	TT_PLUS     int8 = iota // +
	TT_MINUS    int8 = iota // -
	TT_ASTERISK int8 = iota // *
	TT_SLASH    int8 = iota // /
	TT_MODULO   int8 = iota // %

	// control operators
	TT_COMMA              int8 = iota // ,
	TT_SEMICOLON          int8 = iota // ;

	// Assign/Access operators
	TT_ARROW          int8 = iota // ->
	TT_ASSIGNMENT     int8 = iota // =
	TT_PLUS_EQUALS    int8 = iota // +=
	TT_MINUS_EQUALS   int8 = iota // -=
	TT_TIMES_EQUALS   int8 = iota // *=
	TT_DIVIDED_EQUALS int8 = iota // /=

	// Integer / Floats
	TT_DOT   int8 = iota // .
	TT_INT   int8 = iota // a number (10)
	TT_FLOAT int8 = iota // a floatingpoint number (10.5)

	// Parentheses '()', '[]', and '{}'
	TT_LEFT_PARENTHESIS  int8 = iota // (
	TT_RIGHT_PARENTHESIS int8 = iota // )
	TT_LEFT_BRACKET      int8 = iota // [
	TT_RIGHT_BRACKET     int8 = iota // ]
	TT_LEFT_BRACE        int8 = iota // {
	TT_RIGHT_BRACE       int8 = iota // }

	// Other tokens
	TT_ID          int8 = iota // a id is a name that is neither in single or double quotes
	TT_STRING      int8 = iota // a string in double quotes ("Hello, World\n")
	TT_CHARACTER   int8 = iota // a character in single quotes ('C')
	TT_END_OF_FILE int8 = iota // the end of the input file

	// Compiler directives (@)
	TT_DEBUG     int8 = iota // debug directive. As of now not defined
	TT_BIG_O     int8 = iota
	TT_ATTRIBUTE int8 = iota
)

// All legal characters that are allowed in a ID
const LEGAL_CHARACTERS_IN_ID string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_" // Numbers are also legal if they are not the first character
const LEGAL_CHARACTERS_IN_NUMBERS string = "0123456789"
const LEGAL_HEXADECIMAL_CHARACTERS string = LEGAL_CHARACTERS_IN_NUMBERS + "ABCDEFabcdef"

// Is used to track the position of the lexer in a section
type SectionPosition struct {
	section_name string
	current_line int
}

// Is used to track the overall position of the lexer
type Position struct {
	sections               []SectionPosition
	previous_section_index int
	current_section        SectionPosition
}

/*
Is the output for each token
	- section_name is the name of the section the token is in (used for error handling)
	- line is the line the token is on (used for error handling)
	- token_type is the type of the token e.g. TT_STRING, etc...
	- value is the literal value of a token so for example the token with the type TT_PLUS has a value of "+"
*/
type Token struct {
	Section_name string
	File string
	Line         int
	Token_type   int8
	Value        string
}
