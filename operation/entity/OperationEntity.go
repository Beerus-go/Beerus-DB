package entity

const (
	NotWhere = "6ca6d99a-2ca3-4734-921d-f3718bb7e179"
)

// Condition Setter
type Condition struct {
	Key string
	Val []interface{}
}

// GetCondition get Condition
func GetCondition(key string, val ...interface{}) *Condition {
	condition := new(Condition)
	condition.Key = key
	condition.Val = val
	return condition
}

// PageParam Paging entity
type PageParam struct {
	CurrentPage int
	PageSize    int
	Params      map[string]interface{}
}

// PageResult Pagination back to the entity
type PageResult struct {
	CurrentPage int
	PageSize    int
	PageCount   int
	PageTotal   int
	DataList    []map[string]string
}

// CalcPageTotal Calculate the total number of pages
func CalcPageTotal(pageSize int, pageCount int) int {
	if pageCount%pageSize > 0 {
		return (pageCount / pageSize) + 1
	}
	return pageCount / pageSize
}
