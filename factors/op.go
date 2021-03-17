package factors

import (
	"fmt"
	"regexp/syntax"
)

var RegexpOpCodeTable = []string{
	"",
	"NoMatch",
	"EmptyMatch",
	"Literal",
	"CC",
	"AnyCharNL",
	"AnyChar",
	"^",
	"$",
	"\\A",
	"$",
	"\\b",
	"\\B",
	"()",
	"*",
	"+",
	"?",
	"RT",
	"・",
	"|",
}

type Op syntax.Op

func (o Op) String() string {
	op := int(o)
	if op >= 0 && op < len(RegexpOpCodeTable) {
		return RegexpOpCodeTable[op]
	}
	return fmt.Sprintf("UNDEF(%d)", op)
}
