### gflow
测试用例编排工具, yaml格式

##### yaml example
```yaml
# choose test case run mode
# serial: executed one by one
# parallel: run test case concurrently
# mode is case insensitive
mode: serial
# retry onece when fail
retry: 1 
# test case steps
workflow:
  # test case step
  - step:
      # step type
      type: nap
      # step label
      label: "auth api"
      # additional customized data that used to run test case step
      # the first character of the field name must be uppercase
      data: 
        # test case selector 
        Selector: "TODO"
        # test case uuid(mongo object id)
        Uuid: "5b4c4ff30207459803bfb3df"
  - step:
       type: nap
       label: "get user list"
       data:
         Uuid: "5b4c50210207459803bfb3e0"
```

##### run workflow
```shell
cd cmds/gwctl/
go build .
./gwctl run workflow-demo.yaml
```
