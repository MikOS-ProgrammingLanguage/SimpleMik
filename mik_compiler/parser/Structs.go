package parser

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/lexer"
	"MicNewDawn/utils"
	"reflect"
)

func Struct(attrs []int8) FirstClass {
	ParserAdvance()

	for _, val := range attrs {
		switch val {
		case utils.PACKED: {}
		default:
			errors.UnexpectedAttributeError(ReverseStringInt8Map(utils.ATTRIBUTES)[val], GetLocation())
		}
	}

	VariableBehaviorDescriptors := VariableBehaviorDescriptorValidityCheck(GetVariableBehaviorDescriptors())
	Name := CheckNameTaken(GetName())
	is_global := utils.StringInArray("global", VariableBehaviorDescriptors)

	if ControlVariables.current_token.Token_type == lexer.TT_LEFT_BRACE {
		ParserAdvance()
		fields := ParseBody(lexer.TT_RIGHT_BRACE)
		ParserAdvance()

		// check if all fields are allowed and if their bodies are empty
		for _, val := range fields {
			switch reflect.TypeOf(val).Name() {
			case "AssignmentNode":
				{
					temp := val.(AssignmentNode)
					if temp.Variable_body != nil {
						errors.NoInitializedExpressionExpected(GetLocation())
					}
				}
			case "ArrayAssignmentNode":
				{
					// structs need fixed size only
				}
			default:
				// error
				errors.UnexpectedExpressionInStruct(GetLocation())
			}
		}

		// create struct to return
		to_return := StructNode{
			Name:       Name,
			VARIABLES:  fields,
			Attributes: attrs,
			Local:      is_global,
			Global:     is_global,
		}

		// append struct to types
		NamesAndScopes.Struct_names = append(NamesAndScopes.Struct_names, Name)
		NamesAndScopes.STRUCTS[Name] = to_return
		utils.CUSTOM_TYPES = append(utils.CUSTOM_TYPES, Name)

		// return finished struct
		return to_return

	} else {
		errors.OpeningBraceExpected(GetLocation())
	}

	return nil
}
