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
By default, parent stack parameters override children stack parameters.
Command line parameters overrides any parameters given in the files.

For environment parameter values such as development, staging, production, end users should use `--parameters` flag to provide values that specific to the environment or use `--env` flag to specify the value file for the environment.

### Unit tested all components
All components must be unit tested

### Support multi-region deployments
Provide facility that can apply the same CloudFormation or changes to multple regions in one command

### AWS profile management
- Allow using environment variable
- Allow using AWS profile
- Handle MFA
- Allow profile configured in cfctl without the needs to install awscli
- Ordering: ENV > profile > cfctl configuration

Load order:
<1> cfctl credential provider
<2> env crednetial provider
<3> shared credential provider
<4> remote credential provider (ec2 or ec2 roles)

<2> to <4> are provided by AWS SDK default chained credential provider

### Folder structure
```  
  - Your Git Repo
    |- templates
        |- ec2-general.yaml
        |- ec2-kubernetes.yaml
        |- s3-encrypted.yaml
        ...
    |- modules
       |- jump-host.yaml
       |- web-server.yaml
        ...
    |- project
       |- rocket-man
          |- env
              |- dev.yaml
              |- stage.yaml
              |- prod.yaml
          |- params.yaml
          |- packages.yaml
          |- config.yaml
       |- super-woman
          |- params.yaml
          |- package.yaml
          |- config.yaml
```

### config.yaml
The is the configurations file sets default values for cfctl

### params.yaml
This file contains all stack specific default parameter-value pairs

### package.yaml
This is the dependancy management file.

- This file does not track package versions
- Package must be local to the project at this stage. Inter-repo packaging will be considered in the future when there are enough user demands.

### ~/.cfctl/
This is the directory for cfctl management.

### Commands
1. AWS Profile
```
  $ cfctl profile create
  
  $ cfctl profile update
  
  $ cfctl profile delete
```  
2. CloudFormation
```
  $ cfctl stack create
  
  $ cfctl stack update
  
  $ cfctl stack delete
  
  $ cfctl stack output
  
  $ cfctl stack validate
  
  $ cfctl stack list
  
  $ cfctl stack get
```  
### Flags

Global
```
  --profile 
  --region
```


Trello Board
---
[github-cfctl](https://trello.com/b/3etT9edo/github-cfctl)
