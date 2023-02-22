package errors

import "fmt"

var ShowOpts bool

func UnusedVariableWarning(name string) {
	ThrowWarning(fmt.Sprintf("UnusedVariable. The variable \"%s\" is not used!", name), false)
}

func UnusedFunctionWarning(name string) {
	ThrowWarning(fmt.Sprintf("UnusedFunction. The function \"%s\" is not called", name), false)
}

func ManualPointerModification(location string) {
	if WALL_FLAG {
		ThrowWarning(fmt.Sprintf("Modifying the address of a pointer manually is dangerous!\n---> %s", location), false)
	}
}

func OptimizedFunctionCallWarning(name string, passes int) {
	if ShowOpts {
		ThrowWarning(fmt.Sprintf("In pass %d Optimized: Call to function       -> %s()", passes, name), false)
	}
}

func OptimizedFunctionDefinition(name string, passes int) {
	if ShowOpts {
		ThrowWarning(fmt.Sprintf("In pass %d Optimized: Definition of function -> %s", passes, name), false)
	}
}

func OptimizedAtLocationWarning(name string, passes int) {
	if ShowOpts {
		ThrowWarning(fmt.Sprintf("In pass %d Optimized: Code at location       -> %s", passes, name), false)
	}
}