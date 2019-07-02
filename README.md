# `cfctl` - a Simple AWS CloudFormation DevOps Tool
cfctl is a streamline command line utility that helps to organise and manage AWS stacks that created by using CloudFormation. 

## Reason of Creation
**TL;DR**: 

I need a simple command line tool can
- facilitates writing plain CloudFormation.
- have similar fashion like in [Ansible](https://www.ansible.com/) to manage parameters (variables) and deployment environment.
- easy command to manage CloudFormation lifecycles.
- support yaml parameter file.

<!-- **Long story**: You can check out my article [From lmdo to cfctl, a journey of two worlds](). -->

## Features
- Detect circular dependency amongst stacks during deployments.
- Provide file ecryption facility for secrets used in variables and automatically decrypt them during deployment.
- Configuration over convention. Provide flexibility to suit different needs in folder structure.
- No vendor lock-in. You won't loose the ability to re-use your CloudFormations even you decide to switch to a different tool.
- Handling nested stacks auto-uploading.
- Fetching stack output on the fly.

## API References  [![API References](https://godoc.org/github.com/liangrog/cfctl?status.svg)](https://godoc.org/github.com/liangrog/cfctl)

## Getting Started
### Installing
1. Download the desired version base on your OS from the [releases page](https://github.com/liangrog/cfctl/releases)
2. Move it to the executables folder. For example for linux amd64: `chmod +x cfctl-linux-amd64 && sudo mv cfctl-linux-amd64 /usr/local/bin/cfctl`

cfctl piggy-backs your existing [AWSCLI](https://aws.amazon.com/cli/) credential setting. If you don't have one, there are a few options:
- Use awscli environment [variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)

or

- Create two files: `~/.aws/credentials` and `~/.aws/config` as per [instruction](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). 

### Repository Structure
The repository structure is very flexible. It's up to users' preference how they want to structure their templates, parameters and variables as long as the required values are provided in the stack file (see [StackFile Anatomy](#stack-file-anatomy)).

A simple example:
```
simple/
├── params
├── stacks.yaml
├── templates
└── vars
```

A more complex example:
```
complex/
├── team-a
│   ├── param
│   │   ├── db
│   │   └── web-server
│   ├── stacks.yaml
│   └── vars
│       ├── dev
│       └── prod
├── team-b
│   ├── param
│   │   ├── db
│   │   └── web-server
│   ├── stacks.yaml
│   └── vars
│       ├── dev
│       └── prod
└── templates
    ├── app
    ├── eks
    └── rds
```

### Varaibles
`cfctl` provide a simple templating system that allows
1. Use different values for the same parameter based on deployment environment. 
2. Source values from environment variables and other stack outputs.
3. Automatically upload nested templates and pass the s3 url to the value.

Using above complex example, supposed you have a parameter file in `team-a/param/db/web.yaml`:
```yaml
DBName: '{{ .WebDBName }}'
```
Then you could provide different database names for different environments:

`team-a/vars/dev/vars.yaml`
```yaml
WebDBName: dev-db
```
`team-a/vars/prod/vars.yaml`
```yaml
WebDBName: prod-db
```

When deploying stacks, using `--env prod` for example, the `team-a/vars/prod/vars.yaml` will be used.

**Tips:** 
1. You can set a globle value to a var by creating a folder with name `default` under `vars` folder for example. Every variables in the `default` folder will be loaded first.
2. The var files` name can be anything. However if there are multiple var files in the same folder, the var files will be loaded in lexical order. The later one will override the previous one. For example, override the vars in `default` folder.

Variable files can be encrypted using `cfctl vault encrypt` command. The encrypted files will be automatically decrypted during deployment. Details please see [Encrypt or Decrypt Secret Variable Files](#encrypt-or-decrypt-secret-variable-files).

### Variable Functions
#### Get values from another stack's outputs
`'{{ stackOutput "stack-name" "value name in the outputs"}}'`

#### Get values from environment variables
`'{{ env "variable name" }}'`

#### Auto-uploading nested template and pass url to value
`'{{ tpl "rds/mysql.yaml" }}'`

**NOTE:** Please strictly following the example for the quotes and space when using them. Otherwise you might encounter errors.

**Tips:** 
You could combine multiple functions. For example:
`'{{ env "variable name" }}, {{ stackOutput "stack-name" "value name in the outputs"}}'`

You could also nested functions. For example, you might have `{{ env "DB_TEMPLATE" }}`, then you could use `export DB_TEMPLATE={{ tpl "rds/mysql.yaml" }}`. 

### Stack File Anatomy
The default stack file name is stacks.yaml. You can use custom names as long as your provide it to `-f` in command.
```yaml
# Required: true
#
# AWS S3 bucket name, where the nested stack templates will be uploaded into.
# If the bucket doesn't exist, cfctl will create it for you as long as the IAM
# has the correct permission.
s3Bucket: my-bucket

# Required: true
# 
# The relative (to stack file) path of the directory where all your Cloudformation template files reside.
templateDir: relative/path/to/template/folder

# Required: true
#
# The relative (to stack file) path of the directory where all your templates` parameter files reside.
paramDir: relative/path/to/parameter/folder

# Required: true
#
# The relative (to stack file) path of the directory where all your deployment specific variables are.
envDir: relative/path/to/deployment/vars/folder

# Required: true
#
# The stack list
stacks:
  - name: stack-a           # Stack name. 
    tpl: web-server.yaml    # Stack template file. Relative path to "templateDir": [templateDir]/web-server.yaml.
    param: web/server.yaml  # Template parameter file. Relative path to "paramDir": [paramDir]/web/server.yaml.
    tags:                   # Tags for the stack.
      component: web
  - name: stack-b           # Stack name.
    tpl: rds/mysql.yaml     # Stack template file. Relative path to "templateDir": [templateDir]/rds/mysql.yaml.
    param: web/db.yaml      # Template parameter file. Relative path to "paramDir": [paramDir]/web/db.yaml.
    tags:                   # Tags for the stack.
      component: web
```

**NOTE:** Variable function `env` is available for stack file.

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
Creating and updating a stack shares the same command `cfctl stack deploy`.
```sh
# Deploy all stacks without using variable.
$ cfctl stack deploy 

# Deploy all stacks from a specific stack file
$ cfctl stack deploy -f stack-file.yaml

# Deploy particular stacks using variables from specific environment
$ cfctl stack deploy --stack stack1,stack2 --env production

# Deploy stacks using variables from specific environment that contains secrets and providing password file
$ cfctl stack deploy --env production --vault-password-file path/to/password/file

# Override environment values
$ cfctl stack deploy --env production --vault-password-file path/to/password/file --vars name1=value1,name2=value2

# Output parameters only for all stacks
$ cfctl stack deploy --env production --param-only

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

### Upload Files to S3 Bucket
```sh
# Upload one file
$ cfctl s3 upload file-1 --bucket my-bucket

# Upload multiple files
$ cfctl s3 upload file-1 file-2 --bucket my-bucket

# Upload everything in a folder recursively
$ cfctl s3 upload template/web --bucket my-bucket -r

# Upload everything in a folder recursively except fileA
$ cfctl s3 upload template/web --bucket my-bucket -r --exclude-files fileA 

# Upload files and folder
$ cfctl s3 upload file-1 template/web --bucket my-bucket -r
```

### Encrypt or Decrypt Secret Variable Files
cfctl provides file encryption/decryption implementation as per [ansible-vault 1.1 spec](https://docs.ansible.com/ansible/latest/user_guide/vault.html#vault-payload-format-1-1). The encrypted files are interchangable with ansible-vault, in other words, the files encrypted by cfctl or ansible-vault can be decrypted by either one of them.

The password lookup order is defined as below:
1. CLI option `--vault-password`
2. CLI option `--vault-password-file`
3. Environment variable `CFCTL\_VAULT\_PASSWORD`
4. Environment variable `CFCTL\_VAULT\_PASSWORD\_FILE`
5. Default password file `$HOME/.cfctl\_vault\_password`
6. Shell prompt

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

