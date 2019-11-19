### Encrypt or Decrypt Secret Variable Files
cfctl provides file encryption/decryption implementation as per [ansible-vault 1.1 spec](https://docs.ansible.com/ansible/latest/user_guide/vault.html#vault-payload-format-1-1). The encrypted files are interchangable with ansible-vault, in other words, the files encrypted by cfctl or ansible-vault can be decrypted by either one of them.

The password lookup order is defined as below:
1. CLI option `--vault-password`
2. CLI option `--vault-password-file`
3. Environment variable `CFCTL_VAULT_PASSWORD`
4. Environment variable `CFCTL_VAULT_PASSWORD_FILE`
5. Default password file `$HOME/.cfctl_vault_password`
6. Shell prompt

For decryption, multiple passwords can be seperated by using **comma delimiter (,)**. For example:
```
    password1,password2,password3...
```

All passwords will be tried until one that works. 

Here are some simple examples how to use the command:
```
# To encrypt
$ cfctl vault encrypt file1 

# To decrypt
$ cfctl vault decrypt file1 file2 file3 --password secret
```

