package executablegenerator

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"fmt"
	"reflect"
)

func GenerateExpression(expression parser.SecondClass, Type utils.Type) ExpressionReturn_t {
	switch reflect.TypeOf(expression).Name() {
	case "BinaryOperationNode":
		{
			operation := expression.(parser.BinaryOperationNode)

			left := GenerateExpression(operation.Left_branch, Type)
			operator := GetOperatorLLVMEncoding(operation.Operator_token, Type)

			if operator != "" && Type.BaseType == utils.T_STRING {
				// operation not defined
				errors.OperationNotDefinedForType(Type.TypeName, operator)
			}

			right := GenerateExpression(operation.Right_branch, Type)

			if operator == "and" {
				CodeAdd(fmt.Sprintf("\t%%%d = icmp eq i1 %s, %s\n", EMULATED_REGISTER_COUNTER, left.return_value, right.return_value), false)
				EMULATED_REGISTER_COUNTER++
			} else if operator == "or" {
				CodeAdd(fmt.Sprintf("\t%%%d = or i1 %s, %s\n", EMULATED_REGISTER_COUNTER, left.return_value, right.return_value), false)
				EMULATED_REGISTER_COUNTER++
			} else {
				CodeAdd(fmt.Sprintf("\t%%%d = %s %s %s, %s\n", EMULATED_REGISTER_COUNTER, operator, GetTypeLLVMEncoding(Type), left.return_value, right.return_value), false)
				EMULATED_REGISTER_COUNTER++
			}
			return ExpressionReturn_t{
				return_value:                fmt.Sprintf("%%%d", EMULATED_REGISTER_COUNTER-1),
				store_expression_return:     left.store_expression_return || right.store_expression_return,
				return_value_is_a_literal:   right.return_value_is_a_literal || left.return_value_is_a_literal,
				return_value_is_a_reference: right.return_value_is_a_reference || left.return_value_is_a_reference,
			}
		}
	case "DirectNode":
		{
			operation := expression.(parser.DirectNode)

			if Type.BaseType == utils.T_FLOAT && operation.Type == lexer.TT_INT && !in_array {
				operation.Value += ".0"
			} else if Type.BaseType == utils.T_STRING && operation.Type == lexer.TT_STRING {
				*CODE = fmt.Sprintf("@.str.%d = private unnamed_addr constant [%d x i8] c\"%s\\00\", align 1\n", EMULATED_REGISTER_STRING_COUNTER, len(operation.Value)+1, operation.Value) + *CODE
				EMULATED_REGISTER_STRING_COUNTER++
				return ExpressionReturn_t{
					return_value:              fmt.Sprintf("getelementptr inbounds ([%d x i8], [%d x i8]* @.str.%d, i64 0, i64 0)", len(operation.Value)+1, len(operation.Value)+1, EMULATED_REGISTER_STRING_COUNTER-1),
					return_value_is_a_literal: true,
				}
			}

			return ExpressionReturn_t{
				return_value:              operation.Value,
				store_expression_return:   true,
				return_value_is_a_literal: true,
			}
		}
	case "VariableNameNode":
		{
			operation := expression.(parser.VariableNameNode)

			if reflect.TypeOf(parser.NamesAndScopes.VARIABLES[operation.Name]).Name() == "ArrayAssignmentNode" {
				in_array = true
			}

			a := "%"

			if operation.Is_global {
				a = "@"
			}

			if Type.BaseType == utils.T_STRING {
				EMULATED_REGISTER_COUNTER++
			} else {
				EMULATED_REGISTER_COUNTER++
				CodeAdd(fmt.Sprintf(
					"\t%%%d = load %s, %s* %s%s, align 4\n",
					EMULATED_REGISTER_COUNTER-1,
					GetTypeLLVMEncoding(operation.Type),
					GetTypeLLVMEncoding(operation.Type),
					a,
					operation.Name,
				), false)
				return ExpressionReturn_t{
					return_value: fmt.Sprintf("%%%d", EMULATED_REGISTER_COUNTER-1),
				}
			}

			// Integer interoperability
			if GetTypeLLVMEncoding(Type) != GetTypeLLVMEncoding(operation.Type) && utils.StringInArray(GetTypeLLVMEncoding(operation.Type), utils.INT_TYPES) && utils.StringInArray(GetTypeLLVMEncoding(Type), utils.INT_TYPES) {
				CodeAdd(fmt.Sprintf(
					"%%%d = %s %s %%%d to %s\n",
					EMULATED_REGISTER_COUNTER,
					GetTypeCastLLVMEncoding(GetTypeLLVMEncoding(operation.Type), GetTypeLLVMEncoding(Type)),
					GetTypeLLVMEncoding(operation.Type),
					EMULATED_REGISTER_COUNTER-1,
					GetTypeLLVMEncoding(Type),
					), false)
					EMULATED_REGISTER_COUNTER++
			}

			in_array = false
			return ExpressionReturn_t{
				return_value: fmt.Sprintf("%%%d", EMULATED_REGISTER_COUNTER-1),
			}
		}
	case "FactorExpression":
		{
			operation := expression.(parser.FactorExpression)
			expression := GenerateExpression(operation.Binary_operation_node, Type)
			return expression
		}
	case "ArraySliceNode":
		{
			in_array = true
			var VariableNamePrefix string
			operation := expression.(parser.ArraySliceNode)

			if operation.Is_global {
				VariableNamePrefix = "@"
			} else {
				VariableNamePrefix = "%"
			}
			var_name := operation.Name

			// Evaluate array length expression
			for idx, i := range operation.Positions {
				dimensions := GenerateArrayDimensionString(operation.Lengths[:len(operation.Lengths)-idx], Type)
				// rework to dereference multiple dimensions
				index_expression_return := GenerateExpression(i, Type)
				var index string
				if index_expression_return.store_expression_return {
					index = index_expression_return.return_value
				} else {
					index = fmt.Sprintf("%%%d", EMULATED_REGISTER_COUNTER-1)
				}
			
				// %1 = get element ptr inbounds. (load the array at index)
				CodeAdd(fmt.Sprintf(
					"\t%%%d = getelementptr inbounds %s, %s* %s%s, i64 0, i64 %s\n",
					EMULATED_REGISTER_COUNTER,
					dimensions,
					dimensions,
					VariableNamePrefix,
					var_name,
					index,
				), false)
				EMULATED_REGISTER_COUNTER++
				VariableNamePrefix = "%"
				var_name = fmt.Sprint(EMULATED_REGISTER_COUNTER-1)
			}

			CodeAdd(fmt.Sprintf("\t%%%d = load %s, %s* %%%d, align 4\n", EMULATED_REGISTER_COUNTER, GetTypeLLVMEncoding(Type), GetTypeLLVMEncoding(Type), EMULATED_REGISTER_COUNTER-1), false)
			EMULATED_REGISTER_COUNTER++

			in_array = false
			return ExpressionReturn_t{
				return_value: fmt.Sprint(EMULATED_REGISTER_COUNTER-1),
			}
		}
	case "WholeArrayAsLiteral":
		{
			operation := expression.(parser.WholeArrayAsLiteral)
			in_array = true

			// Get variable prefix
			var VariableNamePrefix string = "%"
			if parser.NamesAndScopes.VARIABLES[operation.Name].(parser.ArrayAssignmentNode).Is_global {
				VariableNamePrefix = "@"
			}

			// Evaluate array length expression
			dimensions := GenerateArrayDimensionString(operation.Lengths, Type)

			CodeAdd(fmt.Sprintf(
				"\t%%%d = getelementptr inbounds %s, %s* %s%s, i64 0, i64 0\n",
				EMULATED_REGISTER_COUNTER,
				dimensions,
				dimensions,
				VariableNamePrefix,
				operation.Name,
			), false)
			EMULATED_REGISTER_COUNTER++
			
			in_array = false
			return ExpressionReturn_t{
				return_value: fmt.Sprint(EMULATED_REGISTER_COUNTER),
			}
		}
	case "FunctionCallNode":
		{
			operation := expression.(parser.FunctionCallNode)

			var function_type utils.Type
			if reflect.TypeOf(parser.NamesAndScopes.FUNCTIONS[operation.Function_name]).Name() == "FunctionNode" {
				function_type = parser.NamesAndScopes.FUNCTIONS[operation.Function_name].(parser.FunctionNode).Return_type
			}

			var assign_to string
			call_body := GetFunctionCallBody(operation.Called_function_arguments, operation.Function_name)

			if function_type.BaseType == utils.T_VOID {
				assign_to = "   "
			} else {
				assign_to = fmt.Sprintf("%%%d = ", EMULATED_REGISTER_COUNTER)
				EMULATED_REGISTER_COUNTER++
			}

			CodeAdd(fmt.Sprintf(
				"\t%scall %s @%s(%s)\n",
				assign_to,
				GetTypeLLVMEncoding(function_type),
				operation.Function_name,
				call_body,
			), false)

			return ExpressionReturn_t{
				return_value: assign_to[:len(assign_to)-3],
			}
		}
	case "TypeCastNode":
		{
			operation := expression.(parser.TypeCastNode)
			expression_return := GenerateExpression(operation.Expression, operation.From)
			CodeAdd(fmt.Sprintf(
				"\t%%%d = %s %s %s to %s\n",
				EMULATED_REGISTER_COUNTER,
				GetTypeCastLLVMEncoding(GetTypeLLVMEncoding(operation.From), GetTypeLLVMEncoding(operation.To)),
				GetTypeLLVMEncoding(operation.From),
				expression_return.return_value,
				GetTypeLLVMEncoding(operation.To),
			), false)
			EMULATED_REGISTER_COUNTER++

			return ExpressionReturn_t{
				return_value: fmt.Sprintf("%%%d", EMULATED_REGISTER_COUNTER-1),
			}
		}
	}

	return ExpressionReturn_t{}
}
