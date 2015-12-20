// Evaluators for text blocks

package blockly

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
