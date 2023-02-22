package parser

// Merges a map of string:FirstClass into another map of the same type
func MergeFirstClassMaps(a, b map[string]FirstClass) map[string]FirstClass {
	for key, elem := range b {
		a[key] = elem
	}
	return a
}

// returns true if a given string (a) is in a map of type string:FirstClass
func ItemInFirstClassMap(a string, b map[string]FirstClass) bool {
	_, c := b[a]
	return c
}
func ItemInStringMap(a string, b map[string]string) bool {
	_, c := b[a]
	return c
}

func ItemInIntMap(a string, b map[string]int) bool {
	_, c := b[a]
	return c
}

// returns true if a given string(a) is in a map of arrays
func StringInMapOfArrays(a string, b map[string][]FirstClass) bool {
	_, c := b[a]
	return c
}

func ItemInInt8StringMap(a int8, b map[int8]string) bool {
	_, c := b[a]
	return c
}

func ReverseStringInt8Map(a map[string]int8) map[int8]string {
	n := make(map[int8]string, len(a))
	for k, v := range a {
		n[v] = k
	}
	return n
}

func ReverseStringMap(b map[string]string) map[string]string {
	n := make(map[string]string, len(b))
	for k, v := range b {
		n[v] = k
	}
	return n
}
