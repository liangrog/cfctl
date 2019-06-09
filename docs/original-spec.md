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

### Provide diff facility
Provide diff to stack before applying the changes. Default requires confirmation. Allow pass through with flag.

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
    |- stacks
       |- jump-host.yaml
       |- web-server.yaml
        ...
    |- deployments
       |- foo
          |- env
              |- default
                 |- vars.yaml
              |- dev
                 |- secret.yaml
                 |- env.yaml
              |- prod
                 |- secret.yaml
                 |_ env.yaml
          |- params
             |- jump-host
                |- default.yaml
                |- ap-southeast-2.yaml
          |- stacks.yaml
```
`env` sub-folder convention:
* `default` folder contains yaml files that will be applied to any environment as the default.
* Specific folder given by the command line will override values defined by `default` folder.
* ansible-vault encrypted file will be auto dectected and decrypted during variable merger.


`templates` folder contains generic purpose templates. They can be used individually or being use by modules. Multi-level folder is allowed. The templates in this filder can be referred by using in the parameters in the format of `{{ template file/path/in/template/folder }}`.

`modules` folder contains modules templates that represent a infrastructure function. It can be used individually or being referred by other stacks in the format of `{{ stack file/path/in/modules/folder }}` in the parameters file.

secret can be used in parameters in the form of `{{ secret key-in-secret-file }}`. secret file cannot use any helpers.
environment values can use secret helper.
parameters file can use template and stack helpers.

`{{ stackout stack-output-key }}` is the helper to get a output value of a stack

For multi-region, create a seperate parameters file and name it with the AWS region name. The values in this file will override the default parameters file if exists.

stacks.yaml contains stacks for the foo deployments. It will have creation orders, otherwise it will concurrently create the stacks.


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
