package main

import "fmt"

type Status int

const (
	Susceptive Status = iota
	Infective
	Removed
)

func (s Status) String() string {
	switch s {
	case Susceptive:
		return "Susceptive"
	case Infective:
		return "infective"
	case Removed:
		return "Removed"
	default:
		return fmt.Sprintf("Status(%d)", s)
	}
}
