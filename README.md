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

Create an archive of `<bundle directory>` with any images found by scanning the manifests together with any listed in `bundle.json`.

### Stage Bundle

`sheaf stage <archive path> <prefix>`

Stage the images located in the archive to a new registry with `<prefix>`. Images will be renamed.

### Generate Manifest

`sheaf gen-manifest <archive path> [--prefix=<prefix>]`

Generate manifests stored in the archive to stdout. If `<prefix>` is specified, the images in the manifests will be
rewritten to the prefixed location. 

## Finding images

There are myriad ways to specify an image in a manifest. `sheaf` can detect images defined in pod specs that are in
Pods themselves or in a pod Spec template (e.g., in a Deployment). This heuristic works in a large number of cases. 
With Kubernetes and Custom Resource Definitions it is possible to define images in other locations as well. `sheaf`
has a method called "user defined images", that allows custom locations to be created. 

Example:

```json
{
    "name": "knative-serving-0.12",
    "version": "0.1.0",
    "schemaVersion": "v1alpha1",
    "images": [],
    "userDefinedImages": [
        {
            "apiVersion": "caching.internal.knative.dev/v1alpha1",
            "kind": "Image",
            "jsonPath": "{.spec.image}",
            "type": "single"
        }
    ]
}
```
 
In this case, when `sheaf` is parsing the `Image` kind API version `caching.internal.knative.dev/v1alpha1` it will
use the JSON path `{.spec.image}` to locate an image. 

```yaml
---
apiVersion: caching.internal.knative.dev/v1alpha1
kind: Image
metadata:
  labels:
    serving.knative.dev/release: "v0.12.0"
  name: queue-proxy
  namespace: knative-serving
spec:
  image: gcr.io/knative-releases/knative.dev/serving/cmd/queue@sha256:3932262d4a44284f142f4c49f707526e70dd86317163a88a8cbb6de035a401a9
```

In this case `gcr.io/knative-releases/...` will be identified as an image. There is also a mechanism to find multiple 
images at the same time. If the `userDefinedImage` type is "multiple", then the output from the JSON path should be 
a comma separated value.   

Examples of JSON path queries for multiple values:
* `{range .spec.images[*]}{@}{','}{end}`: Looks for an array of images in `.spec.images`
* `{range ..spec.containers[*]}{.image}{','}{end}`: This is the method that `sheaf` uses to find images in a Pod spec
template
