package funcs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sunshine-dev-code/rest/scopes"
	"github.com/sunshine-dev-code/rest/errs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchConditionMap map[string]string

var ConditionNameMap = SearchConditionMap{
	"eq":  "=",
	"gt":  ">",
	"gte": ">=",
	"lt":  "<",
	"lte": "<=",
	// "ne":  "!=",
	// "in":  "in",
	// "not_in": "not in",
	// "like":   "like",
	"prefix": "like", // 前缀匹配
}

func (scm *SearchConditionMap) GetCondition(key, value string) (string, string, error) {
	var (
		keySplit     = strings.Split(key, "_")
		length       = len(keySplit)
		keyWord      string
		conditionStr string
		condition    string
		ok           bool
	)
	if length == 1 {
		//字符串不包含下划线_,则字符串作为字段名
		keyWord, condition = key, "="
	} else {
		//字符串不包含下划线_,则取切分后的最后一段匹配条件
		// 若存在则取匹配出的值作为条件,key前面部分作为keyWord
		// 若不存在则将整体字符串作为keyWord
		conditionStr = keySplit[length-1]
		if condition, ok = (*scm)[conditionStr]; ok {
			keyWord = strings.Join(keySplit[:len(keySplit)-1], "_")
		} else {
			keyWord = key
			condition = "="
		}
	}

	if conditionStr == "prefix" {
		value = value + "%"
	}

	return fmt.Sprintf("%s %s ?", keyWord, condition), value, nil
}

func underscoreToLowerCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}
func checkKey(key string, model any) bool {
	var (
		typ = reflect.TypeOf(model)
	)
	key = underscoreToLowerCamelCase(key)

	switch typ.Kind() {
	case reflect.Struct:
	case reflect.Pointer:
		if typ.Elem().Kind() != reflect.Struct {
			return false
		}
	}
	_, ok := typ.FieldByName(key)
	return ok
}

func GetCondition[T any](key string, model T) (string, error) {
	var (
		keySplit     = strings.Split(key, "_")
		keyWord      string
		conditionStr string
		condition    string
		ok           bool
	)
	switch len(keySplit) {
	case 1:
		keyWord, conditionStr = key, "eq"
	case 2:
		keyWord, conditionStr = keySplit[0], keySplit[1]
	default:
		keyWord = strings.Join(keySplit[:len(keySplit)-2], "_")
		conditionStr = keySplit[len(keySplit)-1]
	}
	if ok = checkKey(keyWord, model); !ok {
		return "", errs.ErrParameter
	}

	if condition, ok = ConditionNameMap[conditionStr]; !ok {
		return "", errs.ErrParameter
	}
	return fmt.Sprintf("%s %s ?", keyWord, condition), nil
}

// type SearchParams struct {
// 	SearchParams     SearchParams
// 	SearchConditions []*db.Condition
// }

type SearchPageParams struct {
	PageParams
	SearchParams     SearchParams
	SearchConditions []*db.Condition
}

var sqlSymbolsAndKeywords = []string{
	",", ".", ";", "'", "\"", "`", "(", ")", "[", "]", "{", "}", "%", "*", "%", "<", "<=", ">", ">=", "<>", "!=",
	"LIKE", "BETWEEN", "AND", "OR", "SET", "DELETE", "DROP", "SELECT", "ALTER", "INSERT", "UPDATE", "CREATE", "INTO",
}

func (spp *SearchPageParams) ParseCondition() error {
	var (
		key       string
		condition string
		value     string
		err       error
	)
	for key, value = range spp.SearchParams {
		// 检查key内容，防止sql注入
		keyUpper := strings.ToUpper(key)
		for _, symbol := range sqlSymbolsAndKeywords {
			symbol := fmt.Sprintf(" %s ", symbol)
			if strings.Contains(keyUpper, symbol) {
				return errs.ErrParameter
			}
		}

		if condition, value, err = ConditionNameMap.GetCondition(key, value); err != nil {
			return err
		}
		spp.SearchConditions = append(spp.SearchConditions, &db.Condition{condition, value})
	}
	return nil
}

// 获取翻页与搜索参数
func NewSearchPageParamsFromContext(c *gin.Context) (*SearchPageParams, error) {
	var (
		params = NewPageParams()
		err    error
	)

	// 翻页参数
	//
	if err = c.ShouldBindQuery(&params.PageParams); err != nil {
		return nil, err
	}
	// 检查参数
	if params.PageParams.Size > 50 {
		return nil, errs.ErrParameter
	}
	// 设置默认参数
	if params.PageParams.Page == 0 {
		params.PageParams.Page = 1
	}
	if params.PageParams.Size == 0 {
		params.PageParams.Size = 10
	}

	// 搜索条件
	//
	if err = c.BindQuery(&params.SearchParams); err != nil {
		return nil, err
	}
	// 清理条件中出现的翻页参数
	delete(params.SearchParams, "page")
	delete(params.SearchParams, "size")
	delete(params.SearchParams, "order_by")

	return params, err
}

func GetSearchPageContent[T any](tx *gorm.DB,
	searchPageParams *SearchPageParams,
	scopes ...db.ScopeFunc) (total int64, result *[]T, err error) {
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
		searchPageParams.OrderBy=strings.ReplaceAll(searchPageParams.OrderBy, "+", " ")
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

func GetSearchPageDataWithTransaction[T any](tx *gorm.DB,
	searchPageParams *SearchPageParams,
	scopes ...db.ScopeFunc) (*PageResult, error) {
	var (
		models *[]T
		total  int64
		err    error
	)
	if total, models, err = GetSearchPageContent[T](tx, searchPageParams, scopes...); err != nil {
		return nil, err
	}
	result := NewPageResult(searchPageParams.PageParams, len(*models), int(total), models)
	return result, err
}

func GetSearchPageData[T any](searchPageParams *SearchPageParams, scopes ...db.ScopeFunc) (*PageResult, error) {
	var (
		result *PageResult
		err    error
	)
	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if result, err = GetSearchPageDataWithTransaction[T](tx, searchPageParams, scopes...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, err
}
