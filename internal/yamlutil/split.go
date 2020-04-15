/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package yamlutil

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/bryanl/sheaf/third_party/github.com/kubernetes/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/yaml"
)

// Split separates a stream of YAML documents into a slice of single documents.
func Split(in []byte) ([][]byte, error) {

	// The following implementation isn't great (see the tests), but it works
	// some of the time and is probably better than reimplementing from scratch.
	// See https://github.com/kubernetes/apimachinery/issues/91.

	inReader := ioutil.NopCloser(bytes.NewReader(in))

	docReader := yaml.NewDocumentDecoder(inReader)
	defer docReader.Close()

	docs := [][]byte{}
	doc := []byte{}
	buf := make([]byte, 4096)
	for {
		n, err := docReader.Read(buf)
		doc = append(doc, buf[0:n]...)
		switch err {
		case nil: // end of chunk returned
			docs = append(docs, doc)
			doc = []byte{}
		case io.ErrShortBuffer: // partial chunk returned
			continue
		case io.EOF: // either end of chunk or no bytes returned
			if len(doc) > 0 {
				docs = append(docs, doc)
			}
			return docs, nil
		default:
			return [][]byte{}, err
		}
	}
}
