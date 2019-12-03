package main

type Transaction struct {
	ID int
	// Return writeState
	Operation func(map[string]int) (interface{}, map[string]int)
	Reads     map[string]interface{}
}

func State(changeSets chan map[string]int) {
	State := map[string]int{}
	for {
		changeSet := <-changeSets
		for k, v := range changeSet {
			State[k] = v
		}
	}
}

func TransactionHandler(transaction Transaction) {
	for {

	}
}

func TransactionsHandler(pendingTransactions chan Transaction) {
	transactions := map[int]Transaction{}
	for {
		transaction := <-pendingTransactions
		transaction.ID = len(transactions) + 1
		transactions[transaction.ID] = transaction
		// Read the data here mb, then the ID will work i guess
		// View :=
		go TransactionHandler(transaction)
	}
}

func main() {
	transactions := make(chan Transaction)
	t1 := Transaction{}
	t1.Reads = map[string]interface{}{
		"r2d2":           nil,
		"luke skywalker": nil,
	}
	t1.Operation = func(View map[string]int) (interface{}, map[string]int) {
		changeSet := map[string]int{}
		changeSet["r2d2"] = 2
		result := "Done!"
		return result, changeSet
	}
	go TransactionsHandler(transactions)
	transactions <- t1
	for {
	}
}
