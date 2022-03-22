package unit

import "regexp"

var regexpLength *regexp.Regexp = regexp.MustCompile(`(px|em|ex|ch|rem|in|cm|mm|pc|pt|vw|vh|vmin|vmax)$`)

type Unit string

const (
	CM       Unit = "cm"
	EM       Unit = "em"
	EX       Unit = "ex"
	CH       Unit = "ch"
	IN       Unit = "in"
	MM       Unit = "mm"
	PC       Unit = "pc"
	PT       Unit = "pt"
	PX       Unit = "px"
	REM      Unit = "rem"
	SINGULAR Unit = ""
	VH       Unit = "vh"
	VMAX     Unit = "vmax"
	VMIN     Unit = "vmin"
	VW       Unit = "vw"
)

// TODO add other units
func GetUnit(value string) (Unit, bool) {
	unit := regexpLength.FindString(value)

	return Unit(unit), unit != ""
}
