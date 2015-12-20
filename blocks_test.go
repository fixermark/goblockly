// Test for Blockly block XML

package goblockly

import (
	"encoding/xml"
	"runtime"
	"testing"
)

func stringsEqual(t *testing.T, a, b string) {
	if a != b {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "UNKNOWN"
			line = 0
		}
		t.Errorf("%s:%d - Expected %s to equal %s", file, line, a, b)
	}
}

func intsEqual(t *testing.T, a, b int) {
	if a != b {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "UNKNOWN"
			line = 0
		}
		t.Errorf("%s:%d - Expected %d to equal %d", file, line, a, b)
	}
}

func TestDeserializeBlock(t *testing.T) {
	blocks := `
<xml xmlns="http://www.w3.org/1999/xhtml">
  <block type="controls_for" x="76" y="-1861">
    <field name="VAR">i</field>
    <value name="FROM">
      <block type="math_number">
	<field name="NUM">0</field>
      </block>
    </value>
    <value name="TO">
      <block type="math_number">
	<field name="NUM">10</field>
      </block>
    </value>
    <value name="BY">
      <block type="math_number">
	<field name="NUM">2</field>
      </block>
    </value>
    <statement name="DO">
      <block type="text_print">
	<value name="TEXT">
	  <block type="variables_get">
	    <field name="VAR">i</field>
	  </block>
	</value>
	<next>
	  <block type="controls_flow_statements">
	    <field name="FLOW">CONTINUE</field>
	  </block>
	</next>
      </block>
    </statement>
    <next>
      <block type="controls_if">
	<mutation else="1">
	</mutation>
	<value name="IF0">
	  <block type="logic_compare">
	    <field name="OP">GT</field>
	    <value name="A">
	      <block type="variables_get">
		<field name="VAR">i</field>
	      </block>
	    </value>
	    <value name="B">
	      <block type="math_number">
		<field name="NUM">5</field>
	      </block>
	    </value>
	  </block>
	</value>
	<statement name="DO0">
	  <block type="text_print">
	    <value name="TEXT">
	      <block type="text">
		<field name="TEXT">hi</field>
	      </block>
	    </value>
	  </block>
	</statement>
	<statement name="ELSE">
	  <block type="text_print">
	    <value name="TEXT">
	      <block type="text">
		<field name="TEXT">hello</field>
	      </block>
	    </value>
	  </block>
	</statement>
      </block>
    </next>
  </block>
</xml>
`
	var unmarshaled BlockXml
	err := xml.Unmarshal([]byte(blocks), &unmarshaled)

	if err != nil {
		t.Fatal("Could not deserialize block: ", err)
	}

	i := &Interpreter{
		nil,
		func(s string) {
			t.Errorf("Interpreter failed: %s", s)
		},
		nil,
		nil,
	}

	stringsEqual(t, unmarshaled.XMLName.Space, "http://www.w3.org/1999/xhtml")
	stringsEqual(t, unmarshaled.XMLName.Local, "xml")
	intsEqual(t, len(unmarshaled.Blocks), 1)
	loopBlock := unmarshaled.Blocks[0]
	stringsEqual(t, loopBlock.Type, "controls_for")
	stringsEqual(t, loopBlock.X, "76")
	stringsEqual(t, loopBlock.Y, "-1861")

	varField := unmarshaled.Blocks[0].SingleFieldWithName(i, "VAR")
	stringsEqual(t, varField, "i")

	fromBlock := unmarshaled.Blocks[0].SingleBlockValueWithName(i, "FROM")
	stringsEqual(t, fromBlock.Type, "math_number")
	stringsEqual(t, fromBlock.SingleFieldWithName(i, "NUM"), "0")

	body := loopBlock.SingleBlockStatementWithName(i, "DO")
	stringsEqual(t, body.Type, "text_print")
	body2 := body.Next
	stringsEqual(t, body2.Type, "controls_flow_statements")

	ifBlock := loopBlock.Next
	intsEqual(t, ifBlock.Mutation.Else, 1)
	ifBlock.SingleBlockValueWithName(i, "IF0")

}
