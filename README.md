# sheaf

Manages bundles of Kubernetes components.

`sheaf` is a tool that can create a `bundle` of Kubernetes components. It can generate an archive from the bundle
that can be distributed for use in Kubernetes clusters. The initial idea was inspired by
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
sheaf init --bundle-path <bundle directory> --bundle-name <your project name>
```

Initialize a sheaf project:

* Creates the directory bundle directory.
* Creates a bundle configuration in _project_/`bundle.json`
* Creates a manifests directory in _project_/`app/manifests`

### Add Manifests to Bundle

Repeat the following for each manifest (or pass multiple `-f` switches):

`sheaf manifest add --bundle-path <bundle directory> -f <manifest path or URL>`

### Package Bundle

`sheaf archive pack --bundle-path <bundle directory> --dest <archive output directory>`

Create an archive of `<bundle directory>` with any images found by scanning the manifests together with any listed in `bundle.json` and output it in the _<archive output directory>_ directory.
Note that the directory _<archive output directory>_ must exist.

For an example of what appears in the archive, see below.

### Stage Bundle

`sheaf archive relocate --archive <archive path> --prefix <prefix>`

Relocate the images located in the archive to a registry repository with `<prefix>`. Images will be renamed and pushed to the new registry.

### Generate Manifest

`sheaf manifest show --bundle-path <bundle directory> [--prefix=<prefix>]`

Generate manifests stored in the archive to stdout. If `<prefix>` is specified, the images in the manifests will be
rewritten to the prefixed location. 

### Create user defined images

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


## Finding images

There are myriad ways to specify an image in a manifest. `sheaf` can detect images defined in pod specs that are in
Pods themselves or in a pod Spec template (e.g., in a Deployment). This heuristic works in a large number of cases. 
With Kubernetes and Custom Resource Definitions it is possible to define images in other locations as well. `sheaf`
has a method called "user defined images", that allows custom locations to be created. 

Example:

```sh
sheaf config set-udi --bundle-path project-path --api-version caching.internal.knative.dev/v1alpha1 --kind Image --json-path '.spec.image'
```
 
In this case, when `sheaf` is parsing the `Image` kind API version `caching.internal.knative.dev/v1alpha1` it will
use the JSON path `.spec.image` to locate an image. In this case:

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

sheaf will identify `gcr.io/knative-releases/...` as an image.

Examples of JSON path queries:
* `.spec.images[*]`: Looks for an array of images in `.spec.images`
* `..spec.containers[*].image`: This is the method that `sheaf` uses to find images in a Pod spec template

## So what's in an archive?

The following detailed example shows the contents of an archive. This may change over time as `sheaf` evolves, but it should give some insight.
```
$ ./sheaf init --bundle-path scratch/example --bundle-name example
$ ./sheaf manifest add --bundle-path scratch/example -f scratch/deployment.yaml
Adding manifest from scratch/deployment.yaml
$ cat scratch/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-deployment
spec:
  template:
     spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
$ ./sheaf archive pack --bundle-path scratch/example --dest scratch
Staging bundle configuration
Staging manifests
Staging images
  deployment.yaml
  └─ docker.io/library/nginx
  adding docker.io/library/nginx to layout
Creating archive: scratch/example-0.1.0.tgz
$ tar xvzf scratch/example-0.1.0.tgz --directory scratch/example-expanded
$ tree scratch/example-expanded
scratch/example-expanded
├── app
│   └── manifests
│       └── deployment.yaml
├── artifacts
│   └── layout
│       ├── blobs
│       │   └── sha256
│       │       ├── 00dbf0056e4065d3e66bbfbd42551123cabfd4b75859f586bfaea443e279d07a
│       │       ├── 015f0cb643493d580fe878b28ba8b33f17752aa349f417f80c022514e2b50a8c
│       │       ├── 02997b777612b24efedac3f1f219afcd6bcf769108563209573e3471bd534478
│       │       ├── 0598bd3a8f718f1b030930aa692ca6f9095d743c57b3b4aa3a200c362dc0de02
│       │       ├── 06085c693a49adbc2f75ef87074579df1412d39da48d1c4203b0a8370c67125b
│       │       ├── 0901fa9da894a8e9de5cb26d6749eaffb67b373dc1ff8a26c46b23b1175c913a
│       │       ├── 092368e66229d2df31f2f6980d4abe630bc4576cc830485ce9aec1ad6ee39a7b
│       │       ├── 0fc13119a4a310bba2bfeab8e68068dfce512f8871050555a254988397aafa1e
│       │       ├── 18a64c9928215257ea8bcbbc2aee0e41b022bed61791291d141aaa2edcb5aad7
│       │       ├── 1cf27aa8120baed22838ca29cdacc71bb16278cd846d91c62479ed72eeb23463
│       │       ├── 230cf2b89e1f15b73a321a8b7637583cc0c6bbf0f948ef291d1ed844126c4635
│       │       ├── 2b94f09d2d25fffb610349b7e7d26c0b1bc6a7d55795906a34e346c52c464fd1
│       │       ├── 2db5eb6422bd5ac7462c19ab7b1f8d89a41873216bc583a9a6821c1a478337c8
│       │       ├── 2dd003996c9ab82cac8112be0a4c04068e666e7a5d0cce3c65fb8f064de284e7
│       │       ├── 33cc09c9b190539635d7c971301f623d94fda5b4b5647966c6c240902119009f
│       │       ├── 36a1077091ada6e640493ea6da4c304d7b32f7341dcb5be2bbb27c45dd20b77b
│       │       ├── 3d3726baabd4a64466564b8068c1af39ca69f003bd0662f0c30e6253666290ab
│       │       ├── 405e75bf6bb0104d67fcebf58e07cd21bf344589df9c1a41c00354a60ea3a604
│       │       ├── 40969e979bd83d4599eb92444aefa1abeaaf64186a636f4e4c50a181e680d360
│       │       ├── 43362c1f9f7a7afc4e76b3906fb9dd7e605432ce6f00137d5aeb7420fede16ee
│       │       ├── 4584011f2cd184c8af36b3111f56cd7830a7bfda446d2112d28ec0dfff91082a
│       │       ├── 4ae502311710257cf63b1020256729f1f474f320b6e768284be83182877f3244
│       │       ├── 4ce8bb5bca50e30a1b523de188d7538f5f0e17693dda87462c434aa7174d5b57
│       │       ├── 565c5b948bc1162b1b8bb1699830e79345a5fe06ce7a6e97a137b9837e9422e9
│       │       ├── 57a619d34e5582111a25d6e29783e7cef7858c05d82e446c9654698a51d86457
│       │       ├── 57d8e76b728c535237497f00d57f2e5a000eed287fdf975ffcf090cd2a094061
│       │       ├── 5b11fc09c1a26b11a7df7d593adf43baff53c5cdba71cf8a87ae4a6dd17eb52c
│       │       ├── 5e38a7a36675286bba77c3b1a48fe3f994195a672f8f92c4f01635c3f05b0ecf
│       │       ├── 63dbb66c5119bb5086d9e6fb6b154211afc20b44ed136ab7df808f6044cfc6f1
│       │       ├── 644cef618ceea1cc2a7c5aa22d188fcdedeffc2ac89896fffdec74506615a0fc
│       │       ├── 66979ef6cabd7bb5aa8a418f1d1a8127ebd5ac02f910a14def8f33738041c4a7
│       │       ├── 67d252a8c1e1e940e1b2c52e05046403d00e84a33f6ddc94dc9c8870110338d3
│       │       ├── 67ddcee0b1be2b67823856a6aa49c601ca5783d693f4724b6a8f7e413d44a73a
│       │       ├── 68c6bf98048c7f913ee5e59b4f913be0b363d4ca0644ead4287202d56cbdf4d9
│       │       ├── 68ec9ddfb393aee91093cbf859cde906459cbf7218fc988f2990f21b22f8e521
│       │       ├── 6a03ff813dcd7a821ea0d229b54798547a6062c7f26f54dc8c1e2981c01f5e58
│       │       ├── 6a4ba5c72a8ea85ea434c8b2039f9b4a776fb2463fff084b0cff53c07bcc58e6
│       │       ├── 71719749d285a416f8253bf524980b4a96250add89b03e9af3cb8d5ffe293ca9
│       │       ├── 71d2ad811a3bf658598b1f8d635024971b15b673e50891446f6ac97ad9ecfae9
│       │       ├── 76f2a0c8ea98ebd699a77ced9c677e97cd54b038a8c5e89670af78f38b047b33
│       │       ├── 8559a31e96f442f2c7b6da49d6c84705f98a39d8be10b3f5f14821d0ee8417df
│       │       ├── 860f8957d8be856e2235a28e49fc4dca17254951e0eb67d760769755656f5cad
│       │       ├── 8a8b84398062afd8652193b566914890f836f8a9303eeb278a4dce93e2223a65
│       │       ├── 8ff4598873f588ca9d2bf1be51bdb117ec8f56cdfd5a81b5bb0224a61565aa49
│       │       ├── 9a2ce85454764788d4e69016d9a7b74808619a56a00c4a1685d685655d0eb78c
│       │       ├── 9c2b660fcff6200198b5a6dd0f4acb51d27813a959f6a38a2ab089b5d058e1d7
│       │       ├── a52fbcfc43b142e3f86683bf1b2b8da0f58f0f52cf1584d236f47a2a4cadbf6f
│       │       ├── a93c8a0b0974c967aebe868a186e5c205f4d3bcb5423a56559f2f9599074bbcd
│       │       ├── b38125636144fda2e96c14893240ff2a307f3743ea83d297905289e49523a088
│       │       ├── c2cb1aebb3e93bed04cee202f88e8c1932169beae7a164cd2b5a5d8b356b1c35
│       │       ├── c2e53b84630c36ddf2ee1993636b89d678a43f75d27c3e2298f77ccd06b6594e
│       │       ├── cc0e3936130331cb326019018d2761610ae2e018a07fba42d65ce63b1ad3e657
│       │       ├── ccdad9502600ef5f2a2b04a4c5a9d94ae9e0b7bd6f9a1588090193964819916b
│       │       ├── ce8699553d8b3417193659018fb5fb4ada22ed879ec1650f47f5b769e97a9800
│       │       ├── efe6d8b8ae3a94bd3b6e7e68a7ab7749e099231ad9d7b13cbccf8af538fb04dd
│       │       ├── f3e35b5be24177bc7f2c19401e9b45e8e834795815a982187c680643037064ed
│       │       └── f58da03af52f1386c795c25912f0835a91884c5cbe62963d7c1438b2259095a3
│       ├── index.json
│       └── oci-layout
└── bundle.json
$ cat scratch/example-expanded/bundle.json
{
  "schemaVersion": "v1alpha1",
  "name": "example",
  "version": "0.1.0"
}
$ cat scratch/example-expanded/app/manifests/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-deployment
spec:
  template:
     spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
```
The contents of scratch/example-expanded/artifacts/layout is a [standard OCI image layout](https://github.com/opencontainers/image-spec/blob/master/image-layout.md).
