package gflow

import (
	"fmt"
	"net/http"
	"strings"
)

type NapData struct {
	Selector string
	Uuid     string
}

type NapStep struct {
	label  string
	uri    string
	data   *NapData
	client *http.Client
}

func NewNapStep(label, uri string, data map[interface{}]interface{}) (error, Step) {
	step := &NapStep{
		label:  label,
		uri:    strings.TrimSuffix(uri, "/"),
		data:   &NapData{},
		client: &http.Client{},
	}
	// interface to struct
	if err := Fill(step.data, data); err != nil {
		return err, nil
	}
	return nil, step
}
func (s *NapStep) Label() string {
	return s.label
}

func (s *NapStep) UUID() string {
	if s.data != nil {
		return s.data.Uuid
	}
	return "UUID"
}

func (s *NapStep) init() error {
	return nil
}

func (s *NapStep) Run(context Context) error {
	// init
	if err := s.init(); err != nil {
		return err
	}
	if len(s.data.Uuid) > 1 {
		req, err := http.NewRequest("GET", s.uri+"/"+s.data.Uuid, nil)
		if err != nil {
			return err
		}
		res, err := s.client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("status:%d", res.StatusCode)
		}
	} else {
		return fmt.Errorf("not uuid")
	}
	return nil
}

func (s *NapStep) OnFailure(err error, context Context) error {
	// throw
	return err
}
