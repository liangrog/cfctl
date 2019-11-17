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


