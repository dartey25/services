package service

import (
	"fmt"
	"strings"

	"github.com/mdoffice/md-services/internal/eucustoms/model"
)

func (s *EuCustomService) GetCountries() (countries []string, err error) {
	err = s.db.Select(&countries, "select distinct(t.cnt) from admin.aeo_dic t order by t.cnt")
	return
}

func (s *EuCustomService) GetAeoData(q model.AeoQueryParams) (*model.AeoPaginatedData, error) {
	var whereTypesClause string
	var whereClause string

	if q.Holder != "" {
		whereClause += fmt.Sprintf("upper(t.name) LIKE '%%%s%%'", strings.ToUpper(q.Holder))
	}

	if q.Country != "" {
		if whereClause != "" {
			whereClause += " and "
		}
		whereClause += fmt.Sprintf("upper(t.cnt) LIKE '%s%%'", strings.ToUpper(q.Country))
	}

	if typeCnt := len(q.AuthTypes); typeCnt == 1 {
		if whereClause != "" {
			whereClause += " and "
		}
		whereClause += fmt.Sprintf("upper(t.cert) LIKE '%s%%'", strings.ToUpper(q.AuthTypes[0]))
	} else if typeCnt == 2 {
		for index, value := range q.AuthTypes {
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
	query += fmt.Sprintf(" order by to_date(t.effdate, 'dd/mm/yyyy') desc, t.cert, t.name offset %v rows fetch next %v rows only", (q.Page-1)*q.Limit, q.Limit)

	var data model.AeoPaginatedData
	err := s.db.Get(&data.TotalItems, cntQuery)
	if err != nil || data.TotalItems == 0 {
		return nil, err
	}
	data.TotalPages = data.TotalItems / q.Limit

	err = s.db.Select(&data.Data, query)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
