package dam

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	"path/filepath"
)

func ImageBundler(images []interfaces.RetrievedImages, metadata []interfaces.ImageMetadata) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	//Load images into zip file
	for _, image := range images {
		fileName := fmt.Sprintf("%s", filepath.Base(image.ObjectReference))
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create zip writer: %v", err)
		}

		_, err = fileWriter.Write(image.Image)
		if err != nil {
			return nil, fmt.Errorf("failed to write image: %v", err)
		}
	}

	//Load metadata into zip file
	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %v", err)
	}
	metadataFileWriter, err := zipWriter.Create("image_metadata.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata zip writer: %v", err)
	}
	_, err = metadataFileWriter.Write(metadataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to write metadata: %v", err)
	}

	//Closing the zip writer
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %v", err)
	}

	return buf.Bytes(), nil
}
