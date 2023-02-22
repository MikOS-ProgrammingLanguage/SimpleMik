package parser

// UniversalNone implementation
func (universal_none UniversalNone) is_first_class()  {}
func (universal_none UniversalNone) is_second_class() {}

// Root node
func (root_node RootNode) is_first_class() {}

// Struct node
func (struct_node StructNode) is_first_class() {}

// If node
func (if_node IfNode) is_first_class() {}

// while node
func (while_node WhileNode) is_first_class() {}

// return node
func (return_node ReturnNode) is_first_class() {}

// Assignment node
func (assignment_node AssignmentNode) is_first_class() {}

// Array assignment node
func (array_assignment_node ArrayAssignmentNode) is_first_class() {}

// Binary operation node
func (binary_operation_node BinaryOperationNode) is_first_class()  {}
func (Binary_operation_node BinaryOperationNode) is_second_class() {}

// Factor Expression node
func (factor_expression FactorExpression) is_second_class() {}

// Function node
func (function_node FunctionNode) is_first_class() {}

// Function call node
func (function_call_node FunctionCallNode) is_first_class() {}

func (function_call_node FunctionCallNode) is_second_class() {}

// Direct node
func (direct_node DirectNode) is_second_class() {}

// Variable name node
func (variable_name_node VariableNameNode) is_second_class() {}

// reassignment node
func (reassignment_node ReAssignmentNode) is_first_class() {}

// array slice node
func (array_slice_node ArraySliceNode) is_second_class() {}

// array slice reassignment node
func (array_slice_reassignment_node ArraySliceReassignmentNode) is_first_class() {}

// whole array as literal node
func (whole_array_as_literal WholeArrayAsLiteral) is_second_class() {}

// typecast node
func (typecast_node TypeCastNode) is_second_class() {}

// first class function call
func (first_class_function_call FirstClassFunctionCall) is_first_class() {}
