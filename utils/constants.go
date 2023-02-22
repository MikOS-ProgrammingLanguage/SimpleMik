package utils

// is used as a global to later on index the text at a location as (file | row).
var RAW_TEXT map[string]map[int]string = map[string]map[int]string{
	"": {0:""},
}

var INT_TYPES []string = []string{
	"i8",
	"i16",
	"i32",
	"i64",
}

// an array with all basetypes
var TYPES []string = []string{
	"int",
	"flt",
	"str",
	"chr",
	"bool",
	"cock",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int8",
	"int16",
	"int32",
	"int64",
	"void",
}

var TYPE_SIZES map[string]int = map[string]int{
	"i32": 32,
	"i16": 16,
	"i8":  8,
	"i1":  1,
	"i64": 64,
}

// an array with all variable behavior descriptors
var VARIABLE_BEHAVIOR_DESCRIPTORS []string = []string{
	"immutable",
	"local",
	"global",
	"volatile",
	"__noret__",
}

// Every key in this map conflicts logically with it's value
var VARIABLE_BEHAVIOR_DESCRIPTORS_CONFLICT_MAP map[string]string = map[string]string{
	"local": "global",
}

// an array with all type descriptors
var TYPE_DESCRIPTORS []string = []string{
	"long",
	"unsigned",
}

// Every key in this mapconflicts logically with it's value
var TYPE_DESCRIPTORS_CONFLICT_MAP map[string]string = map[string]string{}

// an array with all instructions
var INSTRUCTIONS []string = []string{
	"int",
	"flt",
	"str",
	"chr",
	"cock",
	"bool",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int8",
	"int16",
	"int32",
	"int64",
	"void",
	"mikf",
	"mikas",
	"struct",
	"estruct",
	"if",
	"else",
	"elif",
	"while",
	"true",
	"false",
	"keep",
	"roll",
}

const (
	INT     int8 = iota
	FLT     int8 = iota
	STR     int8 = iota
	CHR     int8 = iota
	COCK    int8 = iota
	BOOL    int8 = iota
	UINT8   int8 = iota
	UINT16  int8 = iota
	UINT32  int8 = iota
	UINT64  int8 = iota
	INT8    int8 = iota
	INT16   int8 = iota
	INT32   int8 = iota
	INT64   int8 = iota
	VOID    int8 = iota
	MIKF    int8 = iota
	MIKAS   int8 = iota
	STRUCT  int8 = iota
	ESTRUCT int8 = iota
	IF      int8 = iota
	ELSE    int8 = iota
	ELIF    int8 = iota
	WHILE   int8 = iota
)

// an array with all custom types the user creates
var CUSTOM_TYPES []string = []string{}

// a map that maps types to indices
var TYPE_HIERARCHY map[int8]int = map[int8]int{
	T_VOID:  1,
	T_BOOL:  2,
	
	T_CHAR:   3,
	T_INT8:   3,
	T_UINT8:  3,

	T_INT16:  4,
	T_UINT16: 4,

	T_INT:    5,
	T_UINT32: 5,

	T_INT64:  6,
	T_UINT64: 6,

	T_FLOAT: 7,
	T_INVALID: 8,
}

var RESERVED_KEYWORD_CONSTANTS []string = []string{
	"true",
	"false",
}

var RESERVED_KEYWORD_CONSTANTS_TYPES map[string]Type = map[string]Type{
	"true":  {BaseType: T_BOOL, Dimension: 0, AdditionalType: "", TypeName: "bool"},
	"false": {BaseType: T_BOOL, Dimension: 0, AdditionalType: "", TypeName: "bool"},
}

var RESERVED_KEYWORD_CONSTANTS_VALUES map[string]string = map[string]string{
	"true":  "1",
	"false": "0",
}

var VALID_ATTRIBUTES []string = []string{
	"@__packed__",
	"@__noexit__",
}

// attribute aliases
const (
	PACKED int8 = iota
	NOEXIT int8 = iota
)

var ATTRIBUTES map[string]int8 = map[string]int8{
	"@__packed__": PACKED,
	"@__noexit__": NOEXIT,
}

// types to be used in type struct
const (
	T_INVALID int8 = iota

	T_VOID int8 = iota
	T_FLOAT int8 = iota
	T_STRING int8 = iota
	T_BOOL int8 = iota
	T_CHAR int8 = iota
	T_COCK int8 = iota

	T_UINT8 int8 = iota
	T_UINT16 int8 = iota
	T_UINT32 int8 = iota
	T_UINT64 int8 = iota

	T_INT8 int8 = iota
	T_INT16 int8 = iota
	T_INT int8 = iota
	T_INT64 int8 = iota
)

var STRING_TO_TYPE map[string]int8 = map[string]int8 {
	"void": T_VOID,
	"flt": T_FLOAT,
	"str": T_STRING,
	"chr": T_CHAR,
	"cock": T_COCK,
	"bool": T_BOOL,
	"uint8": T_UINT8,
	"uint16": T_UINT16,
	"uint32": T_UINT32,
	"uint64": T_UINT64,
	"int8": T_INT8,
	"int16": T_INT16,
	"int": T_INT,
	"int32": T_INT,
	"int64": T_INT64,
}

var OPERATORS []string = []string {
	"==",
	"!=",
	"!",
	"<",
	">",
	"&&",
	"||",
	"<=",
	">=",
}

type Type struct {
	BaseType int8 // 0 means that the additional type should be used
	Dimension int8 // 0 means that it's not an array
	AdditionalType string // specifies other types, like struct names
	TypeName string // stores the actual string representation of the type for error purposes and so on
}

func BoolTypeConstr() Type {
	return Type{T_BOOL, 0, "", "bool"}
}

func VoidTypeConstr() Type {
	return Type{T_VOID, 0, "", "void"}
}

func StringTypeConstr() Type {
	return Type{T_STRING, 0, "", "string"}
}

func IntTypeConstr() Type {
	return Type{T_INT, 0, "", "int"}
}

func FloatTypeConstr() Type {
	return Type{T_FLOAT, 0, "", "float"}
}

func CharTypeConstr() Type {
	return Type{T_CHAR, 0, "", "char"}
}

func InvalidTypeConstr() Type {
	return Type{T_INVALID, 0, "", "inv"}
}

func Int8TypeConstr() Type {
	return Type{T_INT8, 0, "", "int8"}
}

func Int16TypeConstr() Type {
	return Type{T_INT16, 0, "", "int16"}
}

func Int64TypeConstr() Type {
	return Type{T_INT64, 0, "", "int64"}
}
