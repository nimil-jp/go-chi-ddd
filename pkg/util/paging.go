package util

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Paging struct {
	page   int
	length int

	offset int
	limit  int
}

func NewPaging(r *http.Request) *Paging {
	paging := new(Paging)
	paging.page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if paging.page == 0 {
		paging.page = 1
	}
	paging.length, _ = strconv.Atoi(r.URL.Query().Get("length"))

	paging.offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	paging.limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	return paging
}

func (p *Paging) GetCount(query *gorm.DB, model interface{}) (uint, error) {
	var count int64

	copiedQuery := &gorm.DB{
		Config:       query.Config,
		Error:        query.Error,
		RowsAffected: query.RowsAffected,
		Statement:    query.Statement,
	}

	if err := copiedQuery.Model(model).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint(count), nil
}

func (p *Paging) Query() func(*gorm.DB) *gorm.DB {
	offset := (p.page - 1) * p.length
	limit := p.length

	if p.offset > 0 {
		offset = p.offset
	}

	if p.limit > 0 {
		limit = p.limit
	}

	return func(db *gorm.DB) *gorm.DB {
		db.Offset(offset).Limit(limit)
		return db
	}
}
