package main

import (
	"github.com/jaqmol/approx/axmsg"
	"github.com/jaqmol/approx/processorconf"
)

func main() {
	conf := processorconf.NewProcessorConf("approx_merge", []string{"PICK"})
	errMsg := axmsg.Errors{Source: "approx_merge"}

	if len(conf.Outputs) != 1 {
		errMsg.LogFatal(nil, "Merge expects exactly 1 output, but got %v", len(conf.Outputs))
	}
	if len(conf.Inputs) < 2 {
		errMsg.LogFatal(nil, "Merge expects more than 1 input, but got %v", len(conf.Inputs))
	}

	if len(conf.Envs["PICK"]) == 0 {
		errMsg.LogFatal(nil, "Merge expects value for env PICK")
	}

	af := NewApproxMerge(conf)
	af.Start()
}
