package sheaf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ManifestList struct {
	SchemaVersion int     `json:"schemaVersion"`
	Manifests     []Image `json:"manifests"`
}

type Image struct {
	MediaType   string            `json:"mediaType"`
	Size        int               `json:"size"`
	Digest      string            `json:"digest"`
	Annotations map[string]string `json:"annotations"`
}

func LoadFromIndex(indexPath string) ([]Image, error) {
	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}

	fmt.Println(indexPath)
	fmt.Println(string(data))

	var list ManifestList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	return list.Manifests, nil
}
