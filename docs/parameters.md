# Templating for Parameter Files
`cfctl` provide a simple templating system by utilising [go template](https://golang.org/pkg/text/template/) for the parameters file. The parameters file must be written in YAML with "key: value"  format. 

For example:
```
---
ClusterControlPlaneSecurityGroup: ""
NodeGroupName: utility
NodeAutoScalingGroupMinSize: "{{ .EKSNodeGroupMinSize }}"
NodeAutoScalingGroupDesiredCapacity: "{{ .EKSNodeGroupDesiredCapacity }}"
NodeAutoScalingGroupMaxSize: "{{ .EKSNodeGroupMaxSize }}"
NodeInstanceType: "{{ .EKSInstanceType }}"
NodeVolumeSize: "{{ .EKSNodeVolumeSize }}"
KeyName: "{{ .EKSKeyName }}"
BootstrapArguments: "--kubelet-extra-args,node-role.kubernetes.io/utility=true"
NodeExtraSecurityGroupIds: >
  '{{ .SgInternet }}, {{ stackOutput "allow-ssh" "SecurityGroupId"}}, {{ stackOutput "eks-sg-nodes" "SecurityGroupId"}}'

VpcId: "{{ .VpcId }}"
Subnets: "{{ .PrivateSubnet2a }}"
```

The values of those variables will be stored in variable files inside the directory defined in `envDir` folder.


You can set a default globle value to a variable by creating a folder with name `default` under the fodler defined in `envDir` folder. Every variables in the `default` folder will be loaded first.


Variable files can be encrypted using `cfctl vault encrypt` command. The encrypted files will be automatically decrypted during deployment.


## Important
1. When using variables and functions, the string must be quoted.
2. The yaml single line has a limit of 80 chars. If longer than that limit, please use <b>`>`</b> or <b>`|`</b>. The common error you will see if you don't use multi-line: `Error: template: 78723a9a-8820-483b-b451-753d0fb8c229:9: unclosed action`.
3. Functions can be chained using `|`. 
4. The variable file name can be anything. However if there are multiple variable files in the same folder, the files will be loaded in lexical order. The later one will override the previous one.


## Functions
Apart from standard template functions provided by [go template](https://golang.org/pkg/text/template/), it also provides some useful addtional functions:

- <b>env:</b> "{{ env ENV_VARIABLE_NAME}}" will parsing environment variabe.
- <b>awsAccountId:</b> "{{ awsAccountId }}" will get your AWS account ID for your current IAM user.
- <b>hash:</b> "{{ printf "%s" "test" | hash }}" will hash the string "test" using md5
- <b>stackOutput:</b> '{{ stackOutput "stack-name" "value name in the outputs"}}'` will get the value of the output. Note: There can not be a space between value name and the last double curly bracket.
- <b>tpl:</b> '{{ tpl "rds/mysql.yaml" }}' will upload the template to S3 bucket then returns the url.

