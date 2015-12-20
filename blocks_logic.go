// Evaluators for logic blocks

package blockly

func LogicCompareEvaluator(i *Interpreter, b *Block) Value {
	aBlock := b.SingleBlockValueWithName(i, "A")
	bBlock := b.SingleBlockValueWithName(i, "B")
	opField := b.FieldWithName("OP")
	if opField == nil {
		i.Fail("Missing operator in logic operation block.")
		return nilValue
	}
	aValue := i.Evaluate(aBlock)
	bValue := i.Evaluate(bBlock)

	var result bool
	switch opField.Value {
	case "EQ":
		result = aValue.Equals(i, bValue)
	case "NEQ":
		result = !aValue.Equals(i, bValue)
	case "LT":
		result = aValue.IsLessThan(i, bValue)
	case "LTE":
		result = aValue.IsLessThan(i, bValue) || aValue.Equals(i, bValue)
	case "GT":
		result = !aValue.IsLessThan(i, bValue) && !aValue.Equals(i, bValue)
	case "GTE":
		result = !aValue.IsLessThan(i, bValue)
	default:
		i.Fail("Unknown operator: " + opField.Value)
		return nilValue
	}
	return BoolValue(result)
}

func LogicTernaryEvaluator(i *Interpreter, b *Block) Value {
	ifBlock := b.SingleBlockValueWithName(i, "IF")
	thenBlock := b.SingleBlockValueWithName(i, "THEN")
	elseBlock := b.SingleBlockValueWithName(i, "ELSE")

	test := i.Evaluate(ifBlock).AsBoolean(i)
	if test {
		return i.Evaluate(thenBlock)
	} else {
		return i.Evaluate(elseBlock)
	}
}

func LogicBooleanEvaluator(i *Interpreter, b *Block) Value {
	f := b.FieldWithName("BOOL")
	if f == nil {
		i.Fail("Boolean block has no BOOL field")
		return nilValue
	}
	if f.Value == "TRUE" {
		return BoolValue(true)
	} else {
		return BoolValue(false)
	}
}

func LogicOperationEvaluator(i *Interpreter, b *Block) Value {
	aBlock := b.SingleBlockValueWithName(i, "A")
	bBlock := b.SingleBlockValueWithName(i, "B")
	opField := b.FieldWithName("OP")
	if opField == nil {
		i.Fail("Missing operator in logic operation block.")
		return nilValue
	}
	aValue := i.Evaluate(aBlock).AsBoolean(i)
	bValue := i.Evaluate(bBlock).AsBoolean(i)

	var result bool

	switch opField.Value {
	case "AND":
		result = aValue && bValue
	case "OR":
		result = aValue || bValue
	default:
		i.Fail("Unknown operator: " + opField.Value)
		return nilValue
	}
	return BoolValue(result)
}

func LogicNegateEvaluator(i *Interpreter, b *Block) Value {
	argBlock := b.SingleBlockValueWithName(i, "BOOL")
	return BoolValue(!i.Evaluate(argBlock).AsBoolean(i))
}
