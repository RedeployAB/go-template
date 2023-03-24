# Go template

> Go template for new projects

This repository is a template to make use of when creating new projects in
Go. It contains scripts, Dockerfile(s) and workflows.

* [Scripts](#scripts)
* [Dockerfiles](#dockerfiles)
* [Workflows](#workflows)


## Scripts

### `build.sh`

Script to build the binary and optionally build a docker image
containing said binary.

**Note**: The script should have the variable `bin=` at the top modified to match the name of the projects binary, most often the directory name.

**Usage**

```sh
# To build binary.
./scripts/bash/build.sh --version <version>

# To build binary and docker image.
./scripts/bash/build.sh --version <version> --image
```

### `release.sh`

Script to prepare a release. The script makes sure the current branch is the base branch (often `main`), pulls from remote, run tests and then finally creates a Git tag with the version number.

**Usage**

```sh
# Set a version (valid semver without the 'v' prefix).
./scripts/bash/release.sh --version <version>

# If there already is existing tagged versions, the following can be used.

# Patch increment (patch number according to semantic versioning).
./scripts/bash/release.sh --patch

# Minor increment (minor number according to semantic versioning).
./scripts/bash/release.sh --minor

# Major increment (minor number according to semantic versioning).
./scripts/bash/release.sh --major
```

### `go-install.sh`

Script to install dependencies of a Go project when having dependencies in a private repository. When that is not the case, there is no need to make use of this script.

**Usage**

```sh
export PRIVATE_REPO_URL=<url-to-private-repository>
export PRIVATE_REPO_SSH_KEY_BASE64=<base64-of-private-key-to-repository>

./scripts/bash/go-install.sh
```

## Dockerfiles

Two Dockerfiles are provided:

* `Dockerfile` - Needs the binary pre-built and located in `build/`.
* `Dockerfile_build` - Builds the binary during the image build.

**Note**: The Dockerfile(s) needs the following updated:

**First step**

* `ARG BIN` needs to be updated `ARG BIN=<binary-name>` (if not provided during build).

**Second step**

* `ARG BIN` needs to be updated `ARG BIN=<binary-name>` (if not provided during build).
* `ARG PORT` needs to be updated to `ARG PORT=<port-number>` (if not provided during build).

* `ENTRYPOINT` needs to be updated to `ENTRYPOINT [ "/<binary-name>" ]`.

If `ca-certificates` is not needed by the project, the following lines can be deleted:

* `RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates` from the **First step**.
* `COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/` from the **Second step**.


## Workflows

### `test.yaml`

Workflow to be called by other workflows, runs a job with the most common steps for testing a Go module.

### `build.yaml`

Workflow that includes a call to the `test.yaml` workflow. A starter workflow that needs to be modified to suit the project.
