package entity

const (
	NotWhere = "noWhere"
)

// Condition Setter
type Condition struct {
	Key string
	Val interface{}
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
