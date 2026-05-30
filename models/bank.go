package models

import "slices"

type BankName string

const (
	WellsFargo BankName = "WF"
	Fidelity   BankName = "Fidelity"
)

type Bank struct {
	ID   int      `json:"id"`
	Name BankName `json:"name"`
}

// KnownBanks lists all banks that should exist on startup
var KnownBanks = []BankName{
	WellsFargo,
	Fidelity,
}

func (b BankName) IsKnown() bool {
	return slices.Contains(KnownBanks, b)
}
