package main

import "github.com/denlillemand/go_concurrency_notes/strict_2pl_transactions/transactions"

func createTransactions() []transactions.Transaction {
	txs := []transactions.Transaction{}
	keys := []string{"a", "b"}
	for i := 0; i < 4000; i++ {
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
		t := transactions.NewTransaction(tOperation, tReads)
		txs = append(txs, *t)
	}
	return txs

}

func main() {
	txs := createTransactions()

	handler := transactions.NewHandler()
	for _, tx := range txs {
		handler.Execute(tx)
	}
	handler.StateHandler.GetState <- 1
	for {

	}
}
