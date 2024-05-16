package service

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mdoffice/md-services/pkg/eucustoms/model"
)

func (s *EuCustomService) ValidateEori(code string) (data []model.EoriResult, err error) {
	soapReq := model.SoapEnvelope{
		XMLNS:  "http://schemas.xmlsoap.org/soap/envelope/",
		EoriNS: "http://eori.ws.eos.dds.s/",
		Header: "",
		Body: model.SoapBody{
			ValidateEORI: model.ValidateEORI{
				Eori: strings.ToUpper(code),
			},
		},
	}

	xmlData, err := xml.Marshal(soapReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://ec.europa.eu/taxation_customs/dds2/eos/validation/services/validation", "text/xml", bytes.NewBuffer(xmlData))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 500 {
		return nil, errors.New("internal server error on eucustom soap endpoint")
	}

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var envelope model.EoriEnvelope
	err = xml.Unmarshal(resBody, &envelope)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return envelope.Body.EORIResp.Return.Result, nil
}
