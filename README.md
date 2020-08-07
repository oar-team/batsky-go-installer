# batsky-go installer

This is a tool changes the imports of any go project, which enables time
synchronization with the kubernetes infrastrcuture simulator Batkube.

Its purpose it to enable compatibility between Kubernetes schedulers written in
Go and Batkube in order to study scheduling policies by running simulations of
a Kubernetes cluster.

## Usage

- `cd` into target project.
- `go run path/to/this/repo/cmd/main.go /scheduler/source/files/`
- `go mod tidy`
- `go mod vendor`
- `go run path/to/this/repo/cmd/main.go ./vendor/dependency1
./vendor/dependency2 ... [--not ./vendor/dependency3 ...]`

Disclaimer: Patch every dependency you judge critical, but do not replace
batsky-go's dependencies (github.com/pebbe/zmq4 and github.com/google/uuid) as
it would create circular dependencies.

Options:

- `--not` : ignores all the directories or files after this keyword
- To show the files it will replace beforehand, type in "show-files" to get a
    dry run that will just print the path to the files.


Note: If github.com/oar-team has not appeared in the vendor folder after
patching the source files, you need to manually copy over batsky-go source
files and its dependencies. (`clone` and `cd` into batsky-go project, `go mod
vendor`, copy `vendor` content under target project's `vendor` folder and
source files under `vendor/github.com/oar-team/batsky-go`)

## Why is this needed

Batsim controls the time (the simulation time) which means schedulers **must**
base their decisions according to this time rather than the machine time. This
is done by intercepting calls to the Go time library, to re-route time requests
to Batkube. This program modifies source code and vendored dependencies code in
order to redirect specific calls to github.com/oar-team/batsky-go.

## TODO

A script to automate the whole process.
