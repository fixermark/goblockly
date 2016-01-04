GoBlockly Interpreter
=====================

GoBlockly Interpreter is an interpreter for the output of the Blockly visual programming editor.

What is Blockly?
----------------

Blockly is a web-based library for building visual programming editors, in which
users can drag blocks together to build programs. Its focus is beginner
programming education, but it's useful as a simple visual editing tool; the
block-based rules for connecting statements together make syntax errors entirely impossible.

More information is available at https://developers.google.com/blockly.

What is the GoBlockly Interpreter?
----------------------------------

The GoBlockly Interpreter is a go library for interpreting the output of the
`Blockly.Xml.domToText` command from the Blockly library by interpreting the
resulting XML as a program and running that program. It's useful for evaluating
programs server-side and supports extending the interpreter with handling for
your own blocks.

How do I use it?
----------------

Details are provided in the godocs of the library, but the basic overview is:

    import (
    	"bytes"
    	"encoding/xml"
    )

    // XML may have come from client request or database store; it's the output of
    Blockly.Xml.domToText in the Blockly library
    func interpretBlockly(string xml) string {
    	var blocks BlockXml
    	if err := xml.Unmarshal(xmldata, &blocks); err {
    		// Handle malformed XML here
    	}
    	var b bytes.Buffer
    	var i Interpreter
    	i.Console = &b
    	i.FailHandler = func(reason string) {
    		// Handle parser failure here
    	}
    	i.Run(blocks.Blocks)
    	consoleOutput := b.String()
    	return consoleOutput
    }

The code runs server-side; how secure is it?
--------------------------------------------

The GoBlockly interpreter offers no functionality beyond the basic Blockly
blocks out of the box, and there are no exploits known that could allow the
injection of arbitrary behavior via the interpreter. Of course, running
user-provided code server-side always carries some risk; it is always
recommended to run a server process with no more permissions than it needs to
accomplish its goals on the host platform.

Why?
----

Generally, Blockly is converted client-side into code in another language and
executed. Why write an interpreter?

As I rewrite [Belphanior Butler](https://github.com/fixermark/belphanior-butler)
in Go, I wish to recapitulate its ability to hand off the running of all of the
servant manipulation server-side, so it seemed straightforward to recapitulate
the Ruby interpreter for Blockly as a go interpreter (as opposed to, say,
pulling in Node.js and running Blockly converted into JS server-side). This
allows the server to both store the Blockly programs and execute them directly
without having to loop through a web client.

Project Details
---------------

Project copyright 2015 [Mark T. Tomczak](mailto:fixermark@gmail.com).

Project is licensed under the Apache License, version 2.0 (see COPYING for more
information).


