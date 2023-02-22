package optimizer

import (
	"MicNewDawn/mik_compiler/parser"
	"reflect"
)

const (
	Function           int = iota
	Variable           int = iota
	DeadCodeAtLocation int = iota // should be of the pattern idxInRoot:idxInSubNode1:...
)

type RemoveOccurrences struct {
	OccurrenceType int
	Value          interface{}
}

type OptimizationFlags struct {
	ShowOptimizationsMade       bool `json:"show"`
	OptimizeEmptyBlocks         bool `json:"optimize-empty-blocks"`
	OptimizeDeadCode            bool `json:"optimize-dead-code"`
	ShowDifferenceInSizeAndBigO bool `json:"show-size-difference"`
}

func GetAstSize(AST []parser.FirstClass) int {
	if AST == nil {
		return 0
	}
	returnValue := 0
	for _, value := range AST {
		switch reflect.TypeOf(value).Name() {
		case "FunctionNode":
			{
				returnValue += GetAstSize(value.(parser.FunctionNode).CodeBlock) + 1
			}
		case "IfNode":
			{
				returnValue += GetAstSize(value.(parser.IfNode).CodeBlock) + 1
			}
		case "WhileNode":
			{
				returnValue += GetAstSize(value.(parser.WhileNode).CodeBlock) + 1
			}
		default:
			{
				returnValue++
			}
		}
	}

	return returnValue
}
