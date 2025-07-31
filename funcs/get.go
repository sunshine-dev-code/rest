package funcs

import (
	"github.com/sunshine-dev-code/rest/scopes"

	"gorm.io/gorm"
)

func GetDataWithTransaction[T any](tx *gorm.DB, scopes ...scopes.ScopeFunc) (*T, error) {
	var (
		result T
		err    error
	)
	for _, scopeFunc := range scopes {
		tx = scopeFunc(tx)
	}
	if err = tx.First(&result).Error; err != nil {
		return nil, err
	}
	return &result, err
}

func GetData[T any](scopes ...scopes.ScopeFunc) (*T, error) {
	var (
		result *T
		err    error
	)
	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if result, err = GetDataWithTransaction[T](tx, scopes...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, err
}

func GetBatchDataWithTransaction[T any](tx *gorm.DB, scopes ...scopes.ScopeFunc) (*[]T, error) {
	var (
		result []T
		err    error
	)
	for _, scopeFunc := range scopes {
		tx = scopeFunc(tx)
	}
	if err = tx.Find(&result).Error; err != nil {
		return nil, err
	}
	return &result, err
}

func GetBatchData[T any](scopes ...scopes.ScopeFunc) (*[]T, error) {
	var (
		result *[]T
		err    error
	)
	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if result, err = GetBatchDataWithTransaction[T](tx, scopes...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, err
}
