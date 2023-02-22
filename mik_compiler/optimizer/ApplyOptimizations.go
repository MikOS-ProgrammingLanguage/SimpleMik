package optimizer

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/parser"
	"fmt"
	"reflect"
)

func ApplyOptimizations(x []RemoveOccurrences, AST *parser.RootNode) {
	fmt.Println(x)
	for _, value := range x {
		//fmt.Println(AST)
		switch value.OccurrenceType {
		case Function:
			{
				*AST = RemoveFunction(value.Value.(string), AST)
			}
		case DeadCodeAtLocation:
			{
				*AST = RemoveDeadCodeAtLocation(value.Value.([]int), AST.Nodes, 0)
			}
		}
	}
}

// only needs to look out for some specific node.
// thats because a function that has an empty code block can only be of the signature "any -> .. -> void"
// and hence this function can only be called freely
func RemoveFunction(name string, AST *parser.RootNode) parser.RootNode {
	var newAST parser.RootNode
	for _, node := range AST.Nodes {
		switch reflect.TypeOf(node).Name() {
		case "FunctionNode":
			{
				temp := node.(parser.FunctionNode)
				if temp.Function_name == name {
					errors.OptimizedFunctionDefinition(temp.Function_name, passes)
				} else {
					newAST.Nodes = append(newAST.Nodes, temp)
				}
			}
		case "AssemblyFunctionNode":
			{
				functionName := node.(parser.AssemblyFunctionNode).Function_name
				if functionName == name {
					errors.OptimizedFunctionDefinition(functionName, passes)
				} else {
					newAST.Nodes = append(newAST.Nodes, node)
				}
			}
		}
	}

	return newAST
}

// optimizing dead code should be done directly. Not in later passes as it has no effect if removed directly
// use given offset to access variables at the correct location. Needs to be incremented each time something is removed
// given offset has to be separate for each layer of depth
func RemoveDeadCodeAtLocation(path []int, AST []parser.FirstClass, offsetIndex int) parser.RootNode {
	if len(path)-1 == offsetIndex {
		var newAST []parser.FirstClass

		value := path[offsetIndex]
		typeName := reflect.TypeOf(AST[value]).Name()

		if reflect.TypeOf(AST[value]).Name() == "IfNode" && AST[value].(parser.IfNode).Else_body != nil {
			newAST = append(AST[:value], AST[value].(parser.IfNode).Else_body...)
			//fmt.Println(newAST)
			newAST = append(newAST, AST[value:]...)
			//fmt.Println(AST[value:])
		} else {
			// increment at location
			newAST = append(AST[:value], AST[value+1:]...)
		}
		errors.OptimizedAtLocationWarning(typeName, passes)
		return parser.RootNode{Nodes: newAST}
	} else {
		currentIndexedVariable := AST[path[offsetIndex]]

		switch reflect.TypeOf(currentIndexedVariable).Name() {
		case "FunctionNode":
			{
				newFunction := AST[path[offsetIndex]].(parser.FunctionNode)
				newFunction.CodeBlock = RemoveDeadCodeAtLocation(path, currentIndexedVariable.(parser.FunctionNode).CodeBlock, offsetIndex+1).Nodes
				newAST := append(AST[:path[offsetIndex]], newFunction)
				newAST = append(newAST, AST[path[offsetIndex]+1:]...)
				return parser.RootNode{Nodes: newAST}
			}
		case "IfNode":
			{
				newIf := AST[path[offsetIndex]].(parser.IfNode)
				newIf.CodeBlock = RemoveDeadCodeAtLocation(path, currentIndexedVariable.(parser.IfNode).CodeBlock, offsetIndex+1).Nodes
				newAST := append(AST[:path[offsetIndex]], newIf)
				newAST = append(newAST, AST[path[offsetIndex]+1:]...)
				return parser.RootNode{Nodes: newAST}
			}
		case "WhileNode":
			{
				newWhile := AST[path[offsetIndex]].(parser.WhileNode)
				newWhile.CodeBlock = RemoveDeadCodeAtLocation(path, currentIndexedVariable.(parser.WhileNode).CodeBlock, offsetIndex+1).Nodes
				newAST := append(AST[:path[offsetIndex]], newWhile)
				newAST = append(newAST, AST[path[offsetIndex]+1:]...)
				return parser.RootNode{Nodes: newAST}
			}
		}

		return parser.RootNode{Nodes: AST}
	}
}

func TraverseRemoveFunction(name string, AST []parser.FirstClass) []parser.FirstClass {
	var returnArray []parser.FirstClass
	for _, node := range AST {
		if (reflect.TypeOf(node).Name()) == "FirstClassFunctionCall" {
			if node.(parser.FirstClassFunctionCall).Node.(parser.FunctionCallNode).Function_name == name {
				errors.OptimizedFunctionCallWarning(node.(parser.FirstClassFunctionCall).Node.(parser.FunctionCallNode).Function_name, passes)
				continue
			}
		}
		returnArray = append(returnArray, node)
	}
	return returnArray
}
