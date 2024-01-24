# {{ cookiecutter.module_name }}

## Development

This project assumes you are using [mise](https://mise.jdx.dev) to manage
development tools. Once installed, simply run `make setup` to perform a
one-time setup installation of required development dependencies.

The versions of Go and all required tools are specified in the
`.tool-versions` file of this project.

### Makefile

This project makes use of a Makefile to run all development and operations
related tasks. Simply run `make` or `make help` to see the menu of targets that are available.

### Environment Variables

Any project-specific environment variables is handled automatically by
[direnv](https://direnv.net) and the `.envrc` file in this project.
