# batsky-go installer

This repo is a tool aiming at changing imports of a Kubernetes scheduler, so as
to enable compatilibity with the Batkube interface to be able to run
simulations of a Kubernetes cluster.

The simulation are backed up by Batsim. Batkube is a software doing the link
between Batsim and Kubernetes schedulers.

## Usage
- `go build cmd/main.go`. I will refer to the built binary as `./main`.
- `cd` into target project.
- `./main /scheduler/source/files/`
- `go mod tidy`
- `go mod vendor`
- `./main ./vendor/dependency1 ./vendor/dependency2 ...`

Options :
- --not : ignores all the directories or files after this keyword
- To show the files it will replace beforehand, type in "showFiles" to get a
    dryRun that will just print the path to the files.

Notes :
- If github.com/oar-team has not appeared in the vendor folder, you need to
    manually copy over batsky-go source files and its dependencies. (to
    simplify the process, cd into batsky-go project, `go mod vendor`, copy over
    vendor content)
- Do not replace the entire vendor folder. It is not needed, and will create
    circular dependencies. In particular, do not replace github.com/pebbe/zmq4
    and github.com/google/uuid

## Why is this needed?
The Kubernetes ecosystem revolves around its central API server asynchronously,
which is in contradiction with the synchronous paradigm around which Batsim
revolves when exchanging with the schedulers.

More importantly, Batsim controls the time (the simulation time) which means
schedulers **must** base their decisions according to this time rather than the
machine time. This is done by hijacking calls to the Go time library, to
re-route time requests to Batkube.  This program modifies the source code and
vendored dependencies code in order to change specific calls to
github.com/oar-team/batsky-go rather than the standard time library.
