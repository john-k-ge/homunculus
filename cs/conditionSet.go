package cs

type ConditionSet interface {
	//New() *Cache
	SaveCondition(string, int64) error
	CheckCondition(string) (int64, error)
	ConditionExists(string) (bool, error)
	DeleteCondition(string) error
	ClearConditions() error
	IncrementCondition(string) (int64, error)
	Size() (int64, error)
}
