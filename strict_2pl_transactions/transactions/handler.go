package transactions

import (
	"github.com/denlillemand/go_concurrency_notes/strict_2pl_transactions/state"
)

type Handler struct {
	pendingTransactions chan Transaction
	StateHandler        *state.Handler
	validations         chan tResult
}

type tResult struct {
	gtc       int
	result    interface{}
	readSet   map[string]interface{}
	changeSet map[string]int
}

func NewHandler() *Handler {
	pendingTransactions := make(chan Transaction)
	stateHandler := state.New()
	validations := make(chan tResult)
	h := &Handler{
		pendingTransactions: pendingTransactions,
		StateHandler:        stateHandler,
		validations:         validations,
	}
	go h.validationHandler()
	go h.loop()
	return h
}

func (h *Handler) Execute(t Transaction) {
	h.pendingTransactions <- t
}

func (h *Handler) readPhase(t Transaction) {
	readSet := h.StateHandler.Read(t.Reads)
	result, changeSet := t.Operation(readSet.Data)
	tResult := tResult{
		gtc:       readSet.GTC,
		result:    result,
		readSet:   t.Reads,
		changeSet: changeSet,
	}
	h.validations <- tResult
}

func filterTResults(tResults []tResult, gtc int) []tResult {
	filteredTResults := []tResult{}
	for _, tr := range tResults {
		if tr.gtc > gtc {
			tmp := tr
			filteredTResults = append(filteredTResults, tmp)
		}
	}
	return filteredTResults
}

func detectCollisions(tjs []tResult, ti tResult) bool {
	for _, tj := range tjs {
		for k1, _ := range tj.changeSet {
			for k2, _ := range ti.readSet {
				if k1 == k2 {
					return true
				}
			}
		}
	}
	return false
}

// all transactions that want to commit, sends it's GTC, result and readSet
// to the validation phase, this handler should maintain state about all
// transactions that has occurred and active
func (h *Handler) validationHandler() {
	tResults := []tResult{}
	for {
		select {
		case tResult := <-h.validations:
			possibleCollisions := filterTResults(tResults, tResult.gtc)
			if !detectCollisions(possibleCollisions, tResult) {
				h.StateHandler.Write(tResult.changeSet)
			}
		}
	}
}

func (h *Handler) loop() {
	for {
		select {
		case t := <-h.pendingTransactions:
			go h.readPhase(t)
		}
	}
}
