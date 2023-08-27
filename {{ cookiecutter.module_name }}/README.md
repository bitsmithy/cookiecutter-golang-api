# {{ cookiecutter.module_name }}

## Development

This project assumes you are using [asdf](https://asdf-vm.com/) to manage
development tools. Once installed, simply run `make install` to perform a
one-time setup installation of required tools.

The versions of Go and all required tools are specified in the
`.tool-versions` file of this project.

### Makefile

This project makes use of a Makefile to run all development and operations
related tasks. Simply run `make` to see the menu of targets that are available.

### Environment Variables

Any project-specific environment variables is handled automatically by
[direnv](https://direnv.net) and the `.envrc` file in this project.
