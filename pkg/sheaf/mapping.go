package sheaf

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/pivotal/image-relocation/pkg/image"
)

// FlattenRepoPathPreserveTagDigest maps the given Name to a new Name with a given repository prefix.
// It aims to avoid collisions between repositories and to include enough of the original name
// to make it recognizable by a human being. It preserves any tag and/or digest.
func FlattenRepoPathPreserveTagDigest(repoPrefix string, originalImage image.Name) (image.Name, error) {
	rn, err := FlattenRepoPath(repoPrefix, originalImage)
	if err != nil {
		return image.Name{}, err
	}

	// Preserve any tag
	if tag := originalImage.Tag(); tag != "" {
		var err error
		rn, err = rn.WithTag(tag)
		if err != nil {
			panic(err) // should never occur
		}
	}

	// Preserve any digest
	if dig := originalImage.Digest(); dig != image.EmptyDigest {
		var err error
		rn, err = rn.WithDigest(dig)
		if err != nil {
			panic(err) // should never occur
		}
	}

	return rn, nil
}

// FlattenRepoPath maps the given Name to a new Name with a given repository prefix.
// It aims to avoid collisions between repositories and to include enough of the original name
// to make it recognizable by a human being.
func FlattenRepoPath(repoPrefix string, originalImage image.Name) (image.Name, error) {
	hasher := md5.New()

	if _, err := hasher.Write([]byte(originalImage.Name())); err != nil {
		return image.Name{}, fmt.Errorf("hasher write: %w", err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	available := reference.NameTotalLengthMax - len(mappedPath(repoPrefix, "", hash))
	fp := flatPath(originalImage.Path(), available)
	var mp string
	if fp == "" {
		mp = fmt.Sprintf("%s/%s", repoPrefix, hash)
	} else {
		mp = mappedPath(repoPrefix, fp, hash)
	}
	mn, err := image.NewName(mp)
	if err != nil {
		panic(err) // should not happen
	}
	return mn, nil
}

func mappedPath(repoPrefix string, repoPath string, hash string) string {
	return fmt.Sprintf("%s/%s-%s", repoPrefix, repoPath, hash)
}

func flatPath(repoPath string, size int) string {
	return strings.Join(crunch(strings.Split(repoPath, "/"), size), "-")
}

func crunch(components []string, size int) []string {
	for n := len(components); n > 0; n-- {
		comp := reduce(components, n)
		if len(strings.Join(comp, "-")) <= size {
			return comp
		}

	}
	if len(components) > 0 && len(components[0]) <= size {
		return []string{components[0]}
	}
	return []string{}
}

func reduce(components []string, n int) []string {
	if len(components) < 2 || len(components) <= n {
		return components
	}

	tmp := make([]string, len(components))
	copy(tmp, components)

	last := components[len(tmp)-1]
	if n < 2 {
		return []string{last}
	}

	front := tmp[0 : n-1]
	return append(front, "-", last)
}
