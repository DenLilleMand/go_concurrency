package transactions

import (
	"github.com/denlillemand/go_concurrency_notes/strict_2pl_transactions/state"
)

const (
	Aborted   = "Aborted"
	Committed = "Committed"
	Active    = "Active"
	Abort     = "Abort"
	TooLate   = "TooLate"
	Success   = "Success"
)

type Handler struct {
	pendingTransactions chan Transaction
	stateHandler        *state.Handler
	validations         chan tx
	loopStop            chan stopSignal
	reset               chan resetSignal
}

type stopSignal struct {
	resp chan map[string]int
}

type resetSignal struct {
	resp  chan bool
	state map[string]int
}

type executeOperation struct {
	resp chan interface{}
}

type txComm struct {
	resp chan string
	msg  string
}

type tx struct {
	gtc       int
	result    interface{}
	readSet   map[string]interface{}
	changeSet map[string]int
	status    string
	comm      chan txComm
	err       error
}

func newHandler(State map[string]int) *Handler {
	pendingTransactions := make(chan Transaction)
	stateHandler := state.New(State)
	validations := make(chan tx)
	loopStop := make(chan stopSignal)
	reset := make(chan resetSignal)
	h := &Handler{
		pendingTransactions: pendingTransactions,
		stateHandler:        stateHandler,
		validations:         validations,
		loopStop:            loopStop,
		reset:               reset,
	}
	go h.validationHandler()
	go h.loop()
	return h
}

func (h *Handler) Execute(t Transaction) {
	h.pendingTransactions <- t
}

func (h *Handler) Stop() map[string]int {
	resp := make(chan map[string]int)
	h.loopStop <- stopSignal{resp}
	endState := <-resp
	return endState
}

func (h *Handler) Reset(state map[string]int) {
}

func Start(State map[string]int) *Handler {
	handler := newHandler(State)
	return handler
}

// txsHandler will receive a transaction on its pendingTransactions chan
// and start up a txHandler for it. It should keep a list of all of the
// channels that are contained so that we can do some kind of bulk operations,
// or maybe find a single transaction to do a operation.
func (h *Handler) txsHandler() {
	// maps indexed on the tx ID
	txs := map[int]tx{}
	finishedTxs := make(chan tx)
	validatedTxs := make(chan tx)
	for {
		select {
		case t := <-h.pendingTransactions:
			comm := make(chan txComm)
			txs[t.ID] = tx{comm: comm}
			go h.txHandler(t, finishedTxs, comm)
		case r := <-finishedTxs:
			// remove tx from active txs
			// move it into finished, it could have a sort of
			// status on it i guess we could update
			// but then we should not keep separate lists
			if r.err == nil {

			}
			h.validations <- r
		case validatedTx := <-validatedTxs:
			if tx, ok := txs[validatedTx.gtc]; ok {

			} else {
				panic("Smt happended")
			}
		}
	}
}

// transactionHandler has to execute the tx
// but also be able to receive a abort message that somehow
// kills the tx.
//
// Should be possible to wrap it in a function that executes, puts the result
// inside of a channel and sends it back. Question is if just aborting from this
// and letting the executing function realize it after it has been executed is
// good enough.
func (h *Handler) txHandler(t Transaction, resultComm chan tx, comm chan txComm) {
	resp := make(chan tx)
	go h.readPhase(t, resp)
	state := Active
	for {
		select {
		case result := <-resp:
			if state == Active {
				state = Committed
				resultComm <- result
			}
			// right now we do not handle the aborted case,
			// i guess this Tx will just never send a result to the
			// txsHandler, which is Ok
		case txMsg := <-comm:
			// Something that sucks with this one is that
			// we do not stop the actual call
			//
			// @TODO: Get a reference to the go routine and try to kill
			// it off somehow
			if txMsg.msg == Abort {
				if state == Committed {
					txMsg.resp <- TooLate
				} else if state == Active || state == Aborted {
					state = Aborted
					txMsg.resp <- Success
				}
			}
		}
	}
}

func (h *Handler) readPhase(t Transaction, resp chan tx) {
	readSet := h.stateHandler.Read(t.Reads)
	result, changeSet := t.Operation(readSet.Data)
	tx := tx{
		gtc:       readSet.GTC,
		result:    result,
		readSet:   t.Reads,
		changeSet: changeSet,
	}
	resp <- tx
}

func filterTxs(txs []tx, gtc int) []tx {
	filteredTxs := []tx{}
	for _, tr := range txs {
		if tr.gtc > gtc {
			tmp := tr
			filteredTxs = append(filteredTxs, tmp)
		}
	}
	return filteredTxs
}

func detectCollisions(tjs []tx, ti tx) bool {
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
	txs := []tx{}
	for {
		select {
		case tx := <-h.validations:
			possibleCollisions := filterTxs(txs, tx.gtc)
			if !detectCollisions(possibleCollisions, tx) {
				h.stateHandler.Write(tx.changeSet)
			}
		}
	}
}

// When reset is called we need to abort all on going transactions
//
// loop will be the main sync loop, what kind of problems could this cause?
// it feels a little like a bottleneck, to have it handle so many things.
// but maybe it is okay for a first iteration, if this one keeps track of
// all of the pending, active and comitted txs. It should at most handle
// the starting and eventual result of any tx, so something else should
// keep track of those things, so this one is just a messenger.
// Would maybe be easier if the readPhase did not send directly to the
// validationPhase, but went back to a central handler
func (h *Handler) loop() {
	for {
		select {
		case reset := <-h.reset:
			reset.state
		case stop := <-h.loopStop:
			endState := h.stateHandler.Stop()
			stop.resp <- endState
			break
		}
	}
}
