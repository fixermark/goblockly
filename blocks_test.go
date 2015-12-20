// Test for Blockly block XML

package blockly

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestDeserializeBlock(t *testing.T) {
	blocks := `
<xml xmlns="http://www.w3.org/1999/xhtml">
  <block type="text_print" x="29" y="34">
    <value name="TEXT">
      <block type="text">
	<field name="TEXT">hello, world</field>
      </block>
    </value>
  </block>
</xml>
`
	var unmarshaled BlockXml
	err := xml.Unmarshal([]byte(blocks), &unmarshaled)

	if err != nil {
		t.Fatal("Could not deserialize block: ", err)
	}
	xmlName := xml.Name{"http://www.w3.org/1999/xhtml", "xml"}
	blockName := xml.Name{"http://www.w3.org/1999/xhtml", "block"}
	compare := BlockXml{
		XMLName: xmlName,
		Blocks: []Block{
			Block{
				XMLName: blockName,
				Type:    "text_print",
				X:       "29",
				Y:       "34",
				Values: []BlockValue{
					BlockValue{
						Name: "TEXT",
						Blocks: []Block{
							Block{
								XMLName: blockName,
								Type:    "text",
								Fields: []BlockField{
									BlockField{
										Name:  "TEXT",
										Value: "hello, world",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(unmarshaled, compare) {
		t.Error("Cannot match unmarshal. Expected ",
			compare,
			"\nObserved:", unmarshaled)
	}
}
