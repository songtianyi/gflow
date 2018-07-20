package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/songtianyi/gflow"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type flow struct {
	Mode     string `yaml:"mode"`
	Retry    int    `yaml:"retry"`
	Workflow []struct {
		Step struct {
			Type  string      `yaml:"type"`
			Label string      `yaml:"label"`
			Data  interface{} `yaml:"data"`
		} `yaml:"step"`
	} `yaml:"workflow"`
}

func commandRunHandler(c *cli.Context) error {
	path := c.Args().Get(0)
	d, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %s", path, err.Error())
	}

	fileinfo, err := d.Stat()
	if err != nil {
		return fmt.Errorf("stat %s failed, %s", path, err)
	}

	if fileinfo.IsDir() {
		return fmt.Errorf("%s is dir", path)
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s failed, %s", path, err)
	}

	var f flow
	if err := yaml.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("umarshal %s failed, %s", path, err)
	}

	fmt.Println(string(b), "\n", f)
	// validate workflow
	if err := validate(f); err != nil {
		return err
	}
	wf := gflow.New(f.Mode, f.Retry)
	for _, s := range f.Workflow {
		switch strings.ToLower(s.Step.Type) {
		case "nap":
			err, nap := gflow.NewNapStep(s.Step.Label, c.String("uri"), s.Step.Data.(map[interface{}]interface{}))
			if err != nil {
				return fmt.Errorf("init nap step failed, %s", err)
			}
			wf.AddStep(nap)
			break
		default:
			return fmt.Errorf("step type %s not support", s.Step.Type)
		}
	}
	wf.Run()
	return nil
}

func validate(t flow) error {
	t.Mode = strings.ToLower(t.Mode)
	if t.Mode != gflow.SERIAL && t.Mode != gflow.CONCURRENT {
		return fmt.Errorf("invalid mode %s", t.Mode)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = "A cli tool to create and run workflow"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run workflow which described in yaml, if the ",
			Action:  commandRunHandler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "mongo, mgo",
					Value: "mongodb://localhost:27017",
					Usage: "mongo db uri",
				},
				cli.StringFlag{
					Name:  "uri, url",
					Value: "http://localhost:3000/api/nap/case/run",
					Usage: "nap-tcm service address",
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	return
}
