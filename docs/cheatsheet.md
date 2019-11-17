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

# Deploy stacks with specify tag values
$ cfctl stack deploy --stack stack1,stack2 --tags Type=frontend

# Output parameters only for all stacks
$ cfctl stack deploy --env production --param-only

#Keeping stack when creation fails and in ROLLBACK_COMPLETE state
$ cfctl stack deploy --keep-stack-on-failure

# Delete a stack
$ cfctl stack delete stack-1

# Delete multiple stacks
$ cfctl stack delete stack-1 stack-2

# Delete all stacks from a specific stack file
$ cfctl stack delete -f stack-file.yaml --all

# Delete stacks that have specific tag values
$ cfctl stack delete --tags Name=stack-1,Type=frontend

# List all stacks in an AWS account
$ cfctl stack list

# List stacks with specifc status in an AWS account
$ cfctl stack list --status DELETE_COMPLETE

# Get all stack details in config file backend-infra.yaml
$ cfctl stack get -f backend-infra.yaml

# Get a specific stack 'stack-a' detail
$ cfctl stack get --name stack-a 

# Get multiple stack details
$ cfctl stack get --name stack-a,stack-b

# Get stack details with tag Name=frontend
$ cfctl stack get --tags Name=frontend

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
```
    # To encrypt
    $ cfctl vault encrypt file1 file2 file3 --password secret

    # To decrypt
    $ cfctl vault decrypt file1 file2 file3 --password secret
```

