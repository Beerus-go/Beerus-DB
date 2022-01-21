package builder

import "github.com/yuyenews/Beerus-DB/operation/entity"

// ConditionBuilder Used to build conditions more easily
type ConditionBuilder struct {
	conditions []*entity.Condition
}

// Create a ConditionBuilder
func Create() *ConditionBuilder {
	conditionBuilder := new(ConditionBuilder)
	conditionBuilder.conditions = make([]*entity.Condition, 0)
	return conditionBuilder
}

// Add conditions
func (conditionBuilder *ConditionBuilder) Add(key string, val ...interface{}) *ConditionBuilder {

	condition := new(entity.Condition)
	condition.Key = key
	condition.Val = val

	conditionBuilder.conditions = append(conditionBuilder.conditions, condition)
	return conditionBuilder
}

// Build get conditions
func (conditionBuilder *ConditionBuilder) Build() []*entity.Condition {
	return conditionBuilder.conditions
}
