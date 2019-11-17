### Varaibles
`cfctl` provide a simple templating system that allows
1. Use different values for the same parameter based on deployment environment. 
2. Source values from environment variables and other stack outputs.
3. Automatically upload nested templates and pass the s3 url to the value.

Notes: the yaml single line has a limit of 80 chars. If longer than that limit, please use <b>`>`</b> or <b>`|`</b>. The common error you will see if you don't use multi-line:
```
Error: template: 78723a9a-8820-483b-b451-753d0fb8c229:9: unclosed action
```

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


