package blockly

import (
	"fmt"
)

// ControlIfEvaluator evaluates if / elseif / else blocks.
func ControlIfEvaluator(i *Interpreter, b *Block) Value {
	var elseIfs int
	var elses int
	if b.Mutations != nil {
		elseIfs = b.Mutations[0].ElseIf
		elses = b.Mutations[0].Else
	}

	for idx := 0; idx <= elseIfs; idx++ {
		ifBlock := b.SingleBlockValueWithName(i, fmt.Sprintf("IF%d", idx))
		doBlock := b.SingleBlockStatementWithName(i, fmt.Sprintf("DO%d", idx))
		if i.Evaluate(ifBlock).AsBoolean(i) {
			return i.Evaluate(doBlock)
		}
	}

	if elses > 0 {
		return i.Evaluate(b.SingleBlockStatementWithName(i, "ELSE"))
	}

	return nilValue
}

func ControlRepeatExtEvaluator(i *Interpreter, b *Block) Value {
	repeatedBlock := b.SingleBlockStatementWithName(i, "DO")
	var result Value
	continuing := true
	for repeats := int(i.Evaluate(
		b.SingleBlockValueWithName(i, "TIMES")).AsNumber(i)); repeats > 0 && continuing; repeats-- {
		func() {
			defer i.CheckBreak(&continuing)
			result = i.Evaluate(repeatedBlock)
		}()
	}
	return result
}

func ControlWhileUntil(i *Interpreter, b *Block) Value {
	modeField := b.FieldWithName("MODE")
	if modeField == nil {
		i.Fail("No MODE field in controls_whileUntil")
		return nilValue
	}
	predicate := b.SingleBlockValueWithName(i, "BOOL")
	body := b.SingleBlockStatementWithName(i, "DO")

	running := true
	var v Value
	v = nilValue
	for running {
		switch modeField.Value {
		case "WHILE":
			running = i.Evaluate(predicate).AsBoolean(i)
		case "UNTIL":
			running = !(i.Evaluate(predicate).AsBoolean(i))
		}

		if running {
			func() {
				defer i.CheckBreak(&running)
				v = i.Evaluate(body)
			}()
		}
	}
	return v
}

func ControlForEvaluator(i *Interpreter, b *Block) Value {
	varName := b.FieldWithName("VAR")
	if varName == nil {
		i.Fail("No VAR field in controls_for")
		return nilValue
	}
	fromValue := i.Evaluate(b.SingleBlockValueWithName(i, "FROM")).AsNumber(i)
	toValue := i.Evaluate(b.SingleBlockValueWithName(i, "TO")).AsNumber(i)
	byValue := i.Evaluate(b.SingleBlockValueWithName(i, "BY")).AsNumber(i)
	body := b.SingleBlockStatementWithName(i, "DO")
	running := true
	var v Value
	v = nilValue

	for ; (fromValue <= toValue) && running; fromValue += byValue {
		i.Context[varName.Value] = NumberValue(fromValue)
		func() {
			defer i.CheckBreak(&running)
			v = i.Evaluate(body)
		}()
	}

	return v
}

// ControlForEachEvaluator runs the body of the control with the index variable
// set to each element of a list.
func ControlForEachEvaluator(i *Interpreter, b *Block) Value {
	varName := b.SingleFieldWithName(i, "VAR")
	list := i.Evaluate(b.SingleBlockValueWithName(i, "LIST")).AsList(i)
	body := b.SingleBlockStatementWithName(i, "DO")

	running := true

	for idx := 0; idx < len(*list.Values) && running; idx += 1 {
		i.Context[varName] = (*list.Values)[idx]
		func() {
			defer i.CheckBreak(&running)
			i.Evaluate(body)
		}()
	}

	return nilValue
}

func ControlFlowStatements(i *Interpreter, b *Block) Value {
	breakType := b.FieldWithName("FLOW")
	if breakType == nil {
		i.Fail("No FLOW field in controls_flow_statements")
		return nilValue
	}
	switch breakType.Value {
	case "BREAK":
		panic(BreakEvent{Then: ThenBreak})
	case "CONTINUE":
		panic(BreakEvent{Then: ThenContinue})
	}
	return nilValue
}
