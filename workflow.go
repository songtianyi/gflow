package gflow

import (
	"fmt"
	"strings"
	"time"
)

const (
	SERIAL     = "serial"
	CONCURRENT = "concurrent"
)

type Context interface{}

type Workflow struct {
	Mode        string
	OnFailure   FailureFunc
	Context     Context
	Description string
	FileNmae    string

	queue []Step
}

func New(mode string, retry int, desc string, filename string) *Workflow {
	w := &Workflow{
		Mode:        strings.ToLower(mode),
		queue:       make([]Step, 0),
		Description: desc,
		FileNmae:    filename,
	}
	w.OnFailure = RetryFailure(retry)
	return w
}

func (w *Workflow) Run() error {
	if w.Mode == SERIAL {
		for i, step := range w.queue {
			// output format
			// yaml-index-uuid-label-start-end-status-error
			prefix := fmt.Sprintf("%s-%d-%s-%s-%d", w.FileNmae, i, step.UUID(), step.Label(), time.Now().Unix())
			if err := step.Run(w.Context); err != nil {
				if err := step.OnFailure(err, w.Context); err != nil {
					if err := w.OnFailure(err, step, w.Context); err != nil {
						fmt.Println(fmt.Sprintf("%s-%d-FAILED-%s", prefix, time.Now().Unix(), err))
						return err
					}
				}
			}
			fmt.Println(fmt.Sprintf("%s-%d-COMPLETE-%s", prefix, time.Now().Unix(), "OK"))
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
