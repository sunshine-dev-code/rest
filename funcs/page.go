package funcs

import (
	"github.com/sunshine-dev-code/rest/db"
	"gorm.io/gorm"
)

type PageDataType = string

const (
	AvailableType PageDataType = "avaliable"
	DeletedType   PageDataType = "deleted"
	AllType       PageDataType = "all"
)

// params struct define
type PageParams struct {
	Page    int          `json:"page" form:"page"`         //页码
	Size    int          `json:"size" form:"size"`         //单页尺寸
	OrderBy string       `json:"order_by" form:"order_by"` //排序
	Type    PageDataType `json:"type" form:"type"`         //数据类型
}

type PageResult struct {
	PageParams
	// Page    int
	// Size    int
	Count int `json:"count"` //本页（Data）数据条数
	Total int `json:"total"` //所有数据条数
	// OrderBy string
	Data any `json:"data"`
}

func NewPageParams() *SearchPageParams {
	params := SearchPageParams{}
	params.SearchParams = make(SearchParams)
	return &params
}

func NewPageResult(pageParams PageParams, count, total int, data any) *PageResult {
	return &PageResult{pageParams, count, total, data}
}

func GetPageDataWithTransaction[T any](tx *gorm.DB, pageParam *PageParams, scopes ...db.ScopeFunc) (*PageResult, error) {
	var (
		models *[]T
		total  int64
		pageNo int
		err    error
	)

	if pageParam.OrderBy == "" {
		pageParam.OrderBy = "ID"
	}
	if pageParam.Size == 0 {
		pageParam.Size = 10
	}
	if pageParam.Page == 0 {
		pageParam.Page = 1
	}
	pageNo = pageParam.Page - 1
	tx = tx.Model(&models)
	for _, f := range scopes {
		tx = f(tx)
	}

	switch pageParam.Type {
	case DeletedType:
		tx = tx.Unscoped().Where("deleted_at is not null")
	case AllType:
		tx = tx.Unscoped()
	default:
		pageParam.Type = AvailableType
	}
	tx = tx.Count(&total).
		Limit(pageParam.Size).
		Offset(pageNo * pageParam.Size).
		Order(pageParam.OrderBy)
	tx = tx.Find(&models)
	err = tx.Error

	if err != nil {
		return nil, err
	}
	result := NewPageResult(*pageParam, len(*models), int(total), models)
	return result, err
}

func GetPageData[T any](pageParam *PageParams, scopes ...db.ScopeFunc) (*PageResult, error) {
	var (
		result *PageResult
		err    error
	)
	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if result, err = GetPageDataWithTransaction[T](tx, pageParam, scopes...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, err
}
