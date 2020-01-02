package state

import "fmt"

type Handler struct {
	changeSets chan map[string]int
	readSets   chan readSetHandler
	GetState   chan getState
	stop       chan stopSignal
}

type stopSignal struct {
	resp    chan map[string]int
	success bool
}

type getState struct {
	resp chan map[string]int
}

type ReadSet struct {
	Data map[string]int
	GTC  int
}

type readSetHandler struct {
	resp    chan ReadSet
	readSet map[string]interface{}
}

func New(initialState map[string]int) *Handler {
	changeSets := make(chan map[string]int)
	readSets := make(chan readSetHandler)
	getState := make(chan getState)
	stop := make(chan stopSignal)
	sh := &Handler{
		changeSets: changeSets,
		readSets:   readSets,
		GetState:   getState,
		stop:       stop,
	}
	go sh.loop(initialState)
	return sh
}

func (sh *Handler) loop(initialState map[string]int) {
	State := initialState
	GTC := 0
	totalUpdates := 0
	for {
		select {
		case changeSet := <-sh.changeSets:
			for k, v := range changeSet {
				State[k] = v
			}
			totalUpdates = totalUpdates + 1
			fmt.Printf("totalUpdates = %+v\n", totalUpdates)
		case read := <-sh.readSets:
			data := ReadSet{
				GTC: GTC,
			}

			readState := map[string]int{}
			for k, v := range State {
				readState[k] = v
			}
			read.resp <- data
			GTC += 1
		case getState := <-sh.GetState:
			getState.resp <- State
		case stop := <-sh.stop:
			stop.resp <- State
			break
		}
	}
}

func (sh *Handler) Get() map[string]int {
	resp := make(chan map[string]int)
	getState := getState{resp}
	sh.GetState <- getState
	state := <-resp
	return state
}

func (sh *Handler) Stop() map[string]int {
	resp := make(chan map[string]int)
	stop := stopSignal{resp, true}
	sh.stop <- stop
	endState := <-resp
	return endState
}

// Need a API for reading the state of state given a readset
func (sh *Handler) Read(readSet map[string]interface{}) ReadSet {
	resp := make(chan ReadSet)
	rsh := readSetHandler{resp: resp, readSet: readSet}
	sh.readSets <- rsh
	return <-resp
}

func (sh *Handler) Write(changeSet map[string]int) {
	sh.changeSets <- changeSet
}
