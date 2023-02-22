package parser

import "MicNewDawn/utils"

// Structs that are first class are things like NamesAndScopes.VARIABLES or NamesAndScopes.FUNCTIONS. Though not their body's
type FirstClass interface {
	is_first_class()
}

// Second class structs are body's of first class structs
type SecondClass interface {
	is_second_class()
}

// Is a placeholder that effectively does nothing
type UniversalNone struct{}

// The root node is the entrypoint of the program. Each branch stems from the root node.
type RootNode struct {
	Nodes []FirstClass
}

// for if statements. Can also be else if (FirstClass)
type StructNode struct {
	Name       string
	VARIABLES  []FirstClass
	Attributes []int8
	Local      bool
	Global     bool
}

type IfNode struct {
	Boolean_statement SecondClass
	CodeBlock         []FirstClass
	Type              utils.Type
	Else              bool
	Else_body         []FirstClass
}

// for while loops. (FirstClass)
type WhileNode struct {
	Boolean_statement SecondClass
	CodeBlock         []FirstClass
	Type              utils.Type
}

// for the return statement. (FirstClass)
type ReturnNode struct {
	Return_expression SecondClass
	Goal_type         utils.Type
}

type ReAssignmentNode struct {
	Type               utils.Type
	Is_global          bool
	Variable_name      string
	Body               SecondClass
	Section            string
}

// used for NamesAndScopes.VARIABLES being assigned (FirstClass)
type AssignmentNode struct {
	Variable_behavior_descriptors []string
	Type_descriptors              []string
	Type                          utils.Type
	Is_global                     bool
	Variable_name                 string
	Variable_body                 SecondClass
	Section                       string
	Bounds						  Bound
}

type Bound struct {
	BKeep 						  bool
	BRoll						  bool
	Upper						  SecondClass
	Lower						  SecondClass
}

type ArrayAssignmentNode struct {
	Variable_behavior_descriptors []string
	Type_descriptors              []string
	Type                          utils.Type
	Dimensions int8
	Is_global                     bool
	Variable_name                 string
	Array_length                  []SecondClass
	Section                       string
}

type ArraySliceReassignmentNode struct {
	Type               utils.Type
	Is_global          bool
	Variable_name      string
	Body               SecondClass
	Indices            []SecondClass
	Lengths            []SecondClass
	Dimensions 		   int8
	Section            string
}

type BinaryOperationNode struct {
	Left_branch    SecondClass
	Operator_token string
	Right_branch   SecondClass
}

type FactorExpression struct {
	Binary_operation_node SecondClass
	not                   bool
	minus_prefix          bool
	bitwise_not           bool
}

type FunctionNode struct {
	Declared                      bool
	Function_name                 string
	Arguments                     FunctionAssignArguments
	Return_type                   utils.Type
	CodeBlock                     []FirstClass
	Variable_behavior_descriptors []string
	Section                       string
	VARIABLES                     map[string]FirstClass
	FUNCTIONS                     map[string]FirstClass
	Structs                       map[string]FirstClass
}

type FirstClassFunctionCall struct {
	Node SecondClass
}

type FunctionAssignArguments struct {
	Arguments    []FirstClass
	TypeArray    []utils.Type
}

type FunctionCallNode struct {
	Function_name             string
	Called_function_arguments []SecondClass
	Is_negated                bool
	Bitwise_not               bool
}

type DirectNode struct {
	Type        int8
	Value       string
	Is_negated  bool
	Bitwise_not bool
}

type VariableNameNode struct {
	Name                      string
	Type                      utils.Type
	Not                       bool
	Is_negated                bool
	Bitwise_not               bool
	Is_global                 bool
}

type ArraySliceNode struct {
	Name            string
	Positions       []SecondClass
	Lengths         []SecondClass
	Not             bool
	Is_negated      bool
	Bitwise_not     bool
	Is_global       bool
}

type WholeArrayAsLiteral struct {
	Name            string
	Lengths         []SecondClass
	Dimensions int8
	Not             bool
	Is_negated      bool
	Bitwise_not     bool
	Is_global       bool
}

type TypeCastNode struct {
	Expression    SecondClass
	From          utils.Type
	To            utils.Type
}
