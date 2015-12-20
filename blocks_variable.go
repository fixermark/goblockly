// Evaluators for variable blocks

package goblockly

func VariableSetEvaluator(i *Interpreter, b *Block) Value {
	field := b.FieldWithName("VAR")
	if field == nil {
		i.Fail("No VAR field in variables_set")
		return nilValue
	}
	variableValue := b.SingleBlockValueWithName(i, "VALUE")
	i.Context[field.Value] = i.Evaluate(variableValue)
	return nilValue
}

func VariableGetEvaluator(i *Interpreter, b *Block) Value {
	field := b.FieldWithName("VAR")
	if field == nil {
		i.Fail("No VAR field in variables_get")
		return nilValue
	}
	variableValue, ok := i.Context[field.Value]
	if !ok {
		i.Fail("No variable named '" + field.Value + "'")
		return nilValue
	}
	return variableValue
}
