package dimension

import (
	"fmt"
	"gless/unit"
	"regexp"
	"strconv"
)

var regexpNumber *regexp.Regexp = regexp.MustCompile(`^(\d+\.{0,1}\d*)`)

type Dimension struct {
	Value float64
	Unit  unit.Unit
}

func NewDimension(value string) (*Dimension, error) {
	number, err := strconv.ParseFloat(regexpNumber.FindString(value), 64)
	if err != nil {
		return nil, err
	}

	unit, _ := unit.GetUnit(value)

	return &Dimension{Value: number, Unit: unit}, nil
}

func (dimension *Dimension) String() string {
	return fmt.Sprintf("%g%s", dimension.Value, dimension.Unit)
}
