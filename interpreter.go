// Blockly Interpreter
//
// Interpreter for the Blockly programming language (https://developers.google.com/blockly).
//
// Blockly is a library for building visual programming editors. This
// interpreter is capable of digesting the output of Blockly.Xml.domToText (in
// the Blockly library) by interpreting the resulting XML as a program and
// running that program. It supports most of the basic block types that Blockly
// ships with.
//
// To use:
//
// Create a blockly.Interpreter.
//
// Populate the interpreter with a Console (io.Writer) as an output destination
// and a FailHandler to be run if the Blockly script cannot be interpreted.
//
// Parse the XML into a BlockXml using the tools in encoding/xml.
//
// Run the first block in the XML with:
//   interpreter.Run(blockXml.Blocks[0])
//
// If BlockXml contains more than one block, this implies that the Blockly
// workspace had multiple disconnected blocks in it. You can run the subsequent
// blocks, but generally these are ignored.
package goblockly

import (
	"fmt"
	"io"
	"strings"
)

type BreakType int

// Breaks can either be 'break' (end the current loop) or 'continue' (go to
// next step in the current loop). THese constants indicate which type of break
// occurred.
const (
	ThenBreak BreakType = iota
	ThenContinue
)

// BreakEvent is a special type sent via panic to indicate a break or continue
// occurred. Though Blockly's client-side editor will warn a user if they
// attempt to use a break block outside of a loop, Blockly will still allow the
// user to send a block in that form to the server, so we panic on breaks to
// allow the interpreter to handle this state.
type BreakEvent struct {
	Then BreakType
}

// Error method for BreakEvents (which are panicked as errors).
func (evt BreakEvent) Error() string {
	switch evt.Then {
	case ThenBreak:
		return "Break outside of loop."
	case ThenContinue:
		return "Continue outside of loop."
	default:
		return "BreakEvent outside of loop."
	}
}

// An Evaluator can evaluate a block, in the current context of the interpreter, into a value.
type Evaluator func(*Interpreter, *Block) Value

// The Interpreter maintains interpretation state for Blockly evaluation (such
// as where print operations should go, what to do if evaluation fails, and
// variable values).
type Interpreter struct {
	// The output console. All 'print' operations write to here.
	Console io.Writer
	// Function called if evaluation fails
	FailHandler func(string)
	// Variables in an execution cycle of the interpreter
	Context map[string]Value
	// Custom handlers for specific block prefixes. Inserting an evaluator
	// here will cause all blocks with type "prefix_<something>" to be
	// handled by the evaluator in key "prefix_". Blocks with a type not in
	// prefix get handled by the default evaluators or cause the interpreter
	// to Fail if there is no evaluator for the block type.
	PrefixHandlers map[string]Evaluator
}

// Table of default evaluators for block types.
var evaluators map[string]Evaluator

// PrepareEvaluators populates the default evaluator table. It must be called once before any calls to Interpreter.Run.
func PrepareEvaluators() {
	evaluators = map[string]Evaluator{
		"controls_if":              ControlIfEvaluator,
		"controls_repeat_ext":      ControlRepeatExtEvaluator,
		"controls_whileUntil":      ControlWhileUntil,
		"controls_for":             ControlForEvaluator,
		"controls_forEach":         ControlForEachEvaluator,
		"controls_flow_statements": ControlFlowStatements,

		"logic_compare": LogicCompareEvaluator,
		"logic_ternary": LogicTernaryEvaluator,

		"logic_boolean":   LogicBooleanEvaluator,
		"logic_operation": LogicOperationEvaluator,
		"logic_negate":    LogicNegateEvaluator,

		"math_number":     NumberEvaluator,
		"math_arithmetic": NumberArithmeticEvaluator,

		"text":       TextEvaluator,
		"text_join":  TextJoinEvaluator,
		"text_print": PrintEvaluator,

		"lists_create_empty": ListCreateEmptyEvaluator,
		"lists_create_with":  ListCreateWithEvaluator,
		"lists_repeat":       ListRepeatEvaluator,
		"lists_length":       ListLengthEvaluator,
		"lists_isEmpty":      ListIsEmptyEvaluator,
		"lists_indexOf":      ListIndexOfEvaluator,
		"lists_getIndex":     ListGetIndexEvaluator,
		"lists_setIndex":     ListSetIndexEvaluator,
		"lists_getSublist":   ListGetSublistEvaluator,
		"lists_split":        ListSplitEvaluator,

		"colour_picker": ColourPickerEvaluator,
		"colour_random": ColourRandomEvaluator,
		"colour_rgb":    ColourRgbEvaluator,
		"colour_blend":  ColourBlendEvaluator,

		"variables_set": VariableSetEvaluator,
		"variables_get": VariableGetEvaluator,
	}
}

// Fail causes interpretation to panic. If Fail is called in the context of a
// Run, the Run will recover and call the interpreter's FailHandler function.
func (i *Interpreter) Fail(reason string) {
	panic(reason)
}

// WriteToConsole outputs a string to the interpreter's Console.
func (i *Interpreter) WriteToConsole(s string) {
	i.Console.Write([]byte(s + "\n"))
}

// Run evaluates a top-level block and handles Fail panics by recovering them
// into the interpreter's FailHandler. It also initializes the Interpreter's Context (i.e. clears all variables).
func (i *Interpreter) Run(b *Block) {
	defer func() {
		if r := recover(); r != nil {
			i.FailHandler(fmt.Sprint(r))
			return
		}
	}()

	i.Context = make(map[string]Value)
	i.Evaluate(b)
}

// Evaluate evaluates a specific block by determining what evaluator can consume
// it. Generally, this is called by Run and by other evaluators; it should not
// need to be called directly.
//
// PrefixHandlers may call Evaluate directly if they are evaluating a type of
// block that itself has values or statements.
func (i *Interpreter) Evaluate(b *Block) Value {
	var evaluator Evaluator
	for k, v := range i.PrefixHandlers {
		if strings.HasPrefix(b.Type, k) {
			evaluator = v
			break
		}
	}
	if evaluator == nil {
		evaluator = evaluators[b.Type]
	}
	if evaluator == nil {
		i.Fail("No evaluator for block '" + b.Type + "'")
		return nilValue
	}
	value := evaluator(i, b)
	if b.Next == nil {
		return value
	} else {
		value = i.Evaluate(b.Next)
		return value
	}
}

// CheckBreak is a helper function that checks if a panic was due to breaking
// from a loop. If it was not, it re-panicks. If it was, it sets continuing to
// true if the enclosing loop should contiue with the next iteration, or false
// if the enclosing loop should break out.
func (i *Interpreter) CheckBreak(continuing *bool) {
	r := recover()

	if r != nil {
		if breakEvent, ok := r.(BreakEvent); ok {
			*continuing = breakEvent.Then == ThenContinue
		} else {
			panic(r)
		}
	}
}
