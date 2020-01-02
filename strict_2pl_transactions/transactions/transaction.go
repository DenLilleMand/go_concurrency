package transactions

type Transaction struct {
	ID int
	// Return writeState
	Operation func(map[string]int) (interface{}, map[string]int)
	Reads     map[string]interface{}
}

func NewTransaction(
	operation func(map[string]int) (interface{}, map[string]int),
	reads map[string]interface{},
	ID int) *Transaction {
	return &Transaction{
		ID:        ID,
		Operation: operation,
		Reads:     reads,
	}
}
