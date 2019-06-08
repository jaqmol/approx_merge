package main

import (
	"bufio"
	"io"

	"github.com/jaqmol/approx/errormsg"
	"github.com/jaqmol/approx/processorconf"
)

// TODO: test & debug

// NewApproxMerge ...
func NewApproxMerge(conf *processorconf.ProcessorConf) *ApproxMerge {
	errMsg := &errormsg.ErrorMsg{Processor: "approx_merge"}
	pickEnv := conf.Envs["PICK"]
	var pick Pick
	if "as_comes" == pickEnv {
		pick = PickAsComes
	} else if "round_robin" == pickEnv {
		pick = PickRoundRobin
	} else {
		errMsg.LogFatal(nil, "Merge expects env PICK to be either as_comes or round_robin, but got %v", pickEnv)
	}
	return &ApproxMerge{
		errMsg:     errMsg,
		conf:       conf,
		output:     conf.Outputs[0],
		inputs:     conf.Inputs,
		pick:       pick,
		inputIndex: 0,
	}
}

// ApproxMerge ...
type ApproxMerge struct {
	errMsg     *errormsg.ErrorMsg
	conf       *processorconf.ProcessorConf
	output     *bufio.Writer
	inputs     []*bufio.Reader
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
		_, err := a.output.Write(msgBytes)
		if err != nil {
			a.errMsg.Log(nil, "Error writing response to output: %v", err.Error())
			return
		}
		err = a.output.Flush()
		if err != nil {
			a.errMsg.Log(nil, "Error flushing response to output: %v", err.Error())
			return
		}
	}
}

func (a *ApproxMerge) pickAsComes(msgChan chan<- []byte) {
	for i := 0; i < len(a.inputs); i++ {
		input := a.inputs[i]
		go a.pickAsComesFromInput(input, msgChan)
	}
}

func (a *ApproxMerge) pickAsComesFromInput(input *bufio.Reader, msgChan chan<- []byte) {
	var hardErr error
	for hardErr == nil {
		var msgBytes []byte
		msgBytes, hardErr = input.ReadBytes('\n')
		if hardErr != nil {
			return
		}
		msgChan <- msgBytes
	}
	a.processHardInputReadingError(hardErr)
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
		msgBytes, hardErr = input.ReadBytes('\n')
		if hardErr != nil {
			return
		}
		msgChan <- msgBytes
	}
	a.processHardInputReadingError(hardErr)
}

func (a *ApproxMerge) processHardInputReadingError(hardErr error) {
	if hardErr == io.EOF {
		a.errMsg.LogFatal(nil, "Unexpected EOL listening for response input")
	} else {
		a.errMsg.LogFatal(nil, "Unexpected error listening for response input: %v", hardErr.Error())
	}
}
