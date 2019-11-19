# Cheat Sheet
## CloudFormation Template Validation
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

## Manage CloudFormation Stack Lifecyle

## Stack Creation and Updates
```sh
# Deploy all stacks in the default stacks.yaml file with default.
$ cfctl stack deploy 

# Deploy all stacks from a specific stack file
$ cfctl stack deploy -f stack-file.yaml

# Deploy particular stacks using variables from specific environment
$ cfctl stack deploy --stack stack1,stack2 --env production

# Deploy stacks using variables from production environment that contains secrets and providing password file
$ cfctl stack deploy --env production --vault-password-file path/to/password/file

# Override environment values
$ cfctl stack deploy --env production --vault-password-file path/to/password/file --vars name1=value1,name2=value2

# Deploy stacks have specify tag values
$ cfctl stack deploy --stack stack1,stack2 --tags Type=frontend

# Output parameters only for all stacks
$ cfctl stack deploy --env production --param-only

# Keeping stack when creation fails and in ROLLBACK_COMPLETE state, otherwise the stack will be deleted.
$ cfctl stack deploy --keep-stack-on-failure
```

## Stack Deletion
```sh
# Delete a stack
$ cfctl stack delete stack-1

# Delete multiple stacks
$ cfctl stack delete stack-1 stack-2

# Delete all stacks from a specific stack file
$ cfctl stack delete -f stack-file.yaml --all

# Delete stacks have specific tag values
$ cfctl stack delete --tags Name=stack-1,Type=frontend
```

## Stack Queries
```sh
# List all stacks in an AWS account
$ cfctl stack list

# List stacks with specifc status in an AWS account
$ cfctl stack list --status DELETE_COMPLETE

# Get all stack details in from stack file  backend-infra.yaml
$ cfctl stack get -f backend-infra.yaml

# Get a specific stack 'stack-a' detail
$ cfctl stack get --name stack-a 

# Get multiple stack details
$ cfctl stack get --name stack-a,stack-b

# Get stack details with tag Name=frontend
$ cfctl stack get --tags Name=frontend

# Get all stack resources in config file backend.yaml
$ cfctl stack get-resources --file backend.yaml

# Get a specific stack 'stack-a' resources
$ cfctl stack get-resources --name stack-a

# Get multiple stacks' resources
$ cfctl stack get-resources --name stack-a,stack-b

# Get stack resources details with tag Name=frontend
$ cfctl stack get-resources --tags Name=frontend
```

## S3 Upload
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

## Encrypt or Decrypt Secret Files
```sh
# Encrypt multiple files
$ cfctl vault encrypt file1 file2 file3

# Encrypt using environment value
$ export CFCTL_VAULT_PASSWORD=xxxx
$ cfctl vault encrypt filename 

# Encrypt using password file
$ cfctl vault encrypt filename --vault-password-file path/to/password/file

# Edit encrypted vault file
$ cfctl vault edit filename

# Decrypt multiple files
$ cfctl vault decrypt file1 file2 file3
```

