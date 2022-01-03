package Velocity

import (
	"fmt"
	"math"
)

type Unit int8

const (
	Undefined Unit = iota
	MetresPerSecond
	KilometresPerHour
	MilesPerHour
)

func (u Unit) String() string {
	switch u {
	case MetresPerSecond:
		return "m/s"
	case KilometresPerHour:
		return "km/h"
	case MilesPerHour:
		return "mph"
	}
	return "unknown"
}

type Velocity struct {
	value float64
	unit  Unit
}

func New(value float64, unit Unit) Velocity {
	v := Velocity{
		value: value,
		unit:  unit,
	}
	return v
}

func (d Velocity) Set(value float64, unit Unit) {
	d.value = value
	d.unit = unit
}

func (d Velocity) Get(unit Unit) float64 {
	switch unit {
	case MetresPerSecond:
		return toMetresPerSecond(d.value, d.unit)
	case KilometresPerHour:
		return toMetresPerSecond(d.value, d.unit) * 3.6
	case MilesPerHour:
		return toMilesPerHour(d.value, d.unit)
	}

	return math.NaN()
}

func (d Velocity) ToStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f", d.Get(unit))
}

func (d Velocity) ToFullStringAs(unit Unit) string {
	return fmt.Sprintf("%.2f%s", d.Get(unit), d.unit.String())
}

func toMetresPerSecond(value float64, from Unit) float64 {
	switch from {
	case MetresPerSecond:
		return value
	case KilometresPerHour:
		return (value / 3.6)
	case MilesPerHour:
		return (value / 2.237)
	}

	return math.NaN()
}

func toMilesPerHour(value float64, from Unit) float64 {
	switch from {
	case MetresPerSecond:
		return (value * 2.237)
	case KilometresPerHour:
		return (value / 1.609)
	case MilesPerHour:
		return value
	}

	return math.NaN()
}
