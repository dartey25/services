package model

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"time"
)

type AeoQueryParams struct {
	Page      int
	Limit     int
	Holder    string
	Country   string
	AuthTypes []string
}

type AeoPaginatedData struct {
	XMLName    xml.Name  `xml:"results"`
	Page       int       `json:"page" xml:"page"`
	Limit      int       `json:"-" xml:"-"`
	TotalItems int       `json:"totalItems" db:"TOTAL" xml:"totalItems"`
	TotalPages int       `json:"totalPages" xml:"totalPages"`
	Data       []AeoData `json:"data" db:"DATA" xml:"item"`
}

type AeoData struct {
	Id              int          `json:"id" db:"ID" xml:"id"`
	Holder          string       `json:"holder,omitempty" db:"NAME" xml:"holder"`
	HolderHighlight string       `xml:"-" db:"-" xml:"-"`
	IssCountry      string       `json:"issCountry,omitempty" db:"CNT" xml:"issCountry,omitempty"`
	CusCode         string       `json:"cusCode,omitempty" db:"CUSTOM" xml:"cusCode,omitempty"`
	AuthType        string       `json:"authType,omitempty" db:"CERT" xml:"authType,omitempty"`
	EffDate         string       `json:"effDate,omitempty" db:"EFFDATE" xml:"effDate,omitempty"`
	CreatedAt       time.Time    `json:"createdAt" db:"DATE_CREATE" xml:"-"`
	DeletedAt       sql.NullTime `json:"deletedAt,omitempty" db:"DATE_DELETE" xml:"-"`
}

type AeoType struct {
	Code        string
	Description string
	Checked     bool
}

type SoapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	XMLNS   string   `xml:"xmlns:soapenv,attr"`
	EoriNS  string   `xml:"xmlns:eori,attr"`
	Header  string   `xml:"soapenv:Header"`
	Body    SoapBody
}

type SoapBody struct {
	XMLName      xml.Name `xml:"soapenv:Body"`
	ValidateEORI ValidateEORI
}

type ValidateEORI struct {
	XMLName xml.Name `xml:"eori:validateEORI"`
	Eori    string   `xml:"eori:eori"`
}

type EoriEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    EoriBody `xml:"Body"`
}

type EoriBody struct {
	XMLName  xml.Name     `xml:"Body"`
	EORIResp EORIResponse `xml:"validateEORIResponse"`
}

type EORIResponse struct {
	XMLName xml.Name   `xml:"validateEORIResponse"`
	Return  EoriReturn `xml:"return"`
}

type EoriReturn struct {
	XMLName     xml.Name     `xml:"return"`
	RequestDate string       `xml:"requestDate"`
	Result      []EoriResult `xml:"result"`
}

type EoriResult struct {
	XMLName     xml.Name `xml:"result" json:"-"`
	EORI        string   `xml:"eori" json:"eori"`
	Status      int      `xml:"status" json:"status"`
	StatusDescr string   `xml:"statusDescr" json:"statusDescription"`
	Name        string   `xml:"name,omitempty" json:"name,omitempty"`
	Address     string   `xml:"address,omitempty" json:"address,omitempty"`
	Street      string   `xml:"street,omitempty" json:"street,omitempty"`
	PostalCode  string   `xml:"postalCode,omitempty" json:"postalCode,omitempty"`
	City        string   `xml:"city,omitempty" json:"city,omitempty"`
	Country     string   `xml:"country,omitempty" json:"country,omitempty"`
}

type EoriDisplayResult struct {
	Eori    string
	Status  int
	Name    string
	Address string
}

func NewEoriDisplayResult(data EoriResult) EoriDisplayResult {
	address := data.Street
	if data.PostalCode != "" {
		address += fmt.Sprintf(", %s", data.PostalCode)
	}
	if data.City != "" {
		address += fmt.Sprintf(", %s", data.City)
	}
	if data.Country != "" {
		address += fmt.Sprintf(", %s", data.Country)
	}

	return EoriDisplayResult{
		Eori:    data.EORI,
		Status:  data.Status,
		Name:    data.Name,
		Address: address,
	}
}
