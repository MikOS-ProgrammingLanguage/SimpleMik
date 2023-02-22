package parser

import "MicNewDawn/utils"

func Return() FirstClass {
	Flags.return_hit = true
	ParserAdvance()
	if ControlVariables.expected_function_return_type.BaseType == utils.T_VOID {
		return ReturnNode{
			Return_expression: nil,
			Goal_type:         ControlVariables.expected_function_return_type,
		}
	}

	ControlVariables.current_expected_type = ControlVariables.expected_function_return_type
	return_expression := Expression()
	return ReturnNode{
		Return_expression: return_expression,
		Goal_type:         ControlVariables.expected_function_return_type,
	}
}
