package onepassword

import (
	"os/exec"
	"strings"
	"sync"
)

type mockOnePassConfig struct {
	runCmd             func() (string, error)
	execCommandResults []string // Populated when execCommand is executed so you can assert against the args passed in
}

func mockOnePassClient(params *mockOnePassConfig) *OnePassClient {
	ret := &OnePassClient{
		PathToOp: "op",
		mutex:    &sync.Mutex{},
	}

	if params.runCmd != nil {
		ret.execCommand = func(binary string, args ...string) *exec.Cmd {
			params.execCommandResults = append([]string{binary}, args...)

			out, err := params.runCmd()
			if err != nil {
				return exec.Command("sh", "-c", "echo "+strings.ReplaceAll(err.Error(), `"`, `\"`)+" && false")
			}
			return exec.Command("sh", "-c", "echo "+strings.ReplaceAll(out, `"`, `\"`))
		}
	}

	return ret
}
