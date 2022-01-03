package Moisture

import (
	"fmt"
)

type Moisture struct {
	value int64
}

func New(value int64) Moisture {
	m := Moisture{
		value: value,
	}
	return m
}

func (d Moisture) Set(value int64) {
	d.value = value
}

func (d Moisture) Get() int64 {
	return d.value
}

func (d Moisture) ToString() string {
	return fmt.Sprintf("%d", d.value)
}

func (d Moisture) ToFullString() string {
	return fmt.Sprintf("%d%s", d.value, "%")
}
