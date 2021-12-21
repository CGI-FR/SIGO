# Go Template

This repository is a Go project template. This is not a real application, the go source code is an empty shell.

List of folders

- `cmd` contains source code of all compiled binaries, each sub-folder correspond to a binary.
- `internal` contains source code that cannot be linked in an external project.
- `pkg` contains source code that can be linked in an external project.
- `test` contains integration tests source code (run with venom).
- `githooks` contains git hooks for better automation.

## Usage

### Prerequisites

You need :

- Visual Studio Code ([download](https://code.visualstudio.com/)) with the [Remote - Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed.
- Docker Desktop (Windows, macOS) or Docker CE/EE (Linux)

Details are available on [the official Visual Studio documentation](https://code.visualstudio.com/docs/remote/containers#_getting-started).

### Initialize a new repository using Github

Use the button `Use this template` at the top of this page. It will automatically initialize your new repository with this template.

### Initialize a new repository without using Github

```console
$ cd </your/project/root>
/your/project/root$ wget -nv -O- https://github.com/adrienaury/go-template/archive/refs/heads/main.tar.gz | tar --strip-components=1 -xz
/your/project/root$ git init -b main
/your/project/root$ git add .
/your/project/root$ git commit -m "chore: init repository from go template"
```

### Modify an existing repository

Warning: do this in a branch where to isolate the changes

```console
$ cd </your/project/root>
/your/project/root$ wget -nv -O- https://github.com/adrienaury/go-template/archive/refs/heads/main.tar.gz | tar --strip-components=1 -xz
```

### Things you might or should delete (or replace)

- might: the `.vscode` folder contains a VSCode color theme, you might want to remove this folder, or customize the colors
- might: the `githooks/commit-msg` file contains a commit message check to enforce [semantic commit message](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)
- might: the `Makefile` if you want to only use the neon build tool, but I think it's nice to give a well known handle to the build workflow
- might: the `.github` folder if you don't want to use GitHub actions
- should: the `Dockerfile` and the `Dockerfile.webserver` if your project does not produce Docker images (DO NOT DELETE the `docker-compose.yml` file, you can move it under `.devcontainer` and adapt paths)
- should: the `LICENSE` file as your project is not copyrighted by me ! Replace it by your own license
- should: everything under `test/suites`, replace it with your own tests suites or remove the `test` folder completely if you don't do integration tests
- should: every `.go` file will have to be deleted, as well as folders under `cmd`, `pkg` and `internal`
- should: this `README.md` file, of course ;)

### Run your workspace

When opening the folder with [Visual Studio Code](https://code.visualstudio.com/), the [Remote - Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) will detect the devcontainer configuration and ask you to reopen the project in a container.

Accept and enjoy !

## Features

- Structured Go code following folder names convention
- VSCode devcontainer pre-parameterized
- Auto code formatting, using EditorConfig extension
- Changelog, Contrib : all initialized with default templates
- Git commit message semantic validation
- Docker compatible (docker client and docker-compose are available inside the devcontainer)
- Pre-configured GitHub actions : CI on pull requests, Release on tag
- Build targets (run with `make` or `neon`) :
  - [`help`](#help-target) : default target, print help message
  - [`info`](#info-target) : print information on the build pipeline
  - [`promote`](#promote-target) : promote the project to a new tag using semantic versioning
  - [`refresh`](#refresh-target) : refresh go modules dependencies
  - [`compile`](#compile-target) : compile sources
  - [`lint`](#lint-target) : check the code for suspicious constructs
  - [`test`](#test-target) : run the unit tests
  - [`release`](#release-target) : compile binaries, with production flags
  - [`test-int`](#test-int-target) : run integration tests with venom
  - [`publish`](#publish-target) : publish binaries on Github with goreleaser
  - [`docker`](#docker-target) : build docker images
  - [`docker-tag`](#docker-tag-target) : tag docker images using semantic versioning
  - [`docker-push`](#docker-push-target) : publish docker images on Dockerhub
  - [`license`](#license-target) : scan binaries for 3rd party licenses and generate a notice file

### Build targets

Run a build target by using the neon command.

```console
neon target
```

This text bloc show how target are related to each other. E.g. running the target `lint` will also run `info` and `refresh`.

```text
→ help
→ promote
→ info ┰ docker → docker-tag → docker-push
       ┖ refresh ┰ compile → license
                 ┖ lint → test → release → test-int → publish
```

Multiple targets can be run in the same command, e.g. `neon release docker-tag`.

Neon targets are also mapped to a Makefile, so running `make compile` will produce the same result as running `neon compile`.

#### Help target

```console
$ neon help
----------------------------------------------- help --
Available targets

help         Print this message
info         Print build informations
promote      Promote the project with a new tag based on git log history
refresh      Refresh go modules (add missing and remove unused modules) [info]
compile      Compile binary files locally [info->refresh]
lint         Examine source code and report suspicious constructs [info->refresh]
test         Run all tests with coverage [info->refresh->lint]
release      Compile binary files for production [info->refresh->lint->test]
test-int     Run all integration tests [info->refresh->lint->test->release]
publish      Publish tagged binary to Github [info->refresh->lint->test->release->test-int]
docker       Build docker images [info]
docker-tag   Tag docker images [info->docker]
docker-push  Publish docker images to Dockerhub [info->docker->docker-tag]
license      Scan licenses from binaries and generate notice file [info->refresh->compile]

Example: neon -props '{latest: true}' promote publish

Target dependencies

→ help
→ promote
→ info ┰ docker → docker-tag → docker-push
       ┖ refresh ┰ compile → license
                 ┖ lint → test → release → test-int → publish
OK
```

#### Info target

Print build informations, like the author or the current tag.

```console
$ neon info
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = 00424c8c67bca5b11ed99efa0d45902f1143cbd7
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
OK
```

#### Promote target

Promote the project with a new tag based on git log history, or based on the parameter passed with `-props` flag.

Without parameter, the tag name will be determined by the git commit history since the last tag (or will be equal to `v0.1.0` if there is no existing tag). This is base on the [`svu`](https://github.com/caarlos0/svu) tool.

```console
$ neon promote
--------------------------------------------- promote --
Promoted to v0.2.0
OK
```

It's possible to use the `-props` flag to override the name of the tag.

```console
$ neon -props '{tag: "v0.2.1-alpha"}' promote
--------------------------------------------- promote --
Promoted to v0.2.1-alpha
OK
```

#### Refresh target

Refresh go modules (add missing and remove unused modules).

This target will keep your `go.mod` and `go.sum` files clean.

```console
$ neon refresh
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = 00424c8c67bca5b11ed99efa0d45902f1143cbd7
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
OK
```

#### Compile target

Compile binary files locally.

By default, the `cmd` folder is scanned and each subfolder will create a binary with the name of the subfolder.

```console
$ neon compile
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = b0b418aa05db7f386275249ea641f14b295cf3ab
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
--------------------------------------------- compile --
Building cmd/cli
Building cmd/webserver
OK
```

It's possible to use the `-props` flag to specify a list of folders to compile. Be aware that if one of these folders does not have a `main` package, the result file will not be executable.

```console
$ neon -props '{buildpaths: ["internal/helloservice"]}' compile
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = b0b418aa05db7f386275249ea641f14b295cf3ab
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
--------------------------------------------- compile --
Building internal/helloservice
OK
```

#### Lint target

Examine source code and report suspicious constructs. Under the hood, the [`golangci-lint`](https://github.com/golangci/golangci-lint) tool is used.

```console
$ neon lint
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = b0b418aa05db7f386275249ea641f14b295cf3ab
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
------------------------------------------------ lint --
Running command: golangci-lint run --fast --enable-all --disable scopelint --disable forbidigo
OK
```

By default, all fast linters are enabled (`--fast` and `--enable-all` flags) but you can change this with the following build properties :

- `linters` : an array of linters to enable, if left empty then all fast linters are enabled.
- `lintersno` : an array of linters to disable, by default `scopelint` (deprecated) and `forbidigo` are disabled.

To change the default values, edit the `build.yml` file and look for the properties names `linters` or `lintersno`.

These build properties can also be set by the `neon -props` flag.

```console
$ neon -props '{linters: ["deadcode"]}' lint
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = b0b418aa05db7f386275249ea641f14b295cf3ab
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
------------------------------------------------ lint --
Running command: golangci-lint run --enable deadcode --disable scopelint --disable forbidigo
OK
```

#### Test target

Run all tests with coverage.

```console
$ neon test
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = 6ac8d1ab1facb9969a84b330d08e0f3efac55819
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
------------------------------------------------ lint --
Running command: golangci-lint run --fast --enable-all --disable scopelint --disable forbidigo
------------------------------------------------ test --
?       github.com/adrienaury/go-template/cmd/cli       [no test files]
?       github.com/adrienaury/go-template/cmd/webserver [no test files]
?       github.com/adrienaury/go-template/internal/helloservice [no test files]
?       github.com/adrienaury/go-template/pkg/nameservice       [no test files]
OK
```

#### Release target

Compile binary files for production.

The only difference with the [`compile`](#compile) target is with the `ldflags` passed to the Go linker (it will produce a smaller binary) and the dependency to other targets ([`lint`](#lint) and [`test`](#test)).

```console
$ neon release
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = b0b418aa05db7f386275249ea641f14b295cf3ab
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
------------------------------------------------ lint --
Running command: golangci-lint run --fast --enable-all --disable scopelint --disable forbidigo
------------------------------------------------ test --
?       github.com/adrienaury/go-template/cmd/cli       [no test files]
?       github.com/adrienaury/go-template/cmd/webserver [no test files]
?       github.com/adrienaury/go-template/internal/helloservice [no test files]
?       github.com/adrienaury/go-template/pkg/nameservice       [no test files]
--------------------------------------------- release --
Calling target 'compile'
--------------------------------------------- compile --
Building cmd/cli
Building cmd/webserver
OK
```

The build properties are the same as the [`compile`](#compile) target.

#### Test-int target

Run all integration tests. Under the hood the tool [`venom`](https://github.com/ovh/venom) is used.

By default it will run every test suites under the folder `test/suites`.

```console
neon test-int
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = 9791b0c79b55f2f34517d7e6b64d4900c8c7f2ce
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
------------------------------------------------ lint --
Running command: golangci-lint run --fast --enable-all --disable scopelint --disable forbidigo
------------------------------------------------ test --
?       github.com/adrienaury/go-template/cmd/cli       [no test files]
?       github.com/adrienaury/go-template/cmd/webserver [no test files]
?       github.com/adrienaury/go-template/internal/helloservice [no test files]
?       github.com/adrienaury/go-template/pkg/nameservice       [no test files]
--------------------------------------------- release --
Calling target 'compile'
--------------------------------------------- compile --
Building cmd/cli
Building cmd/webserver
-------------------------------------------- test-int --
 • run cli (test/suites/01-run-cli.yml)
        • no-arguments SUCCESS
 • run webserver (test/suites/02-run-webserver.yml)
        • no-arguments SUCCESS
OK
```

#### Publish target

Publish tagged binary to Github (as a Release). Under the hood, the [`goreleaser`](https://github.com/goreleaser/goreleaser) tool is used.

Edit the file `./goreleaser.template.yml` to configure the [`goreleaser`](https://github.com/goreleaser/goreleaser) build.

A prerequisite to this target is that a file named `.github.yml` at the home directory (`~/.github.yml`) contains a `GITHUB_TOKEN` property.

```console
neon publish
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = 00a4bdbf147a4394aa1e7f0483802f94658e9ce3
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- refresh --
go: creating new go.mod: module github.com/adrienaury/go-template
go: to add module requirements and sums:
        go mod tidy
------------------------------------------------ lint --
Running command: golangci-lint run --fast --enable-all --disable scopelint --disable forbidigo
------------------------------------------------ test --
?       github.com/adrienaury/go-template/cmd/cli       [no test files]
?       github.com/adrienaury/go-template/cmd/webserver [no test files]
?       github.com/adrienaury/go-template/internal/helloservice [no test files]
?       github.com/adrienaury/go-template/pkg/nameservice       [no test files]
--------------------------------------------- release --
Calling target 'compile'
--------------------------------------------- compile --
Building cmd/cli
Building cmd/webserver
-------------------------------------------- test-int --
 • run cli (test/suites/01-run-cli.yml)
        • no-arguments SUCCESS
 • run webserver (test/suites/02-run-webserver.yml)
        • no-arguments SUCCESS
--------------------------------------------- publish --

   • releasing...
   • loading config file       file=bin/.goreleaser.yml
   • loading environment variables
   • getting and validating git state
      • releasing v0.1.0, commit 00a4bdbf147a4394aa1e7f0483802f94658e9ce3
      • pipe skipped              error=disabled during snapshot mode
   • parsing tag
   • running before hooks
      • running go mod download
   • setting defaults
      • snapshotting
      • github/gitlab/gitea releases
      • project name
      • loading go mod information
      • building binaries
      • creating source archive
      • archives
      • linux packages
      • snapcraft packages
      • calculating checksums
      • signing artifacts
      • docker images
      • artifactory
      • blobs
      • homebrew tap formula
      • scoop manifests
      • milestones
   • snapshotting
   • checking ./dist
      • --rm-dist is set, cleaning it up
   • loading go mod information
   • writing effective config file
      • writing                   config=bin/dist/config.yaml
   • generating changelog
      • pipe skipped              error=not available for snapshots
   • building binaries
      • building                  binary=/workspaces/go-template/bin/dist/cmd/cli_windows_386/cli.exe
      • building                  binary=/workspaces/go-template/bin/dist/cmd/cli_darwin_arm64/cli
      • building                  binary=/workspaces/go-template/bin/dist/cmd/cli_linux_arm64/cli
      • building                  binary=/workspaces/go-template/bin/dist/cmd/cli_linux_386/cli
      • building                  binary=/workspaces/go-template/bin/dist/cmd/cli_windows_amd64/cli.exe
      • building                  binary=/workspaces/go-template/bin/dist/cmd/cli_darwin_amd64/cli
      • building                  binary=/workspaces/go-template/bin/dist/cmd/cli_linux_amd64/cli
      • building                  binary=/workspaces/go-template/bin/dist/cmd/webserver_windows_386/webserver.exe
      • building                  binary=/workspaces/go-template/bin/dist/cmd/webserver_linux_amd64/webserver
      • building                  binary=/workspaces/go-template/bin/dist/cmd/webserver_linux_386/webserver
      • building                  binary=/workspaces/go-template/bin/dist/cmd/webserver_windows_amd64/webserver.exe
      • building                  binary=/workspaces/go-template/bin/dist/cmd/webserver_darwin_amd64/webserver
      • building                  binary=/workspaces/go-template/bin/dist/cmd/webserver_darwin_arm64/webserver
      • building                  binary=/workspaces/go-template/bin/dist/cmd/webserver_linux_arm64/webserver
   • archives
      • creating                  archive=bin/dist/go-template_v0.1.0-SNAPSHOT-00a4bdb_windows_386.tar.gz
      • creating                  archive=bin/dist/go-template_v0.1.0-SNAPSHOT-00a4bdb_linux_386.tar.gz
      • creating                  archive=bin/dist/go-template_v0.1.0-SNAPSHOT-00a4bdb_darwin_arm64.tar.gz
      • creating                  archive=bin/dist/go-template_v0.1.0-SNAPSHOT-00a4bdb_windows_amd64.tar.gz
      • creating                  archive=bin/dist/go-template_v0.1.0-SNAPSHOT-00a4bdb_linux_amd64.tar.gz
      • creating                  archive=bin/dist/go-template_v0.1.0-SNAPSHOT-00a4bdb_linux_arm64.tar.gz
      • creating                  archive=bin/dist/go-template_v0.1.0-SNAPSHOT-00a4bdb_darwin_amd64.tar.gz
   • creating source archive
   • linux packages
   • snapcraft packages
   • calculating checksums
      • checksumming              file=go-template_v0.1.0-SNAPSHOT-00a4bdb_windows_amd64.tar.gz
      • checksumming              file=go-template_v0.1.0-SNAPSHOT-00a4bdb_darwin_arm64.tar.gz
      • checksumming              file=go-template_v0.1.0-SNAPSHOT-00a4bdb_linux_amd64.tar.gz
      • checksumming              file=go-template_v0.1.0-SNAPSHOT-00a4bdb_darwin_amd64.tar.gz
      • checksumming              file=go-template_v0.1.0-SNAPSHOT-00a4bdb_linux_arm64.tar.gz
      • checksumming              file=go-template_v0.1.0-SNAPSHOT-00a4bdb_windows_386.tar.gz
      • checksumming              file=go-template_v0.1.0-SNAPSHOT-00a4bdb_linux_386.tar.gz
   • signing artifacts
   • docker images
   • publishing
      • blobs
      • http upload
      • custom publisher
      • artifactory
      • docker images
         • pipe skipped              error=publishing is disabled
      • docker manifests
         • pipe skipped              error=publishing is disabled
      • snapcraft packages
         • pipe skipped              error=publishing is disabled
      • github/gitlab/gitea releases
         • pipe skipped              error=publishing is disabled
      • homebrew tap formula
      • scoop manifests
         • pipe skipped              error=publishing is disabled
      • milestones
         • pipe skipped              error=publishing is disabled
   • release succeeded after 1.39s
OK
```

The build properties are the same as the [`compile`](#compile) target. Additionaly the `snapshot` property can be set to `true` to run the target without uploading anything to GitHub.

#### Docker target

Build docker images locally.

Dockerfiles and build contexts are configured by the `dockerfiles` map. It'a map with Dockerfiles paths as keys and context paths as values. The default value is `{"Dockerfile": ".", "Dockerfile.webserver", "."}` which will build both Dockerfiles present at the root of this template with a build context equal to the current directory (root directory).

The images are named according to a few rules :

- if the source Dockerfile has an extension (e.g. `Dockerfile.webserver`) then the image built will be named `<DOCKERHUB_USER>/<PROJECT>-<EXTENSION>` (e.g. : `Dockerfile.webserver` will produce an image named `<DOCKERHUB_USER>/<PROJECT>-webserver`)
- if the source Dockerfile hasn't an extension (e.g. `Dockerfile`) then the image built will be named `<DOCKERHUB_USER>/<PROJECT>`

A prerequisite to this target is that a file named `.dockerhub.yml` at the home directory (`~/.dockerhub.yml`) contains a `DOCKERHUB_USER` property.

```console
neon docker
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = 9791b0c79b55f2f34517d7e6b64d4900c8c7f2ce
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- docker --
adrienaury/go-template                                 refactor              sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   Less than a second ago   20.6MB
adrienaury/go-template-webserver                       refactor              sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   1 second ago    20.6MB
OK
```

To configure the docker target, edit the `build.yml` file and look for the property named `dockerfiles`.

The `dockerfiles` map can also be passed by the `neon -props` flag.

```console
$ neon -props '{dockerfiles: {"Dockerfile": "."}}' docker
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = refactor
COMMIT  = 9791b0c79b55f2f34517d7e6b64d4900c8c7f2ce
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = no
--------------------------------------------- docker --
adrienaury/go-template                                 refactor              sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   7 minutes ago   20.6MB
OK
```

#### Docker-tag target

Tag docker images. This target will run only if the tag being built is a release tag (vX.Y.Z).

It will tag all docker images with [semantic docker tags](https://medium.com/@mccode/using-semantic-versioning-for-docker-image-tags-dfde8be06699).

The build properties are the same as the [`docker`](#docker) target (a `dockerfiles` map).

```console
$ neon docker-tag
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = v0.2.0
COMMIT  = 00a4bdbf147a4394aa1e7f0483802f94658e9ce3
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = yes
VERSION = 0.2.0
--------------------------------------------- docker --
adrienaury/go-template                                 refactor              sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0.2.0                sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       refactor              sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0.2.0                sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
----------------------------------------- docker-tag --
adrienaury/go-template                                 refactor              sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0                    sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0.2                  sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0.2.0                sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       refactor              sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0                    sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0.2                  sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0.2.0                sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
OK
```

Use `-props '{latest: true}'` to include the latest tag.

```console
$ neon -props '{latest: true}' docker-tag
----------------------------------------------- info --
MODULE  = github.com/adrienaury/go-template
PROJECT = go-template
TAG     = v0.2.0
COMMIT  = 00a4bdbf147a4394aa1e7f0483802f94658e9ce3
DATE    = 2021-05-02
BY      = adrienaury@gmail.com
RELEASE = yes
VERSION = 0.2.0
--------------------------------------------- docker --
adrienaury/go-template                                 refactor              sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0.2.0                sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       refactor              sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0.2.0                sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
----------------------------------------- docker-tag --
adrienaury/go-template                                 latest                sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 refactor              sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0                    sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0.2                  sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template                                 v0.2.0                sha256:d05a1e1e5119aab03f3e3e33fa56d7db66ae5634beb53827b0e69fa168e3c595   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       latest                sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       refactor              sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0                    sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0.2                  sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
adrienaury/go-template-webserver                       v0.2.0                sha256:14b333b3679a64b3255e7c88e7211fa4b7502e2664e7b482373b392d5615414c   20 minutes ago   20.6MB
OK
```

#### Docker-push target

Publish tagged docker images to Dockerhub.

A prerequisite to this target is that a file named `.dockerhub.yml` at the home directory (`~/.dockerhub.yml`) contains a `DOCKERHUB_USER` property and a `DOCKERHUB_PASS` property.

The build properties are the same as the [`docker`](#docker) target and the [`docker-tag`](#docker-tag) target combined (a `dockerfiles` map and the `latest` boolean).

#### License target

Scan binaries for 3rd party licenses and generate a notice file, see the [example notice](NOTICE.md) generated by this target. Under the hood, the [golicense](https://github.com/mitchellh/golicense) tool is used.

This target work best if a file named `.github.yml` is present at the home directory (`~/.github.yml`) and contains a valid `GITHUB_TOKEN` property.

By default, this target write in `./NOTICE.md`. Use `-props '{noticefile: "a/different/filename.md"}'` to change the location and name of the notice file.

License scanning can report an error if an unallowed license is detected. Configure the build property `license`, or use the `neon -props` flag to exclude/include license by [SPDX identifier](https://spdx.org/licenses/), or to map unknown license (see [this link](https://github.com/mitchellh/golicense#configuration-file) for more information).

Example : `neon -props '{license: {deny: ["BSD-1-Clause"]}}' license`

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
