package Rainfall

import (
	"fmt"
	"math"
)

type Unit int8

const (
	Undefined Unit = iota
	Millimetre
	Centimetre
	Inch
)

func (u Unit) String() string {
	switch u {
	case Millimetre:
		return "mm"
	case Centimetre:
		return "cm"
	case Inch:
		return "in"
	}
	return "unknown"
}

type Rainfall struct {
	value float64
	unit  Unit
}

func New(value float64, unit Unit) Rainfall {
	r := Rainfall{
		value: value,
		unit:  unit,
	}
	return r
}

func (d Rainfall) Set(value float64, unit Unit) {
	d.value = value
	d.unit = unit
}

func (d Rainfall) Get(unit Unit) float64 {
	switch unit {
	case Millimetre:
		return toMillimetre(d.value, d.unit)
	case Centimetre:
		return toMillimetre(d.value, d.unit) * 10.0
	case Inch:
		return toInch(d.value, d.unit)
	}

	return math.NaN()
}

func (d Rainfall) ToStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f", d.Get(unit))
}

func (d Rainfall) ToFullStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f%s", d.Get(unit), d.unit.String())
}

func toMillimetre(value float64, from Unit) float64 {
	switch from {
	case Millimetre:
		return value
	case Centimetre:
		return (value * 10.0)
	case Inch:
		return (value * 25.4)
	}

	return math.NaN()
}

func toInch(value float64, from Unit) float64 {
	switch from {
	case Millimetre:
		return (value / 25.4)
	case Centimetre:
		return (value / 2.54)
	case Inch:
		return value
	}

	return math.NaN()
}
