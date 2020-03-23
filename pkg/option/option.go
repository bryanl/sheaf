/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package option

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/codec"
	"github.com/bryanl/sheaf/pkg/fs"
	"github.com/bryanl/sheaf/pkg/remote"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// Runner is a sheaf command runner.
type Runner func(options ...sheaf.Option) error

// Options returns options that can be passed to the runner.
func Options(user ...sheaf.Option) []sheaf.Option {
	options := []sheaf.Option{
		sheaf.WithWriter(os.Stdout),
	}

	return append(options, user...)
}

// Generator generates options for a sheaf command.
type Generator struct {
	cmd    *cobra.Command
	prefix string
	m      map[string]func() []sheaf.Option
}

// NewGenerator creates an instance of Generator.
func NewGenerator(cmd *cobra.Command, runner Runner, prefix string) *Generator {
	r := reporter.New(reporter.WithWriter(os.Stdout))
	bundleImager := fs.NewBundleImager(fs.BundleImagerReporter(r))

	f := Generator{
		cmd:    cmd,
		prefix: prefix,
		m: map[string]func() []sheaf.Option{
			"default": func() []sheaf.Option {
				return []sheaf.Option{
					sheaf.WithBundleConfigCodec(fs.NewBundleConfigCodec()),
					sheaf.WithImageReplacer(fs.NewImageReplacer()),
					sheaf.WithBundleConfigWriter(fs.NewBundleConfigWriter()),
					sheaf.WithArchiver(archiver.New()),
					sheaf.WithBundleImager(bundleImager),
					sheaf.WithBundlePacker(fs.NewBundlePacker()),
					sheaf.WithCodec(codec.Default),
				}
			},
		},
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runner(f.Options()...)
	}

	return &f
}

// Options returns options for the Generator.
func (g Generator) Options() []sheaf.Option {
	var list []sheaf.Option
	for _, v := range g.m {
		list = append(list, v()...)
	}
	return Options(list...)
}

// WithArchive sets up the archive flag.
func (g Generator) WithArchive() {
	name := "archive"
	g.stringFlag(name, "", "archive path")
	g.setOptions(name, func() []sheaf.Option {
		archive := viper.GetString(g.flagName(name))
		return []sheaf.Option{
			sheaf.WithArchive(archive),
		}
	})
}

// WithBundleConfigFactory sets up options for creating a bundle config.
func (g Generator) WithBundleConfigFactory() {
	name := "bundle-config-factory"
	g.setOptions(name, func() []sheaf.Option {
		bundlePath := viper.GetString(g.flagName("bundle-path"))
		return []sheaf.Option{
			sheaf.WithBundleConfigFactory(func() sheaf.BundleConfig {
				// TODO: why does this need the default version
				bc := fs.NewBundleConfig(bundlePath, sheaf.BundleConfigDefaultVersion)
				return bc
			}),
		}
	})
}

// WithBundlePath sets a bundle path option.
func (g Generator) WithBundlePath() {
	name := "bundle-path"
	g.stringFlag(name, "", "bundle path")
	g.setOptions(name, func() []sheaf.Option {
		bundlePath := viper.GetString(g.flagName(name))
		if bundlePath == "" {
			bundlePath, _ = os.Getwd()
		}

		return []sheaf.Option{
			sheaf.WithBundlePath(bundlePath),
			sheaf.WithBundleFactory(func(bp string) (bundle sheaf.Bundle, err error) {
				return fs.NewBundle(bp)
			}),
		}
	})
}

// WithBundleName sets up a bundle name option.
func (g Generator) WithBundleName() {
	name := "bundle-name"
	g.stringFlag(name, "", "bundle name")
	g.setOptions(name, func() []sheaf.Option {
		s := viper.GetString(g.flagName(name))
		return []sheaf.Option{sheaf.WithBundleName(s)}
	})
}

// WithBundleVersion sets up a bundle version option.
func (g Generator) WithBundleVersion() {
	name := "bundle-version"
	g.stringFlag(name, "", "bundle version")
	g.setOptions(name, func() []sheaf.Option {
		s := viper.GetString(g.flagName(name))
		return []sheaf.Option{sheaf.WithBundleVersion(s)}
	})
}

// WithFilePaths sets up file path options.
func (g Generator) WithFilePaths() {
	name := "filename"
	g.stringSliceP(name, "f", nil, "filename (can specify multiple times)")
	g.setOptions(name, func() []sheaf.Option {
		s := viper.GetStringSlice(g.flagName(name))
		return []sheaf.Option{sheaf.WithFilePaths(s)}
	})
}

// WithForce sets up force option.
func (g Generator) WithForce() {
	name := "force"
	g.boolFlag(name, false, "force")
	g.setOptions(name, func() []sheaf.Option {
		tf := viper.GetBool(g.flagName(name))
		return []sheaf.Option{sheaf.WithForce(tf)}
	})
}

// WithImages sets up image options.
func (g Generator) WithImages() {
	name := "image"
	g.stringSliceP(name, "i", nil, "image (can specify multiple times)")
	g.setOptions(name, func() []sheaf.Option {
		s := viper.GetStringSlice(g.flagName(name))
		return []sheaf.Option{sheaf.WithImages(s)}
	})
}

// WithInitBundlePath sets up bundle path for a sheaf init.
func (g Generator) WithInitBundlePath() {
	name := "bundle-path"
	g.stringFlag(name, "", "bundle path")
	g.setOptions(name, func() []sheaf.Option {
		bundlePath := viper.GetString(g.flagName(name))

		if bundlePath == "" {
			bundlePath = viper.GetString(g.flagName("bundle-name"))
		}

		return []sheaf.Option{sheaf.WithBundleCreator(fs.BundleCreator(bundlePath))}
	})
}

// WithUserDefinedImage sets up user defined image options.
func (g Generator) WithUserDefinedImage() {
	g.stringFlag("api-version", "", "api version")
	g.stringFlag("kind", "", "kind")
	g.stringFlag("json-path", "", "json path")
	g.stringFlag("type", string(sheaf.SingleResult),
		fmt.Sprintf("type of user defined image (valid types: %s)",
			strings.Join(sheaf.UserDefinedImageTypes, ",")))
	g.setOptions("udi", func() []sheaf.Option {
		udi := sheaf.UserDefinedImage{
			APIVersion: viper.GetString(g.flagName("api-version")),
			Kind:       viper.GetString(g.flagName("kind")),
			JSONPath:   viper.GetString(g.flagName("json-path")),
			Type:       sheaf.UserDefinedImageType(viper.GetString(g.flagName("type"))),
		}

		return []sheaf.Option{
			sheaf.WithUserDefinedImage(udi),
		}
	})
}

// WithUserDefinedImageKey sets up user defined image key options.
func (g Generator) WithUserDefinedImageKey() {
	g.stringFlag("api-version", "", "api version")
	g.stringFlag("kind", "", "kind")
	g.setOptions("udi", func() []sheaf.Option {
		udiKey := sheaf.UserDefinedImageKey{
			APIVersion: viper.GetString(g.flagName("api-version")),
			Kind:       viper.GetString(g.flagName("kind")),
		}

		return []sheaf.Option{
			sheaf.WithUserDefinedImageKey(udiKey),
		}
	})
}

// WithPrefix sets up registry prefix option.
func (g Generator) WithPrefix() {
	name := "prefix"
	g.stringFlag(name, "", "registry prefix")
	g.setOptions(name, func() []sheaf.Option {
		prefix := viper.GetString(g.flagName(name))

		return []sheaf.Option{sheaf.WithRepositoryPrefix(prefix)}
	})
}

// WithInsecureRegistry sets up insecure registry option.
func (g Generator) WithInsecureRegistry() {
	name := "insecure-registry"
	g.boolFlag(name, false, "insecure registry")
	g.setOptions(name, func() []sheaf.Option {
		forceInsecure := viper.GetBool(g.flagName(name))
		irOpts := []remote.ImageReaderOption{
			remote.WithInsecure(forceInsecure),
		}

		opts := []remote.Option{
			remote.WithInsecureRegistry(forceInsecure),
		}

		ir := remote.NewImageReader(irOpts...)
		iw := remote.NewImageWriter(opts...)

		dryRun := viper.GetBool(g.flagName("dry-run"))

		return []sheaf.Option{
			sheaf.WithImageReader(ir),
			sheaf.WithImageWriter(iw),
			sheaf.WithImageRelocator(
				fs.NewImageRelocator(fs.ImageRelocatorDryRun(dryRun))),
		}
	})
}

// WithReference sets up a registry reference option.
func (g Generator) WithReference() {
	name := "ref"
	g.stringFlag(name, "", "reference")
	g.setOptions(name, func() []sheaf.Option {
		return []sheaf.Option{
			sheaf.WithReference(viper.GetString(g.flagName(name))),
		}
	})
}

// WithDestination sets up a destination option.
func (g Generator) WithDestination() {
	name := "dest"
	g.stringFlag(name, "", "destination")
	g.setOptions(name, func() []sheaf.Option {
		return []sheaf.Option{
			sheaf.WithDestination(viper.GetString(g.flagName(name))),
		}
	})
}

// WithDryRun sets up a dry run option.
func (g Generator) WithDryRun() {
	name := "dry-run"
	g.boolFlag(name, false, "dry run")
	g.setOptions(name, func() []sheaf.Option {
		return []sheaf.Option{
			sheaf.WithDryRun(viper.GetBool(g.flagName(name))),
		}
	})
}

func (g Generator) boolFlag(name string, value bool, usage string) {
	g.cmd.Flags().Bool(name, value, usage)
	g.bindFlag(name)
}

func (g Generator) stringFlag(name, value, usage string) {
	g.cmd.Flags().String(name, value, usage)
	g.bindFlag(name)
}

func (g Generator) stringSliceP(name, shorthand string, value []string, usage string) {
	g.cmd.Flags().StringSliceP(name, shorthand, value, usage)
	g.bindFlag(name)
}

func (g Generator) flagName(name string) string {
	return fmt.Sprintf("%s-%s", g.prefix, name)
}

func (g Generator) setOptions(name string, fn func() []sheaf.Option) {
	g.m[name] = fn
}

func (g Generator) bindFlag(name string) {
	if err := viper.BindPFlag(g.flagName(name), g.cmd.Flags().Lookup(name)); err != nil {
		panic(fmt.Sprintf("unable to bind %s in %s", name, g.prefix))
	}
}
