# Go template

> Go templates for new projects

This repository is a template to make use of when creating new projects in
Go. It contains scripts, Dockerfile(s) and workflows.

Four templates are provided:

* [base](templates/base/) - Base project.
* [http-server](templates/http-server/) - Basic HTTP server project.
* [server](templates/server/) - General server project.
* [service](templates/service/) - General service project.

## Getting started

To create a new project based on one of the templates:

1. Install `gonew` (if it is not already installed):
```sh
go install golang.org/x/tools/cmd/gonew@latest
```
2. Download the template and create the new project:
```sh
# Assuming the project will be hosted in GitHub.
# If not replace github.com/<owner>/<repo> with the correct path.
gonew github.com/RedeployAB/go-template/templates/<name> github.com/<owner>/<repo>
```

## Example

To create a new project based on the [http-server](templates/http-server/) template
to a GitHub repository with the user/organization name `YourUser` and the project/module name `myproject`

```sh
gonew github.com/RedeployAB/go-template/templates/http-server github.com/YourUser/myproject
```

**Note**: The directory and files will be created in the current working directory.
Given the example path `development/go`, running the `gonew` command will create
the new `myproject` into `development/go/myproject`.

This will download the module, it's assets and rewrite the module name and import paths
to match the new module name (project).
