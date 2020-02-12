# sheaf

Inspired by [CNAB](https://cnab.io/), sheaf manages bundles of Kubernetes components.

Features:

* initialize a bundle
* package a bundle into an archive
* relocate images in a bundle archive to another registry
* generate manifests contained in a bundle archive

Limitations:

* detect images only in pod specs (apart from any added using `sheaf add-image`)
* replace images globally in manifests

[![asciicast](https://asciinema.org/a/yNVzkkpsVsUjT2jvSidfeqVHz.svg)](https://asciinema.org/a/yNVzkkpsVsUjT2jvSidfeqVHz)

## Getting started

This tool **IS A POC** (mistakes will be made)

### Install/Upgrade

`go get -u github.com/bryanl/sheaf/cmd/sheaf`

### Initialize Bundle

`sheaf init <bundle directory>`

Initial a sheaf project:
* create the directory if it does not exist, 
* create a bundle configuration
* create a manifests directory

### Add Manifests to Bundle

Repeat the following for each manifest (or pass multiple `-f` switches):

`sheaf add-manifest <bundle directory> -f <manifest path or URL>`

### Add Extra Images to Bundle

This is only necessary for images that `sheaf` can't find by scanning manifests, but it won't
do any harm if you add an image that can be found.

Repeat the following for each image (or pass multiple `-i` switches):

`sheaf add-image <bundle directory> -i <image>`

### Package Bundle

`sheaf pack <bundle directory> <archive path>`

Create an archive of `<bundle directory>` with all the images referenced in the manifests.

### Stage Bundle

`sheaf stage <archive path> <prefix>`

Stage the images located in the archive to a new registry with `<prefix>`. Images will be renamed.

### Generate Manifest

`sheaf gen-manifest <archive path> [--prefix=<prefix>]`

Generate manifests stored in the archive to stdout. If `<prefix>`, the images in the manifests will be
rewritten to the prefixed location. 
