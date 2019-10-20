package cmd

const (
	// Root
	CMD_ROOT_OUTPUT = "output"

	// Valut

	// Environment variable name for vault password.
	ENV_VAULT_PASSWORD = "CFCTL_VAULT_PASSWORD"

	// Environment variable name for vault password file.
	ENV_VAULT_PASSWORD_FILE = "CFCTL_VAULT_PASSWORD_FILE"

	// Command line flag for vault password.
	CMD_VAULT_PASSWORD = "vault-password"

	// Command line flag for vault password file.
	CMD_VAULT_PASSWORD_FILE = "vault-password-file"

	// Default vault password file name.
	DEFAULT_VAULT_PASSWORD_FILE = ".cfctl_vault_password"

	// S3

	// Command line flag for bucket.
	CMD_S3_UPLOAD_BUCKET = "bucket"

	// Command line flag for bucket prefix.
	CMD_S3_UPLOAD_PREFIX = "prefix"

	// Command line flag for recursive.
	CMD_S3_UPLOAD_RECURSIVE = "recursive"

	// Command line flag for excluding files.
	CMD_S3_UPLOAD_EXCLUDE_FILES = "exclude-files"

	// Stack

	// Command line flag for stacks.
	CMD_STACK_DEPLOY_STACK = "stack"

	// Command line flag for configuration file.
	CMD_STACK_DEPLOY_FILE = "file"

	// Command line flag for dry run.
	CMD_STACK_DEPLOY_DRY_RUN = "dry-run"

	// Command line flag for envoirnment folder.
	CMD_STACK_DEPLOY_ENV = "env"

	// Command line flag for tag grouping.
	CMD_STACK_DEPLOY_TAGS = "tags"

	// Parameter parsing.
	CMD_STACK_DEPLOY_PARAM_ONLY = "param-only"

	// Variable override
	CMD_STACK_DEPLOY_VARS = "vars"

	// Default environment folder name.
	STACK_DEPLOY_ENV_DEFAULT_FOLDER = "default"

	// Command line flag for stack delete all.
	CMD_STACK_DELETE_ALL = "all"

	// Command line flag for stack get name.
	CMD_STACK_GET_NAME = "name"

	// Command line flag for stack list status.
	CMD_STACK_LIST_STATUS = "status"

	// Template

	// Command line flag for template validate recursively.
	CMD_TEMPLATE_VALIDATE_RECURSIVE = "recursive"
)
