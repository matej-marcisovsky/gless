package unit

import "regexp"

var regexpUnit *regexp.Regexp = regexp.MustCompile(`(px|em|ex|ch|rem|in|cm|mm|pc|pt|vw|vh|vmin|vmax)$`)

type Unit string

const (
	PX     Unit = "px"
	EM     Unit = "em"
	EX     Unit = "ex"
	CH     Unit = "ch"
	REM    Unit = "rem"
	IN     Unit = "in"
	CM     Unit = "cm"
	MM     Unit = "mm"
	PC     Unit = "pc"
	PT     Unit = "pt"
	VW     Unit = "vw"
	VH     Unit = "vh"
	VMIN   Unit = "vmin"
	VMAX   Unit = "vmax"
	NUMBER Unit = "number"
)

func GetUnit(value string) (Unit, bool) {
	_unit := regexpUnit.FindString(value)

	return Unit(_unit), _unit != ""
}
