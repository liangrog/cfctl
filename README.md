# cfctl
AWS CloudFormation DevOp tool

## Requirements
[Desgin princples and requirements](docs/requirements.md)


Trello Board
---
[github-cfctl](https://trello.com/b/3etT9edo/github-cfctl)

## Ansible Vault
cfctl provides file encryption/decryption implementation as per [ansible-vault 1.1 spec](https://docs.ansible.com/ansible/latest/user_guide/vault.html#vault-payload-format-1-1). The encrypted files are interchangable with ansible-vault.

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
