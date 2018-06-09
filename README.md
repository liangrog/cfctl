# cfctl
AWS CloudFormation DevOp tool

## Design Principles
1. Retain CloudFormation's independence to the tool itself
2. Dynamic, on-the-fly state management without the need to use persistent media
3. Facilitate modularity

## Requirements

### Allow multiple sources of parameters
- Environment variable
- Stack output
- Local CloudFormation

Multiple sources for one parameter is allow

### Providing parameter scoping
By default, parent parameter will be overriden by children parameter.
If given `--override` flag, parent parameter override children parameter.

### Unit tested all components
All components must be unit tested

### Support multi-region deployments
Provide facility that can apply the same CloudFormation or changes to multple regions in one command

### Folder structure
```  
  - project
    |- modules
       |- templateA.yaml
       |- templateB.yaml
       |- folderA
       |- templateC.yaml
       |- templateD.yaml
       |- folderB
          |- templateE.yaml
          ...
    |- stacks
       |- global.yaml
       |- stackA
          |- params.yaml
          |- config.yaml
```
### global.yaml

### config.yaml

### Commands
  
Trello Board
---
[github-cfctl](https://trello.com/b/3etT9edo/github-cfctl)
