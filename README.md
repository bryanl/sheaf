# sheaf

Manages bundles of Kubernetes components.

`sheaf` is a tool that can create a `bundle` of Kubernetes components. It can generate an archive from the bundle
that can be distributed for use in Kubernetes clusters. The initial idea was inspired by inspired by 
[CNAB](https://cnab.io/). It answers the question: _how can I distribute Kubernetes manifests with their 
associated images_?

## Why sheaf?

* You want to create a distributable archive of Kubernetes manifests and their images
* The Kubernetes clusters you are working with are potentially air-gapped and don't have access to the Internet.

Features:

* Manages a `bundle` of Kubernetes manifests
* Package a bundle into an archive
* Relocate images referenced in a bundle archive to another registry
* Generate manifests contained in a bundle archive with optionally relocated images


## Getting started

`sheaf` is currently in an alpha state. We are releasing the tool early as a preview.

### Install/Upgrade

`go get -u github.com/bryanl/sheaf/cmd/sheaf`

### Initialize Bundle

```sh
sheaf init <bundle directory>
```

Initialize a sheaf project:

* Creates the directory bundle directory.
* Creates a bundle configuration in _project_/`bundle.json`
* Creates a manifests directory in _project_/`app/manifests`

### Add Manifests to Bundle

Repeat the following for each manifest (or pass multiple `-f` switches):

`sheaf manifest add <bundle directory> -f <manifest path or URL>`

### Package Bundle

`sheaf archive pack <bundle directory> <archive path>`

Create an archive of `<bundle directory>` with any images found by scanning the manifests together with any listed in `bundle.json`.

### Stage Bundle

`sheaf archive reloate <archive path> <prefix>`

Relocate the images located in the archive to a registry repository with `<prefix>`. Images will be renamed.

### Generate Manifest

`sheaf manifest show <archive path> [--prefix=<prefix>]`

Generate manifests stored in the archive to stdout. If `<prefix>` is specified, the images in the manifests will be
rewritten to the prefixed location. 

### Create user defined images.

With Custom Resource Definitions, it is possible to define locations that `sheaf` cannot detect automatically. `sheaf`
allows the user to create user defined images. For example, if you have a custom resource with API version
`x.bryanl.dev/v1` and kind `Config` with a container image specified at `.spec.image`, you can use the following 
command to update the `sheaf` configuration.

```sh
sheaf config set-udi --bundle-path project-path --api-version x.bryanl.dev/v1 --kind Config --json-path '{.spec.image}'
```

There is also support for arrays of images as well.

```sh
sheaf config set-udi --bundle-path project-path \
  --api-version x.bryanl.dev/v1 \
  --kind SecondaryConfig \
  --json-path '{.spec.images}'
  --type multiple
```

When `sheaf` is building a bundle archive or generating manifests, it will use the user defined mappings.


### Add Extra Images to Bundle

This is only necessary for images that `sheaf` can't find by scanning manifests, but it won't
do any harm if you add an image that can be found.

Repeat the following for each image (or pass multiple `-i` switches):

`sheaf config add-image <bundle directory> -i <image>`

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

## Notices/Limitations:

* Detected images are replaced everywhere in your document in order to preserve comments.
