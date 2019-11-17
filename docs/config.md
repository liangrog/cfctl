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

**NOTE:** Variable function `env`, `awsAccountId` and `hash` is available for stack file.


