package Pressure

import (
	"fmt"
	"math"
)

type Unit int8

const (
	Undefined Unit = iota
	Pascal
	Hectopascal
	Kilopascal
	InchOfMercury
)

func (u Unit) String() string {
	switch u {
	case Pascal:
		return "P"
	case Hectopascal:
		return "hP"
	case Kilopascal:
		return "kP"
	case InchOfMercury:
		return "InHg"
	}
	return "unknown"
}

type Pressure struct {
	value float64
	unit  Unit
}

func New(value float64, unit Unit) Pressure {
	p := Pressure{
		value: value,
		unit:  unit,
	}

	return p
}

func (d Pressure) Set(value float64, unit Unit) {
	d.value = value
	d.unit = unit
}

func (d Pressure) Get(unit Unit) float64 {
	switch unit {
	case Pascal:
		return toPascal(d.value, d.unit)
	case Hectopascal:
		return toPascal(d.value, d.unit) / 100.0
	case Kilopascal:
		return toPascal(d.value, d.unit) / 1000.0
	case InchOfMercury:
		return toInchOfMercury(d.value, d.unit)
	}

	return math.NaN()
}

func (d Pressure) ToStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f", d.Get(unit))
}

func (d Pressure) ToFullStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f%s", d.Get(unit), d.unit.String())
}

func toPascal(value float64, from Unit) float64 {
	switch from {
	case Pascal:
		return value
	case Hectopascal:
		return (value * 100.0)
	case Kilopascal:
		return (value * 1000.0)
	case InchOfMercury:
		return (value * 3386.3886666667)
	}

	return math.NaN()
}

func toInchOfMercury(value float64, from Unit) float64 {
	switch from {
	case Pascal:
		return value / 3386.3886666667
	case Hectopascal:
		return value / 33.863886666667
	case Kilopascal:
		return value / 3.3863886666667
	case InchOfMercury:
		return value
	}

	return math.NaN()
}
