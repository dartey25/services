package model

import (
	"encoding/xml"
)

type ApiResponse struct {
	XMLName xml.Name    `xml:"root" json:"-"`
	Data    interface{} `xml:"results>_" json:"data"`
	Error   string      `xml:"Error,omitempty" json:"message,omitempty"`
}

func NewApiResponse(data interface{}) ApiResponse {
	return ApiResponse{Data: data}
}

func NewApiError(message string) ApiResponse {
	return ApiResponse{Error: message}
}
