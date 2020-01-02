package main

import (
	"fmt"

	"github.com/denlillemand/go_concurrency_notes/strict_2pl_transactions/transactions"
)

func createTransactions() []transactions.Transaction {
	txs := []transactions.Transaction{}
	keys := []string{"a", "b"}
	for i := 0; i < 8000; i++ {
		keyIndex := i % 2
		key := keys[keyIndex]
		tReads := map[string]interface{}{
			key: nil,
		}
		tOperation := func(View map[string]int) (interface{}, map[string]int) {
			changeSet := map[string]int{}
			changeSet[key] = i
			result := "Done!"
			return result, changeSet
		}
		t := transactions.NewTransaction(tOperation, tReads, i)
		txs = append(txs, *t)
	}
	return txs

}

func main() {
	txs := createTransactions()
	initialState := map[string]int{
		"a": 0,
		"b": 0,
	}
	handler := transactions.Start(initialState)
	for _, tx := range txs {
		handler.Execute(tx)
	}
	endState := handler.Stop()
	fmt.Printf("endState = %+v\n", endState)
	for {

	}
}
