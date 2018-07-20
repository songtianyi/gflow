package gflow

import (
	"fmt"
	"strings"
)

const (
	SERIAL     = "serial"
	CONCURRENT = "concurrent"
)

type Context interface{}

type Workflow struct {
	Mode      string
	OnFailure FailureFunc
	Context   Context

	queue []Step
}

func New(mode string, retry int) *Workflow {
	w := &Workflow{
		Mode:  strings.ToLower(mode),
		queue: make([]Step, 0),
	}
	if retry > 0 {
		w.OnFailure = RetryFailure(retry)
	}
	return w
}

func (w *Workflow) Run() error {
	if w.Mode == SERIAL {
		for i, step := range w.queue {
			prefix := fmt.Sprintf("STEP(%d)(%s):", i, step.Label())
			if err := step.Run(w.Context); err != nil {
				if err := step.OnFailure(err, w.Context); err != nil {
					if err := w.OnFailure(err, step, w.Context); err != nil {
						fmt.Println(prefix, "[FAILED]", err)
						return err
					}
				}
			}
			fmt.Println(prefix, "COMPLETE")
		}
	}
	return nil
}

func (w *Workflow) AddSteps(steps []Step) {
	w.queue = append(w.queue, steps...)
}

func (w *Workflow) AddStep(step Step) {
	w.queue = append(w.queue, step)
}
