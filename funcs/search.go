package funcs

import (
	"github.com/sunshine-dev-code/rest/scopes"
	"gorm.io/gorm"
)

type SearchParams = map[string]string

func GetSearchContent[T any](tx *gorm.DB,
	searchPageParams *SearchPageParams,
	scopes ...scopes.ScopeFunc) (total int64, result *[]T, err error) {
	var (
		models = []T{}
	)
	tx = tx.Model(&models)
	// 处理搜索条件
	if err = searchPageParams.ParseCondition(); err != nil {
		return 0, nil, err
	}
	for _, condition := range searchPageParams.SearchConditions {
		tx = tx.Where(condition.Condition, condition.Value)
	}

	// 处理翻页参数
	tx = tx.Count(&total).
		Limit(searchPageParams.Size).
		Offset((searchPageParams.Page - 1) * searchPageParams.Size)

	// 处理排序
	if searchPageParams.OrderBy != "" {
		tx = tx.Order(searchPageParams.OrderBy)
	}

	// 处理传入的scopes函数
	for _, f := range scopes {
		tx = f(tx)
	}

	// 获取查询结果及错误信息
	if err = tx.Find(&models).Error; err != nil {
		return 0, nil, err
	}
	return total, &models, err
}

func GetSearchDataWithTransaction[T any](tx *gorm.DB,
	searchPageParams *SearchPageParams,
	scopes ...scopes.ScopeFunc) (*PageResult, error) {
	var (
		models *[]T
		total  int64
		err    error
	)
	if total, models, err = GetSearchContent[T](tx, searchPageParams, scopes...); err != nil {
		return nil, err
	}
	result := NewPageResult(searchPageParams.PageParams, len(*models), int(total), models)
	return result, err
}

func GetSearchData[T any](searchPageParams *SearchPageParams, scopes ...scopes.ScopeFunc) (*PageResult, error) {
	var (
		result *PageResult
		err    error
	)
	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if result, err = GetSearchDataWithTransaction[T](tx, searchPageParams, scopes...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, err
}
