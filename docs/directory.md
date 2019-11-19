# Repository Structure
Because the three directories: `templateDir`, `envDir` and `paramDir` are configurable in [stack file](config.md) and all files in the sub directories of those three main directories are being referenced using relative path, users have all the freedom to decide how they would like to name, or organise directories, files.

A simple example:
```
simple/
├── parameters
├── stacks.yaml
├── templates
└── environments
```

A more complex example:
```
complex/
├── team-a
│   ├── parameters
│   │   ├── db
│   │   └── web-server
│   ├── stacks.yaml
│   └── environments
│       ├── default
│       └── prod
├── team-b
│   ├── parameters
│   │   ├── db
│   │   └── web-server
│   ├── stacks.yaml
│   └── environments
│       ├── dev
│       └── prod
└── templates
    ├── app
    ├── eks
    └── rds
```

There is only one exception that `cfctl` uses convention: if a `default` folder exists in `envDir`, all variables in this folder will be loaded first before overwritten by other variable files on every deployment.

