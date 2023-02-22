package optimizer

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/parser"
	"fmt"
)

var passes int = 0

// depending on which flags were set, add functions to a pipeline
/*
Flag order:
	- optimize empty blocks
	- optimize dead code
*/
func Optimize(optFlags OptimizationFlags, iterations int, AST *parser.RootNode) *parser.RootNode {
	errors.ShowOpts = optFlags.ShowOptimizationsMade
	ASTSize := GetAstSize(AST.Nodes)
	for {
		passes++
		var optimizationsToApply []RemoveOccurrences

		if optFlags.OptimizeEmptyBlocks {
			optimizationsToApply = append(optimizationsToApply, OptimizeEmptyBlocks(AST.Nodes, []int{})...)
		}
		if optFlags.OptimizeDeadCode {
			optimizationsToApply = append(optimizationsToApply, OptimizeDeadCode(AST.Nodes, []int{})...)
		}

		if len(optimizationsToApply) == 0 {
			// quit optimizing cause there's nothing to do
			break
		}
		ApplyOptimizations(optimizationsToApply, AST)
	}
	newASTSize := GetAstSize(AST.Nodes)

	if optFlags.ShowDifferenceInSizeAndBigO {
		var reducedPercentage float32
		if newASTSize != ASTSize {
			reducedPercentage = 100 / float32(ASTSize) * float32(newASTSize)
		} else {
			reducedPercentage = 0
		}
		errors.ThrowInfo(fmt.Sprintf("Optimized size by %.2f%% in %d passes", 100-reducedPercentage, passes), false)
	}
	return AST
}
