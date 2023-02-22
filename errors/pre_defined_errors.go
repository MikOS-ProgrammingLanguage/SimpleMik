package errors

import (
	"fmt"
)

// Pre defined functions for error handeling

// Prints "[ERROR] <your_text>" in red
func UnexpectedTokenError(location string, token rune) {
	ThrowError(fmt.Sprintf("UnexpectedTokenError (%c) \n---> %s", token, location), true)
}

func ClosingParenthesisExpected(location string) {
	ThrowError(fmt.Sprintf("ClosingParenthesisError. A closing parenthesis ')' was expected \n---> %s", location), true)
}

func ClosingBracketExpected(location string) {
	ThrowError(fmt.Sprintf("ClosingBracketExpectedError. A closing Bracket ']' was expected \n---> %s", location), true)
}

func OpeningParenthesisExpected(location string) {
	ThrowError(fmt.Sprintf("OpeningParenthesisError. A opening parenthesis '(' was expected \n---> %s", location), true)
}

func YoinkPackageError(package_name string) {
	ThrowError(fmt.Sprintf("YoinkPackageError. Could not locate or open package: %s", package_name), true)
}

func YoinkError(file_name string) {
	ThrowError(fmt.Sprintf("YoinkError. Could not locate or open file: %s", file_name), true)
}

func TypeMismatchError(expected, got, location string) {
	// TODO: rework to display dimensions as well
	ThrowError(fmt.Sprintf("TypeMismatchError. Expected: %s. But got: %s \n---> %s", expected, got, location), true)
}

func OpeningBraceExpected(location string) {
	ThrowError(fmt.Sprintf("OpeningBraceExpectedError. Expected an opening brace\n---> %s", location), true)
}

func ReturnStatementExpected(location string) {
	ThrowError(fmt.Sprintf("ReturnStatementExpectedError. Expected a return statement\n---> %s", location), true)
}

func PointerAccessError(max_pointers, actual_pointer int, location string) {
	ThrowError(fmt.Sprintf("PointerAccessError. Maximal allowed pointers: %d. But got: %d pointers \n---> %s", max_pointers, actual_pointer, location), true)
}

func PointerMismatchError(location string) {
	ThrowError(fmt.Sprintf("PointerMismatchError. %s", location), true)
}

func ReferenceError(name, location string) {
	ThrowError(fmt.Sprintf("ReferenceError. The referenced variable, function or struct \"%s\" is not defined \n---> %s", name, location), true)
}

func ArgumentExpectedError(location string) {
	ThrowError(fmt.Sprintf("ExpressionExpectedError. After a comma in a function call or struct, an expression is expected \n---> %s", location), true)
}

func ToFewArgumentsException(location, function_name string, expected, got int) {
	ThrowError(fmt.Sprintf("ToFewArgumentsException. In a call to function \"%s\" %d argument(s) were expected but %d were given \n---> %s", function_name, expected, got, location), true)
}

func ToManyArgumentsException(location, function_name string, expected, got int) {
	ThrowError(fmt.Sprintf("ToManyArgumentsException. In a call to function \"%s\" %d argument(s) were expected but %d were given \n---> %s", function_name, expected, got, location), true)
}

func ConflictingTypesError(location, a, b string) {
	ThrowError(fmt.Sprintf("ConflictingTypesError. The type \"%s\" is not compatible with \"%s\" \n---> %s", a, b, location), true)
}

func TypeExpectedError(location string) {
	ThrowError(fmt.Sprintf("TypeExpectedError. A type was expected in an assignment \n---> %s", location), true)
}

func IDExpectedError(location, got string) {
	ThrowError(fmt.Sprintf("IDExpectedError. An ID was expected but \"%s\" was found \n---> %s", got, location), true)
}

func NameAlreadyTakenError(location, name string) {
	ThrowError(fmt.Sprintf("NameAlreadyTakenError. The name \"%s\" is Already taken \n---> %s", name, location), true)
}

func LocalScopeError(location, variable_name, section1, section2 string) {
	ThrowError(fmt.Sprintf("ScopeError. Variable or function \"%s\" is local to \"%s\" but was referenced from \"%s\" \n---> %s", variable_name, section1, section2, location), true)
}

func KeywordExpectedError(keyword, location string) {
	ThrowError(fmt.Sprintf("KeywordExpectedError. Keyword: \"%s\" was expected but not found \n---> %s", keyword, location), true)
}

func OperationNotDefinedForType(type_, operation string) {
	ThrowError(fmt.Sprintf("OperationNotDefinedForType. Operation: \"%s\" is not defined for \"%s\"", operation, type_), true)
}

func ExpectedToken(token, location string) {
	ThrowError(fmt.Sprintf("ExpectedTokenError. Token \"%s\" \n---> %s", token, location), true)
}

func VariableDoesNotExist(variable, location string) {
	ThrowError(fmt.Sprintf("VariableDoesNotExist. The variable \"%s\" does not exist but was referenced \n---> %s", variable, location), true)
}

func FunctionArgumentExpectedError(location string) {
	ThrowError(fmt.Sprintf("ArgumentExpectedError. Expected another Argument after a \",\" in a function assignment \n---> %s", location), true)
}

func UnexpectedReturnType(location, got string) {
	ThrowError(fmt.Sprintf("UnexpectedReturnType. Expected a type as a function but got \"%s\" \n---> %s", got, location), true)
}

func ASMBlockExpected(location string) {
	ThrowError(fmt.Sprintf("ASMBlockExpected. Expected a asm block in assembly function \n---> %s", location), true)
}

func CantRemoveSectionThatIsNotInScope(name, location string) {
	ThrowError(fmt.Sprintf("Can'tRemoveSectionThatIsNotInScope. Failed to remove section %s from scope as it is not in scope \n---> %s", name, location), true)
}

func CantNestFunctions(location string) {
	ThrowError(fmt.Sprintf("CanNotNestFunctions. Tried nesting functions \n---> %s", location), true)
}

func IsNotAssignableError(of_type, location string) {
	ThrowError(fmt.Sprintf("IsNotAssignableError. Expression of type: \"%s\" is not assignable \n---> %s", of_type, location), true)
}

func NoInitializedExpressionExpected(location string) {
	ThrowError(fmt.Sprintf("InitializedVariableNotExpected. A variable that is supposed to be uninitialized was initialized \n---> %s", location), true)
}

func UnexpectedExpressionInStruct(location string) {
	ThrowError(fmt.Sprintf("UnexpectedExpressionInStruct. An expression in a struct is not allowed. Allowed are 1. Assignments; 2. Array assignments. \n---> %s", location), true)
}

func InvalidAttributeError(attribute, location string) {
	ThrowError(fmt.Sprintf("InvalidAttributeError. %s is not a valid attribute \n---> %s", attribute, location), true)
}

func UnexpectedAttributeError(attribute, location string) {
	// convert attribute to name by reversing the map
	ThrowError(fmt.Sprintf("UnexpectedAttributeError. Got an unexpected argument. namingly %s \n---> %s", attribute, location), true)
}

func UnusedAttributesError(length int, attributes string, location string) {
	ThrowError(fmt.Sprintf("UnusedAttributeError. Found %d unused attributes. Namingly: %s) \n---> %s", length, attributes[:len(attributes)-2], location), true)
}

func ClosingBraceExpected(location string) {
	ThrowError(fmt.Sprintf("Closing brace expected!\n---> %s", location), true)
}

// Mip (mik package manager)
func ToFewArgumentError() {
	ThrowError("MIP was called with to few arguments. Run \"mip help\" for help!", true)
}

func UnknownArgumentError(argument string) {
	ThrowError(fmt.Sprintf("Got unexpected argument \"%s\"", argument), true)
}

func PermissionDenied() {
	ThrowError("Permission denied. Execute as root!", true)
}

func PackageAlreadyExistsError(package_name string) {
	ThrowError(fmt.Sprintf("Package %s Already exists!", package_name), true)
}

func PackageNotFoundError(package_name string) {
	ThrowError(fmt.Sprintf("Package %s not installed!", package_name), false)
}

func NotAPackageError(command string) {
	ThrowError(fmt.Sprintf("The directory from which mip was called is not a mik project and can hence not run the \"%s\" command", command), true)
}

// Config
func CouldNotOpenOrLocateFile(file string) {
	ThrowError(fmt.Sprintf("ConfigError. Could not locate or open file: %s", file), true)
}
