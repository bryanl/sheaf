package sheaf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cnabio/duffle/pkg/imagestore/ocilayout"
	"github.com/deislabs/duffle/pkg/imagestore"
	dcopy "github.com/otiai10/copy"
)

// Bundle represents a bundle
type Bundle struct {
	// Path is the path to the bundle directory.
	Path string
	// Config is the BundleConfig for the bundle.
	Config BundleConfig
	// Store is the image store
	Store imagestore.Store

	// tmpDir for temporary things.
	tmpDir string
}

// OpenBundle loads a bundle. Call Bundle.Close() to ensure workspace is cleaned up.
func OpenBundle(path string) (*Bundle, error) {
	// check if directory exists
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("bundle directory %q does not exist", path)
		}

		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", path)
	}

	bundleConfigFilename := filepath.Join(path, BundleConfigFilename)

	bundleConfig, err := LoadBundleConfig(bundleConfigFilename)
	if err != nil {
		return nil, fmt.Errorf("load bundle config: %w", err)
	}

	tmpDir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return nil, fmt.Errorf("create temp directory: %w", err)
	}

	root := filepath.Join(tmpDir, filepath.Base(path))
	if err := dcopy.Copy(path, root); err != nil {
		return nil, fmt.Errorf("stage bundle: %w", err)
	}

	store, err := ocilayout.Create(setStoreLocation(root))
	if err != nil {
		return nil, fmt.Errorf("create image store: %w", err)
	}

	bundle := &Bundle{
		Path:   root,
		Config: bundleConfig,
		Store:  store,
		tmpDir: tmpDir,
	}

	return bundle, nil
}

func setStoreLocation(archiveDir string) imagestore.Option {
	return func(parameters imagestore.Parameters) imagestore.Parameters {
		parameters.ArchiveDir = archiveDir
		return parameters
	}
}

func (b *Bundle) Images() ([]string, error) {
	seen := make(map[string]bool)

	// assume manifests live in `app/manifests`
	manifestsPath := filepath.Join(b.Path, "app", "manifests")
	entries, err := ioutil.ReadDir(manifestsPath)
	if err != nil {
		return nil, fmt.Errorf("read manifests dir %q: %w", manifestsPath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(manifestsPath, entry.Name())
		images, err := ContainerImages(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("find container images for %q: %w", manifestPath, err)
		}

		fmt.Printf("Images in %s: [%s]\n", entry.Name(), strings.Join(images, ","))
		for i := range images {
			seen[images[i]] = true
		}
	}

	var list []string
	for k := range seen {
		list = append(list, k)
	}

	return list, nil
}

// Bundle writes archive to disk.
func (b *Bundle) Write() error {
	outputFile := b.Config.Filename(".")
	fmt.Println("Creating archive: ", outputFile)
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("unable to write %s: %v", outputFile, err)
		}
	}()

	if err := Archive(b.Path, f); err != nil {
		return fmt.Errorf("create archive: %w", err)
	}

	return nil
}

func (b *Bundle) Close() error {
	if err := os.RemoveAll(b.tmpDir); err != nil {
		return fmt.Errorf("remove temporary directory")
	}

	return nil
}
