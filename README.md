# sheaf

Inspired by [CNAB](https://cnab.io/), sheaf manages bundles of Kubernetes components.

Features:

* initialize a bundle
* package a bundle into an archive
* relocate images in a bundle archive to another registry
* generate manifests contained in a bundle archive

## Getting started

This tool **IS A POC** (mistakes will be made)

### Install/Upgrade

`go get -u github.com/bryanl/sheaf/cmd/sheaf`

### Initialize Bundle

`sheaf init <directory>`

Initial a sheaf project:
* create the directory if it does not exist, 
* create a bundle configuration
* create a manifests directory

### Package Bundle

`sheaf pack <directory> <archive path>`

Create an archive of `<directory>` with all the images referenced in the manifests.

### Stage Bundle

`sheaf stage <archive path> <prefix>`

Stage the images located in the archive to a new registry with `<prefix>`. Images will be renamed.

### Generate Manifest

`sheaf gen-manifest <archive path> [--prefix=<prefix>]`

Generate manifests stored in the archive to stdout. If `<prefix>`, the images in the manifests will be rewritten to the prefixed location. 
