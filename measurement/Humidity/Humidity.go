package Humidity

import (
	"fmt"
)

type Humidity struct {
	value int64
}

func New(value int64) Humidity {
	h := Humidity{
		value: value,
	}

	return h
}

func (d Humidity) Set(value int64) {
	d.value = value
}

func (d Humidity) Get() int64 {
	return d.value
}

func (d Humidity) ToString() string {
	return fmt.Sprintf("%d", d.value)
}

func (d Humidity) ToFullString() string {
	return fmt.Sprintf("%d%s", d.value, "%")
}
