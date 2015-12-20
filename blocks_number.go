// Evaluators for number blocks

package goblockly

import (
	"math"
	"strconv"
)

func NumberArithmeticEvaluator(i *Interpreter, b *Block) Value {
	aBlock := b.SingleBlockValueWithName(i, "A")
	bBlock := b.SingleBlockValueWithName(i, "B")
	opField := b.FieldWithName("OP")
	if opField == nil {
		i.Fail("Missing operator in arithmetic block.")
		return nilValue
	}
	aValue := i.Evaluate(aBlock).AsNumber(i)
	bValue := i.Evaluate(bBlock).AsNumber(i)

	var result float64

	switch opField.Value {
	case "ADD":
		result = aValue + bValue
	case "MINUS":
		result = aValue - bValue
	case "MULTIPLY":
		result = aValue * bValue
	case "DIVIDE":
		result = aValue / bValue
	case "POWER":
		result = math.Pow(aValue, bValue)
	default:
		i.Fail("Unknown operator: " + opField.Value)
		return nilValue
	}
	return NumberValue(result)
}

func NumberEvaluator(i *Interpreter, b *Block) Value {
	f := b.FieldWithName("NUM")
	if f == nil {
		i.Fail("Number block has no NUM field")
		return nilValue
	}

	val, err := strconv.ParseFloat(f.Value, 64)
	if err != nil {
		i.Fail(err.Error())
		return nilValue
	}
	return NumberValue(val)
}
