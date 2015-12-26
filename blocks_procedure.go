package goblockly

import (
	"fmt"
)

// Note: procedures_defnoreturn and procedures_defreturn blocks are handled
// specially by Interpreter, since they are only valid as top-level blocks.

// ProceduresFunctionCallEvaluator calls a function, passing arguments and
// returning any return value from the function. Arguments shadow the global
// variable values in the context of the evaluation of the function.evaluates a
// function.
func ProceduresFunctionCallEvaluator(i *Interpreter, b *Block) Value {
	if b.Mutation == nil {
		i.Fail("Missing mutation in procedure call.")
		return nilValue
	}
	function, ok := i.Functions[b.Mutation.Name]
	if !ok {
		i.Fail("Unknown function name '" + b.Mutation.Name + "'")
		return nilValue
	}
	fmt.Printf("Calling function %s\n", b.Mutation.Name)
	argMap := make(map[string]Value)
	for idx, arg := range b.Mutation.Args {
		bv := b.BlockValueWithName(fmt.Sprintf("ARG%d", idx))
		if bv == nil {
			i.Fail(fmt.Sprintf("calling '%s': No value specified for argument '%s'",
				b.Mutation.Name,
				arg))
			return nilValue
		}
		if len(bv.Blocks) != 1 {
			i.Fail(fmt.Sprintf("calling '%s': Missing block for argument '%s'",
				b.Mutation.Name,
				arg))
			return nilValue
		}
		argMap[arg.Name] = i.Evaluate(&bv.Blocks[0])
	}

	shadowedVariables := ShadowVariables(i, argMap)
	if function.Body != nil {
		i.Evaluate(function.Body)
	}
	var retval Value = nilValue
	if function.Return != nil {
		retval = i.Evaluate(function.Return)
	} else {
		fmt.Printf("Function return is nil\n")
	}
	UnshadowVariables(i, argMap, shadowedVariables)
	return retval
}

// ShadowVariables builds a backup copy of variables in the Interpreter context,
// then sets those variables to the shadow values. It returns a map of the shadowed variables.
func ShadowVariables(i *Interpreter, newVariables map[string]Value) map[string]Value {
	shadowedValues := make(map[string]Value)
	for k, v := range newVariables {
		oldValue, ok := i.Context[k]
		if ok {
			shadowedValues[k] = oldValue
		}
		i.Context[k] = v
	}
	return shadowedValues
}

// UnshadowVariables restores the variables in the Interpreter to their
// pre-shadow values.
func UnshadowVariables(i *Interpreter, newVariables map[string]Value, preShadow map[string]Value) {
	for k, _ := range newVariables {
		delete(i.Context, k)
		if oldValue, ok := preShadow[k]; ok {
			i.Context[k] = oldValue
		}
	}
}
