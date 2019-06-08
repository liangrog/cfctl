# `cfctl` - a Simple AWS CloudFormation DevOps Tool
cfctl is a streamline command line utility that helps to organise and manage AWS stacks that created by using CloudFormation. 

## Reason of Creation
**TL;DR**: 

I need a simple command line tool can
- facilitates writing plain CloudFormation.
- have similar fashion like in [Ansible](https://www.ansible.com/) to manage parameters (variables).
- easy command to manage CloudFormation lifecycles.

**Long story**: You can check out my article [From lmdo to cfctl, a journey of two worlds]().

## Features
- Detect circular dependency amongst stacks during deployments.
- Provide file ecryption facility for secrets used in variables and automatically decrypt them during deployment.
- Configuration over convention. Provide flexibility to suit different needs in folder structure.
- No vendor lock-in. You won't loose the ability to re-use your CloudFormations even you decide to switch to a different tool.
- Handling nested stacks auto-uploading.
- Fetching stack output on the fly.

## API References  [![GoDoc](https://godoc.org/github.com/liangrog/cfctl?status.svg)](https://godoc.org/github.com/liangrog/cfctl)

## Getting Started
### Installing
1. Download the desired version base on your OS from the [releases page](https://github.com/liangrog/cfctl/releases)
2. Move it to the executables folder. For example for linux amd64: `chmod +x cfctl-linux-amd64 && sudo mv cfctl-linux-amd64 /usr/local/bin/cfctl`

cfctl piggy-backs your existing [AWSCLI](https://aws.amazon.com/cli/) credential setting. If you don't have one, there are a few options:
- Use awscli environment [variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)

or

- Create two files: `~/.aws/credentials` and `~/.aws/config` as per [instruction](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). 

### Repository Struction Explain

### Cheat Sheet
#### Validate a CloudFormation Template
```sh
 # Validate one local template
$ cfctl template validate ./template-1.yaml     

# Validate multiple local template
$ cfctl template validate ./template-1.yaml ./template-2.yaml

# Validate all templates in a folder recursively
$ cfctl template validate ./template -r

# Validate a template from internet
$ cfctl template validate https://bucket.s3.amazonaws.com/template-a.yaml

# Validate multiple templates reside in local, internet and in a folder
$ cfctl template validate ./template-1.yaml https://bucket.s3.amazonaws.com/template-a.yaml ./template -r
``` 

### Manage CloudFormation Stack Lifecyle

```sh
# Deploy all stacks without using variable.
$ cfctl stack deploy 

# Deploy all stacks from a specific stack file
$ cfctl stack deploy -f stack-file.yaml

# Deploy particular stacks using variables from specific environment
$ cfctl stack deploy --stack stack1,stack2 --env production

# Deploy stacks using variables from specific environment that contains secrets and providing password file
$ cfctl stack deploy --env production --vault-password-file path/to/password/file

# Delete a stack
$ cfctl stack delete stack-1

# Delete multiple stacks
$ cfctl stack delete stack-1 stack-2

# Delete all stacks from a specific stack file
$ cfctl stack delete -f stack-file.yaml --all

# List all stacks in an AWS account
$ cfctl stack list

# List stacks with specifc status in an AWS account
$ cfctl stack list --status DELETE_COMPLETE

# Get a specific stack
$ cfctl stack get --name stack-a 
```

### Encrypt/Decrypt Secret Variables
cfctl provides file encryption/decryption implementation as per [ansible-vault 1.1 spec](https://docs.ansible.com/ansible/latest/user_guide/vault.html#vault-payload-format-1-1). The encrypted files are interchangable with ansible-vault, in other words, the files encrypted by cfctl or ansible-vault can be decrypted by either one of them.

The command group is `cfctl vault`

The password lookup order is defined as below:
1. CLI option `--password`
2. CLI option `--password-file`
3. Environment variable `CFCTL_VAULT_PASSWORD`
4. Environment variable `CFCTL_VAULT_PASSWORD_FILE`
5. Default password file `$HOME/.cfctl_vault_password`
6. Shell prompt


Only **one** password can be used during encryption.

For decryption, multiple passwords can be seperated by using **comma delimiter (,)**. For example:
```
    password1,password2,password3...
```

All passwords will be tried until one that works. 

Here are some simple examples how to use the command:
```
    # To encrypt
    $ cfctl vault encrypt file1 file2 file3 --password secret

    # To decrypt
    $ cfctl vault decrypt file1 file2 file3 --password secret
```

### Upload Files to S3 Bucket

## Requirements
[Desgin princples and requirements](docs/requirements.md)

