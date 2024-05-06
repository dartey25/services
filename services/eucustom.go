package services

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mdoffice/md-services/model"
)

type EuCustomService struct {
	db *sqlx.DB
}

func NewEuCustomService(db *sqlx.DB) *EuCustomService {
	return &EuCustomService{db: db}
}

func (s *EuCustomService) GetCountries() (countries []string, err error) {
	err = s.db.Select(&countries, "select distinct(t.cnt) from admin.aeo_dic t order by t.cnt")
	return
}

func (s *EuCustomService) GetAeoData(holder, country string, types []string, page, limit int) (data model.AeoPaginatedData, err error) {
	var whereTypesClause string
	var whereClause string

	if holder != "" {
		whereClause += fmt.Sprintf("upper(t.name) LIKE '%%%s%%'", strings.ToUpper(holder))
	}

	if country != "" {
		if whereClause != "" {
			whereClause += " and "
		}
		whereClause += fmt.Sprintf("upper(t.cnt) LIKE '%s%%'", strings.ToUpper(country))
	}

	if typeCnt := len(types); typeCnt == 1 {
		if whereClause != "" {
			whereClause += " and "
		}
		whereClause += fmt.Sprintf("upper(t.cert) LIKE '%s%%'", strings.ToUpper(types[0]))
	} else if typeCnt == 2 {
		for index, value := range types {
			if whereTypesClause != "" {
				whereTypesClause += " union "
			}
			whereTypesClause += fmt.Sprintf("select * from admin.aeo_dic t%v where upper(t%v.cert) LIKE '%s%%'", index, index, strings.ToUpper(value))
		}
	}

	query := "select t.* from "
	cntQuery := "select count(*) as TOTAL from "
	if whereTypesClause != "" {
		query += fmt.Sprintf("(%s)", whereTypesClause)
		cntQuery += fmt.Sprintf("(%s)", whereTypesClause)
	} else {
		query += "admin.aeo_dic"
		cntQuery += "admin.aeo_dic"
	}
	query += " t "
	cntQuery += " t "
	if whereClause != "" {
		query += fmt.Sprintf("where %s", whereClause)
		cntQuery += fmt.Sprintf("where %s", whereClause)
	}

	query += fmt.Sprintf(" order by to_date(t.effdate, 'dd/mm/yyyy') desc, t.cert, t.name offset %v rows fetch next %v rows only", page*limit, limit)
	fmt.Printf("QUERY: %s\n\nCNT:%v\n", query, cntQuery)

	err = s.db.Get(&data.Total, cntQuery)
	fmt.Printf("DATA: %v\n", data.Total)
	if err != nil || data.Total == 0 {
		return
	}

	err = s.db.Select(&data.Data, query)
	if err != nil {
		fmt.Printf("DB ERR %s\n", err.Error())
		return
	}
	data.Pages = data.Total / limit
	return
}

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
