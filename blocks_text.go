// Evaluators for text blocks

package goblockly

import (
	"fmt"
	"math/rand"
	"strings"
)

func TextEvaluator(i *Interpreter, b *Block) Value {
	f := b.FieldWithName("TEXT")
	if f == nil {
		i.Fail("Text block has no TEXT field")
		return nilValue
	}
	return StringValue(f.Value)
}

func TextJoinEvaluator(i *Interpreter, b *Block) Value {
	var result string
	for _, v := range b.Values {
		if len(v.Blocks) != 1 {
			i.Fail("A join block socket does not have exactly one block attached to it.")
			return nilValue
		}
		result += i.Evaluate(&v.Blocks[0]).AsString(i)
	}
	return StringValue(result)
}

func PrintEvaluator(i *Interpreter, b *Block) Value {
	v := b.BlockValueWithName("TEXT")
	if v == nil {
		i.Fail("Print block has no TEXT value")
		return nilValue
	}
	if len(v.Blocks) != 1 {
		i.Fail("Print block should have exactly one block attached to it.")
		return nilValue
	}
	i.WriteToConsole(i.Evaluate(&v.Blocks[0]).AsString(i))
	return nilValue
}

// TextLengthEvaluator returns the length of a text block
func TextLengthEvaluator(i *Interpreter, b *Block) Value {
	return NumberValue(len(
		i.Evaluate(b.SingleBlockValueWithName(i, "VALUE")).AsString(i)))
}

// TextIsEmptyEvaluator returns true if the length of the text is 0
func TextIsEmptyEvaluator(i *Interpreter, b *Block) Value {
	return BoolValue(len(
		i.Evaluate(b.SingleBlockValueWithName(i, "VALUE")).AsString(i)) == 0)
}

// TextCharAtEvaluator returns the character at a specific location in a text block.
func TextCharAtEvaluator(i *Interpreter, b *Block) Value {
	textToScan := i.Evaluate(b.SingleBlockValueWithName(i, "VALUE")).AsString(i)
	where := b.SingleFieldWithName(i, "WHERE")
	switch where {
	case "FROM_START":
		idx := int(i.Evaluate(b.SingleBlockValueWithName(i, "AT")).AsNumber(i))
		if idx < 1 {
			i.Fail(fmt.Sprintf("Cannot retrieve character at index %d", idx))
			return nilValue
		}
		if idx > len(textToScan) {
			i.Fail(fmt.Sprintf("Cannot retrieve character at index %d; string length is %d",
				idx,
				len(textToScan)))
			return nilValue
		}
		return StringValue(textToScan[idx-1])
	case "FROM_END":
		idx := int(i.Evaluate(b.SingleBlockValueWithName(i, "AT")).AsNumber(i))
		if idx < 1 {
			i.Fail(fmt.Sprintf("Cannot retrieve character at index %d", idx))
			return nilValue
		}
		if idx > len(textToScan) {
			i.Fail(fmt.Sprintf("Cannot retrieve character at index %d; string length is %d",
				idx,
				len(textToScan)))
			return nilValue
		}
		return StringValue(textToScan[len(textToScan)-idx])
	case "FIRST":
		if len(textToScan) < 1 {
			i.Fail("Cannot return first character from empty string.")
			return nilValue
		}
		return StringValue(textToScan[0])
	case "LAST":
		if len(textToScan) < 1 {
			i.Fail("Cannot return last character from empty string.")
			return nilValue
		}
		return StringValue(textToScan[len(textToScan)-1])
	case "RANDOM":
		if len(textToScan) < 1 {
			i.Fail("Cannot return random character from empty string.")
			return nilValue
		}
		return StringValue(textToScan[rand.Intn(len(textToScan))])
	default:
		i.Fail(fmt.Sprintf("Unknown charAt rule: %s", where))
		return nilValue
	}
}

// TextIndexOfEvaluator returns the first (one-indexed) location at which a
// specified substring occurs in the searched string. Returns 0 if text is not
// found.
func TextIndexOfEvaluator(i *Interpreter, b *Block) Value {

	stringToSearch := i.Evaluate(b.SingleBlockValueWithName(i, "VALUE")).AsString(i)
	substring := i.Evaluate(b.SingleBlockValueWithName(i, "FIND")).AsString(i)

	findFrom := b.SingleFieldWithName(i, "END")
	switch findFrom {
	case "FIRST":
		return NumberValue(strings.Index(stringToSearch, substring) + 1)
	case "LAST":
		return NumberValue(strings.LastIndex(stringToSearch, substring) + 1)
	default:
		i.Fail("text_charAt does not know how to " + findFrom)
		return nilValue
	}
}

// TextGetSubstringEvaluator returns the specified substring.
func TextGetSubstringEvaluator(i *Interpreter, b *Block) Value {
	// TODO(mtomczak): Implement.
	return nilValue
}

// TextAppendEvaluator appends the specified text to the value in the specified
// variable.
func TextAppendEvaluator(i *Interpreter, b *Block) Value {
	varName := b.SingleFieldWithName(i, "VAR")
	toAppend := i.Evaluate(b.SingleBlockValueWithName(i, "TEXT")).AsString(i)
	val, ok := i.Context[varName]
	if !ok {
		i.Context[varName] = StringValue(toAppend)
	} else {
		i.Context[varName] = StringValue(val.AsString(i) + toAppend)
	}
	return nilValue
}
