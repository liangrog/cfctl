package sample

// Sample stack config file
var StackYaml = `
---
# S3 Bucket for uploading templates.
# If it doesn't exist, it will create one.
s3Bucket: cfctl-templates-{{ awsAccountId }}

# Cloudformationn template folder
templateDir: ../../templates

# Values environment folder. It contains
# specific values used in the parameter templates
# for different environment.
envDir: environments

# Cloudformationn Parameters folder. It contains
# all the parameter templates.
paramDir: parameters

# Cloudformation stack list
stacks:
  - name: my-first-cfctl-stack  # Name of your stack
    tpl: s3-encrypted.yaml # Template relative path
    param: s3.yaml # Parameter file relative path
    tags: # AWS stack tags.
      Name: my-first-cfctl-stack
      Group: test`

var S3Template = `
---
AWSTemplateFormatVersion: "2010-09-09"
Description: Template to create a s3 bucket with encryption enabled

Parameters:
  BucketName:
    Description: Name of your bucket
    Type: String

  Version:
    Description: Enables multiple variants of all objects in this bucket.
    Type: String
    Default: Suspended
    AllowedValues:
      - Suspended
      - Enabled

Resources:
  Bucket:
    Type: AWS::S3::Bucket
    Properties:
      AccessControl: BucketOwnerFullControl
      BucketName: !Ref BucketName
      VersioningConfiguration:
        Status: !Ref Version
      BucketEncryption:
        ServerSideEncryptionConfiguration:
        - ServerSideEncryptionByDefault:
            SSEAlgorithm: AES256

Outputs:
  BucketName:
    Description: Name of the bucket
    Value: !Ref Bucket

  DomainName:
    Description: The IPv4 DNS name of the specified bucket
    Value: !GetAtt Bucket.DomainName
`

var SampleParam = `
---
BucketName: "{{ .s3BucketPrefix }}-{{ awsAccountId }}"
`

var EnvVars = `
---
s3BucketPrefix: "cfctl-sample"
`
