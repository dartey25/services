package model

import "time"

type Identificator struct {
	Type  *string `json:"type,omitempty"`
	Value *string `json:"value,omitempty"`
	Note  *string `json:"note,omitempty"`
}

type SanctionsRow struct {
	SID     int `json:"-"`
	Company struct {
		Name           string         `json:"name"`
		TranslitName   string         `json:"translitName"`
		Aliases        []string       `json:"aliases,omitempty"`
		Status         string         `json:"status"`
		Country        string         `json:"country"`
		RegID          *Identificator `json:"regId,omitempty"`
		TaxID          *Identificator `json:"taxId,omitempty"`
		AdditionalInfo []string       `json:"additionalInfo"`
	} `json:"company"`
	Sanctions struct {
		Active  map[string]string `json:"active"`
		Term    string            `json:"term"`
		EndDate *time.Time        `json:"endDate"`
	} `json:"sanctions"`
	LastAction struct {
		Type     string `json:"type"`
		Document struct {
			Num              string     `json:"num"`
			Date             *time.Time `json:"date,omitempty"`
			Appendix         string     `json:"appendix"`
			AppendixPosition string     `json:"appendixPosition,omitempty"`
		} `json:"document"`
	} `json:"lastAction"`
}
