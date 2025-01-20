package utils

import (
	"archive/tar"
	"encoding/json"
)

func AddJSONToTar(tw *tar.Writer, fileName string, data map[string]interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: fileName,
		Size: int64(len(jsonData)),
		Mode: 0644,
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = tw.Write(jsonData)
	return err
}

func AddFileToTar(tw *tar.Writer, fileName string, content []byte) error {
	header := &tar.Header{
		Name: fileName,
		Size: int64(len(content)),
		Mode: 0644,
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err := tw.Write(content)
	return err
}
