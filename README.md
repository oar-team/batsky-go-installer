# Batsky time hijack

This repo is a tool aiming at changing imports of a Kubernetes scheduler, so as
to enable compatilibity with the Batkube interface to be able to run
simulations of a Kubernetes cluster.

The simulation are backed up by Batsim. Batkube is a software doing the link
between Batsim and Kubernetes schedulers.

## Usage
**TODO**

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

**TODO** : lists these calls

