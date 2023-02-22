package executablegenerator

type ExpressionReturn_t struct {
	return_value                string
	store_expression_return     bool
	return_value_is_a_literal   bool
	return_value_is_a_reference bool
}
