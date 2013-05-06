package script

import (
	"bytes"
	"testing"
)

func TestLexer(t *testing.T) {
	src := bytes.NewBuffer([]byte(testText))
	parse(src)
}

const testText = `alpha=1
	save(avg(m), "m.dump", 1e-12)
	run(1e-9)
	b=sin(2,(3,4)) // bye bye;
	c=(1,2,3)
`
