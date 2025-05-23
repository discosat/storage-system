package dim

import (
	"log"
	"os"
	"testing"
)

func TestExtractMetadata(t *testing.T) {

	rawWithMetadata, _ := os.ReadFile("./testImages/forgor.jpg")
	rawNoMetadata, err := os.ReadFile("./testImages/Oense.tif")
	if err != nil {
		log.Println("Files not read")
		t.Fatal(err)
	}

	// Given
	tests := []struct {
		rawObservation []byte
		wantError      bool
	}{
		{rawWithMetadata, false},
		//{rawNoMetadata, true},
		{rawNoMetadata, true},
	}

	for _, testCase := range tests {
		_, _, err := extractMetadata(testCase.rawObservation)
		if (err != nil) != testCase.wantError {
			t.Errorf("extractMetadata() error = %v, wantError = %v", err, testCase.wantError)
		}
	}
}
