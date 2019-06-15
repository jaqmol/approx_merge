package main

import (
	"github.com/jaqmol/approx/axenvs"
	"github.com/jaqmol/approx/axmsg"
)

func main() {
	envs := axenvs.NewEnvs("approx_merge", []string{"PICK"}, nil)
	errMsg := axmsg.Errors{Source: "approx_merge"}

	if len(envs.Outs) != 1 {
		errMsg.LogFatal(nil, "Merge expects exactly 1 output, but got %v", len(envs.Outs))
	}
	if len(envs.Ins) < 2 {
		errMsg.LogFatal(nil, "Merge expects more than 1 input, but got %v", len(envs.Ins))
	}

	if len(envs.Required["PICK"]) == 0 {
		errMsg.LogFatal(nil, "Merge expects value for env PICK")
	}

	af := NewApproxMerge(envs)
	af.Start()
}
