package factors

import (
	"fmt"
	"regexp/syntax"
)

// RegexpOpCodeTable is a set of regexp operations.
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

// Op is systax.Op.
type Op syntax.Op

// String is a string representation of the regexp operation.
func (o Op) String() string {
	op := int(o)
	if op >= 0 && op < len(RegexpOpCodeTable) {
		return RegexpOpCodeTable[op]
	}
	return fmt.Sprintf("UNDEF(%d)", op)
}
