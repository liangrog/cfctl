# `cfctl` - a Simple AWS CloudFormation Stack Lifecycle Management Tool
[![Version](https://img.shields.io/github/v/release/liangrog/cfctl)](https://github.com/liangrog/cfctl/releases)
[![GoDoc](https://godoc.org/github.com/liangrog/cfctl?status.svg)](https://godoc.org/github.com/liangrog/cfctl)
![](https://github.com/liangrog/cfctl/workflows/Development/badge.svg)
![](https://github.com/liangrog/cfctl/workflows/Release/badge.svg)

----

cfctl is a command line tool which helps to facilitate and manage AWS CloudFormation stack lifecycle. 

It supports a simple and highly flexible repository structure for organising Cloudformation templates, parameters and environments. 

The reason for creating this tool you can read in [medium](https://itnext.io/from-lmdo-to-cfctl-the-journey-of-developing-a-devops-tool-13b5d3ba211e)

----

## Features

- No tool lock-in. You can switch back to using awscli without much hassle and vice versa.
- Supports parameter files in YAML format and provides many useful functions for parameter file templating.
- Provide file ecryption facility for secrets used in parameter templates and automatically decrypt them during deployment.
- Configuration over convention. Provide high flexibility to suit different needs in directory structures for manage templates, paramters and environment specific files.
- Auto stack order sorting during deployment based on dependancy.
- Auto detect circular dependency amongst deploying stacks.
- Automatically uploading nested stacks during deployment and return those stack urls for referencing.
- Dynamically retrieving stack outputs for stacks that referencing them.


## API References  
[Please refer to the GoDoc](https://godoc.org/github.com/liangrog/cfctl)

## Getting Started
### Installing
1. Download the desired version base on your OS from the [releases page](https://github.com/liangrog/cfctl/releases)
2. Move it to the executables folder. For example for linux amd64: `chmod +x cfctl-linux-amd64 && sudo mv cfctl-linux-amd64 /usr/local/bin/cfctl`


### Setup AWS Credentials
cfctl piggybacks your existing [AWSCLI](https://aws.amazon.com/cli/) credential setting. If you don't have one, there are a few options:
- Use awscli environment [variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)

or

- Create two files: `~/.aws/credentials` and `~/.aws/config` as per [instruction](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). 

As a minimum, your IAM user must have permission to create S3 bucket. In addition, you will need permissions for AWS resources that your Cloudformation requires.

### Enabling Shell Autocompletion
BASH

```
source <(cfctl completion bash) # setup autocomplete in bash into the current shell, bash-completion package should be installed first.
echo "source <(cfctl completion bash)" >> ~/.bashrc # add autocomplete permanently to your bash shell.
```

You can also use a shorthand alias for cfctl that also works with completion:
```
alias cf=cfctl
complete -F __start_cfctl cf
```

ZSH

```
source <(cfctl completion zsh)  # setup autocomplete in zsh into the current shell
echo "if [ $commands[cfctl] ]; then source <(cfctl completion zsh); fi" >> ~/.zshrc # add autocomplete permanently to your zsh shell
```

### Quick Start
1. Create a sample repository by running below command. The Command will create a default repository structure (which you can change it to your liking).
```
$ cfctl init
$
$ tree cfctl-sample
cfctl-sample
├── deploy
│   └── sample
│       ├── environments
│       │   └── default
│       │       └── var.yaml
│       ├── parameters
│       │   └── s3.yaml
│       └── stacks.yaml
└── templates
    └── s3-encrypted.yaml
```

2. Deploy example stack which creates a new s3 bucket.
```
$ cd cfctl-sample/deploy/sample
$ cfctl stack deploy
``` 

3. Clean up
```
$ cfctl stack delete --all
```

## Documentations
- [Stack Configuration File](docs/config.md)
- [Directory Structure](docs/directory.md)
- [Parameter file template](docs/parameters.md)
- [Secrects Encryption](docs/secrets.md)
- [Cheatsheet](docs/cheatsheet.md)
