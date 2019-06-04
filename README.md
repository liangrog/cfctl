# `cfctl` - a Simple AWS CloudFormation DevOps Tool
cfctl is a command line utility tool that helps to organise and manage AWS stacks that created by using CloudFormation. 

## Reason of Creation
**TL;DR**: 
I just need a simple command line tool can
- facilitates writing plain CloudFormation.
- have similar mechanism like in [Ansible](https://www.ansible.com/) to manage parameters.
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

### Best Practice

### Cheat Sheet




## Requirements
[Desgin princples and requirements](docs/requirements.md)


Trello Board
---
[github-cfctl](https://trello.com/b/3etT9edo/github-cfctl)

## Ansible Vault
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
