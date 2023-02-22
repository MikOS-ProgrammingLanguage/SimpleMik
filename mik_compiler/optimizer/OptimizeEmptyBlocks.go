package optimizer

import (
	"MicNewDawn/mik_compiler/parser"
	"MicNewDawn/utils"
	"reflect"
)

var functionsToRemove []string
var elseBody bool = false
var elseOffset int = 0
var elseIndexOffset int = 0

func OptimizeEmptyBlocks(AST []parser.FirstClass, node []int) []RemoveOccurrences {
	var offset int = 0
	var indexOffset int = 0
	if !elseBody {
		node = append(node, 0)
	} else {
		indexOffset = elseIndexOffset
		offset = elseOffset
	}
	var occurrences []RemoveOccurrences

	for index, val := range AST {
		switch reflect.TypeOf(val).Name() {
		case "FunctionNode":
			{
				temp := val.(parser.FunctionNode)
				node[len(node)-1] = index + indexOffset - offset
				if temp.CodeBlock == nil {
					occurrences = append(occurrences, RemoveOccurrences{
						OccurrenceType: Function,
						Value:          temp.Function_name,
					})
					functionsToRemove = append(functionsToRemove, temp.Function_name)
					offset++
				} else {
					occurrences = append(occurrences, OptimizeEmptyBlocks(temp.CodeBlock, node)...)
				}
			}
		case "AssemblyFunctionNode":
			{

				temp := val.(parser.AssemblyFunctionNode)
				node[len(node)-1] = index + indexOffset - offset
				if len(temp.Assembly_block) == 0 {
					occurrences = append(occurrences, RemoveOccurrences{
						OccurrenceType: Function,
						Value:          temp.Function_name,
					})
					functionsToRemove = append(functionsToRemove, temp.Function_name)
					offset++
				}
			}
		case "IfNode":
			{
				// also need to parse else node for empty blocks like other if statement
				temp := val.(parser.IfNode)
				node[len(node)-1] = index + indexOffset - offset
				if temp.CodeBlock == nil {
					occurrences = append(occurrences, RemoveOccurrences{
						OccurrenceType: DeadCodeAtLocation,
						Value:          node,
					})
					offset++

					if temp.Else && temp.Else_body != nil {
						elseBody = true
						elseOffset = offset
						elseIndexOffset = index
						occurrences = append(occurrences, OptimizeEmptyBlocks(temp.Else_body, node)...)
						offset = elseOffset
						elseIndexOffset = 0
						elseOffset = 0
						elseBody = false
					}
				} else {
					if temp.Else_body != nil {
						elseBody = true
						elseOffset = offset
						elseIndexOffset = index
						occurrences = append(occurrences, OptimizeEmptyBlocks(temp.Else_body, node)...)
						offset = elseOffset
						elseIndexOffset = 0
						elseOffset = 0
						elseBody = false
					}
					occurrences = append(occurrences, OptimizeEmptyBlocks(temp.CodeBlock, node)...)
				}
			}
		case "WhileNode":
			{
				temp := val.(parser.WhileNode)
				node[len(node)-1] = index + indexOffset - offset
				if temp.CodeBlock == nil {
					node[len(node)-1] = index - offset
					occurrences = append(occurrences, RemoveOccurrences{
						OccurrenceType: DeadCodeAtLocation,
						Value:          node,
					})
					offset++
				} else {
					occurrences = append(occurrences, OptimizeEmptyBlocks(temp.CodeBlock, node)...)
				}
			}
		case "FunctionCallNode":
			{
				temp := val.(parser.FunctionCallNode)
				node[len(node)-1] = index - offset
				if utils.StringInArray(temp.Function_name, functionsToRemove) {
					offset++
				}
			}
		case "FirstClassFunctionCall":
			{
				temp := val.(parser.FirstClassFunctionCall)
				node[len(node)-1] = index + indexOffset - offset
				if utils.StringInArray(temp.Node.(parser.FunctionCallNode).Function_name, functionsToRemove) {
					occurrences = append(occurrences, RemoveOccurrences{
						OccurrenceType: DeadCodeAtLocation,
						Value:          node,
					})
					offset++
				}
			}
		}
	}

	if elseBody {
		elseOffset = offset
	}
	return occurrences
}
