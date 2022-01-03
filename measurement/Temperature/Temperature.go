package Temperature

import (
	"fmt"
	"math"
)

type Unit int8

const (
	Undefined Unit = iota
	Kelvin
	Celsius
	Farenheit
)

func (u Unit) String() string {
	switch u {
	case Kelvin:
		return "K"
	case Celsius:
		return "°C"
	case Farenheit:
		return "°F"
	}
	return "unknown"
}

type Temperature struct {
	value float64
	unit  Unit
}

func New(value float64, unit Unit) Temperature {
	t := Temperature{value, unit}
	return t
}

func (d Temperature) Set(value float64, unit Unit) {
	d.value = value
	d.unit = unit
}

func (d Temperature) Get(unit Unit) float64 {
	switch unit {
	case Celsius:
		return toCelsius(d.value, d.unit)
	case Kelvin:
		return toKelvin(d.value, d.unit)
	case Farenheit:
		return toFarenheit(d.value, d.unit)
	}

	return math.NaN()
}

func (d Temperature) ToStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f", d.Get(unit))
}

func (d Temperature) ToFullStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f%s", d.Get(unit), unit.String())
}

func toCelsius(value float64, from Unit) float64 {
	switch from {
	case Celsius:
		return value
	case Kelvin:
		return (value - 273.15)
	case Farenheit:
		return ((value - 32.0) * (5.0 / 9.0))
	}

	return math.NaN()
}

func toFarenheit(value float64, from Unit) float64 {
	switch from {
	case Celsius:
		return ((value * (9.0 / 5.0)) + 32.0)
	case Kelvin:
		return ((toCelsius(value, from) * (9.0 / 5.0)) + 32.0)
	case Farenheit:
		return value
	}

	return math.NaN()
}

func toKelvin(value float64, from Unit) float64 {
	switch from {
	case Celsius:
		return (value + 273.15)
	case Kelvin:
		return value
	case Farenheit:
		return (toCelsius(value, from) + 273.15)
	}

	return math.NaN()
}
