package main

import (
	"io"

	"github.com/jaqmol/approx/axenvs"
	"github.com/jaqmol/approx/axmsg"
)

// NewApproxMerge ...
func NewApproxMerge(envs *axenvs.Envs) *ApproxMerge {
	errMsg := &axmsg.Errors{Source: "approx_merge"}
	pickEnv := envs.Required["PICK"]
	var pick Pick
	if "as_comes" == pickEnv {
		pick = PickAsComes
	} else if "round_robin" == pickEnv {
		pick = PickRoundRobin
	} else {
		errMsg.LogFatal(nil, "Merge expects env PICK to be either as_comes or round_robin, but got %v", pickEnv)
	}
	ins, outs := envs.InsOuts()
	return &ApproxMerge{
		errMsg:     errMsg,
		output:     axmsg.NewWriter(&outs[0]),
		inputs:     axmsg.NewReaders(ins),
		pick:       pick,
		inputIndex: 0,
	}
}

// ApproxMerge ...
type ApproxMerge struct {
	errMsg     *axmsg.Errors
	output     *axmsg.Writer
	inputs     []axmsg.Reader
	pick       Pick
	inputIndex int
}

// Pick ...
type Pick int

// Pick Types
const (
	PickAsComes Pick = iota
	PickRoundRobin
)

// Start ...
func (a *ApproxMerge) Start() {
	msgChan := make(chan []byte)
	if a.pick == PickAsComes {
		a.pickAsComes(msgChan)
	} else if a.pick == PickRoundRobin {
		go a.pickRoundRobin(msgChan)
	}

	for msgBytes := range msgChan {
		err := a.output.WriteBytes(msgBytes)
		if err != nil {
			a.errMsg.Log(nil, "Error writing response to output: %v", err.Error())
			return
		}
	}
}

func (a *ApproxMerge) pickAsComes(msgChan chan<- []byte) {
	for i := 0; i < len(a.inputs); i++ {
		input := a.inputs[i]
		go a.pickAsComesFromInput(&input, msgChan)
	}
}

func (a *ApproxMerge) pickAsComesFromInput(input *axmsg.Reader, msgChan chan<- []byte) {
	var hardErr error
	for hardErr == nil {
		var msgBytes []byte
		msgBytes, hardErr = input.ReadBytes()
		if hardErr != nil {
			break
		}
		msgChan <- msgBytes
	}
	a.processHardInputReadingError(hardErr)
	close(msgChan)
}

func (a *ApproxMerge) pickRoundRobin(msgChan chan<- []byte) {
	var hardErr error
	for hardErr == nil {
		if a.inputIndex >= len(a.inputs) {
			a.inputIndex = 0
		}
		input := a.inputs[a.inputIndex]
		a.inputIndex++

		var msgBytes []byte
		msgBytes, hardErr = input.ReadBytes()
		if hardErr != nil {
			break
		}
		msgChan <- msgBytes
	}
	a.processHardInputReadingError(hardErr)
	close(msgChan)
}

func (a *ApproxMerge) processHardInputReadingError(hardErr error) {
	if hardErr == io.EOF {
		a.errMsg.LogFatal(nil, "Unexpected EOL listening for response input")
	} else {
		a.errMsg.LogFatal(nil, "Unexpected error listening for response input: %v", hardErr.Error())
	}
}
