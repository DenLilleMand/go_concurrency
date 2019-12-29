package state

import "fmt"

type Handler struct {
	changeSets chan map[string]int
	readSets   chan readSetHandler
	GetState   chan int
}

type ReadSet struct {
	Data map[string]int
	GTC  int
}

type readSetHandler struct {
	resp    chan ReadSet
	readSet map[string]interface{}
}

func New() *Handler {
	changeSets := make(chan map[string]int)
	readSets := make(chan readSetHandler)
	getState := make(chan int)
	sh := &Handler{
		changeSets: changeSets,
		readSets:   readSets,
		GetState:   getState,
	}
	go sh.loop()
	return sh
}

func (sh *Handler) loop() {
	State := map[string]int{}
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
		case _ = <-sh.GetState:
			fmt.Printf("State = %+v\n", State)
		}
	}
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
