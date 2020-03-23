/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"io"
	"os"

	"github.com/bryanl/sheaf/pkg/reporter"
)

type options struct {
	reporter reporter.Reporter
	codec    Codec

	bundlePath    string
	bundleFactory BundleFactoryFunc
	bundleName    string
	bundleVersion string

	bundleConfigFactory func() BundleConfig
	bundleConfigCodec   BundleConfigCodec
	bundleConfigWriter  func() (BundleConfigWriter, error)
	bundlePacker        func() (BundlePacker, error)

	archiver Archiver

	createBundle func(bc BundleConfig) error

	repositoryPrefix string

	imageReplacer  ImageReplacer
	imageRelocator ImageRelocator

	userDefinedImage    UserDefinedImage
	userDefinedImageKey UserDefinedImageKey

	bundleImager BundleImager
	imageReader  ImageReader
	imageWriter  ImageWriter

	filePaths   []string
	images      []string
	force       bool
	reference   string
	destination string
	archive     string

	dryRun bool

	writer io.Writer
}

func makeDefaultOptions(list ...Option) options {
	opts := options{
		bundleVersion: BundleConfigDefaultVersion,
		// TODO: combine writer and reporter
		writer:   os.Stdout,
		reporter: reporter.New(reporter.WithWriter(os.Stdout)),
		bundleFactory: func(string) (bundle Bundle, err error) {
			return nil, fmt.Errorf("bundle factory is not configured")
		},
		bundlePacker: func() (packer BundlePacker, err error) {
			return nil, fmt.Errorf("bundle packer is not configured")
		},
		bundleConfigWriter: func() (writer BundleConfigWriter, err error) {
			return nil, fmt.Errorf("bundle config writer is not configured")
		},
	}
	for _, o := range list {
		o(&opts)
	}
	return opts
}

// Option is a functional option for configuring a sheaf command.
type Option func(o *options)

// WithArchive set archive.
func WithArchive(archive string) Option {
	return func(o *options) {
		o.archive = archive
	}
}

// WithArchiver sets archiver.
func WithArchiver(archiver Archiver) Option {
	return func(o *options) {
		o.archiver = archiver
	}
}

// WithBundleImager sets the bundle imager.
func WithBundleImager(bi BundleImager) Option {
	return func(o *options) {
		o.bundleImager = bi
	}
}

// WithBundleName sets bundle name.
func WithBundleName(name string) Option {
	return func(o *options) {
		o.bundleName = name
	}
}

// WithBundleVersion sets bundle version.
func WithBundleVersion(version string) Option {
	return func(o *options) {
		if version == "" {
			version = BundleConfigDefaultVersion
		}
		o.bundleVersion = version
	}
}

// BundleCreatorFunc is a function a function that creates a bundle.
type BundleCreatorFunc func(bc BundleConfig) error

// WithBundleCreator sets bundle creator.
func WithBundleCreator(fn BundleCreatorFunc) Option {
	return func(o *options) {
		o.createBundle = fn
	}
}

// BundleFactoryFunc is a factory that creates bundle instances.
type BundleFactoryFunc func(rootPath string) (Bundle, error)

// WithBundleFactory sets the bundle factory.
func WithBundleFactory(fn BundleFactoryFunc) Option {
	return func(o *options) {
		o.bundleFactory = fn
	}
}

// BundleConfigFactory is a function that creates a bundle config.
type BundleConfigFactory func() BundleConfig

// WithBundleConfigFactory sets the bundle config factory.
func WithBundleConfigFactory(fn BundleConfigFactory) Option {
	return func(o *options) {
		o.bundleConfigFactory = fn
	}
}

// WithBundleConfigWriter sets the bundle config writer..
func WithBundleConfigWriter(bcw BundleConfigWriter) Option {
	return func(o *options) {
		o.bundleConfigWriter = func() (writer BundleConfigWriter, err error) {
			return bcw, nil
		}
	}
}

// WithBundlePacker sets the bundle packer.
func WithBundlePacker(bp BundlePacker) Option {
	return func(o *options) {
		o.bundlePacker = func() (packer BundlePacker, err error) {
			return bp, nil
		}
	}
}

// WithBundlePath set the bundle path.
func WithBundlePath(s string) Option {
	return func(o *options) {
		o.bundlePath = s
	}
}

// WithCodec sets the codec.
func WithCodec(c Codec) Option {
	return func(o *options) {
		o.codec = c
	}
}

// WithDestination sets the destination.
func WithDestination(dest string) Option {
	return func(o *options) {
		o.destination = dest
	}
}

// WithDryRun sets dry run.
func WithDryRun(dryRun bool) Option {
	return func(o *options) {
		o.dryRun = dryRun
	}
}

// WithFilePaths sets file paths.
func WithFilePaths(filePaths []string) Option {
	return func(o *options) {
		o.filePaths = filePaths
	}
}

// WithForce sets force.
func WithForce(force bool) Option {
	return func(o *options) {
		o.force = force
	}
}

// WithImages sets images.
func WithImages(imageList []string) Option {
	return func(o *options) {
		o.images = imageList
	}
}

// WithImageRelocator sets image relocator.
func WithImageRelocator(ir ImageRelocator) Option {
	return func(o *options) {
		o.imageRelocator = ir
	}
}

// WithImageReader sets image reader.
func WithImageReader(ir ImageReader) Option {
	return func(o *options) {
		o.imageReader = ir
	}
}

// WithImageWriter sets image writer.
func WithImageWriter(iw ImageWriter) Option {
	return func(o *options) {
		o.imageWriter = iw
	}
}

// WithRepositoryPrefix sets the repository prefix.
func WithRepositoryPrefix(prefix string) Option {
	return func(o *options) {
		o.repositoryPrefix = prefix
	}
}

// WithImageReplacer sets the image replacer.
func WithImageReplacer(ir ImageReplacer) Option {
	return func(o *options) {
		o.imageReplacer = ir
	}
}

// WithReference sets a registry reference.
func WithReference(ref string) Option {
	return func(o *options) {
		o.reference = ref
	}
}

// WithReporter sets the reporter.
func WithReporter(r reporter.Reporter) Option {
	return func(o *options) {
		o.reporter = r
	}
}

// WithUserDefinedImage sets user defined image.
func WithUserDefinedImage(udi UserDefinedImage) Option {
	return func(o *options) {
		o.userDefinedImage = udi
	}
}

// WithUserDefinedImageKey sets user defined image key.
func WithUserDefinedImageKey(udiKey UserDefinedImageKey) Option {
	return func(o *options) {
		o.userDefinedImageKey = udiKey
	}
}

// WithWriter sets writer.
func WithWriter(w io.Writer) Option {
	return func(o *options) {
		o.writer = w
	}
}

// WithBundleConfigCodec sets bundle config codec.
func WithBundleConfigCodec(bcc BundleConfigCodec) Option {
	return func(o *options) {
		o.bundleConfigCodec = bcc
	}
}
