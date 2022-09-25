package convert

import (
	"fmt"

	"github.com/martinlindhe/unit"
)

var unitMap = map[string]unit.Length{
	"mi":         unit.Mile,
	"mile":       unit.Mile,
	"miles":      unit.Mile,
	"km":         unit.Kilometer,
	"kilometer":  unit.Kilometer,
	"kilometers": unit.Kilometer,
	"kilometre":  unit.Kilometer,
	"kilometres": unit.Kilometer,
	"m":          unit.Meter,
	"meter":      unit.Meter,
	"meters":     unit.Meter,
	"metre":      unit.Meter,
	"metres":     unit.Meter,
}

func ToMeters(value float64, originalUnit string) (float64, error) {
	lengthUnit := unitMap[originalUnit]
	if lengthUnit == 0 {
		return 0, fmt.Errorf("unknown unit provided: %#v", originalUnit)
	}
	return (unit.Length(value) * lengthUnit).Meters(), nil
}
