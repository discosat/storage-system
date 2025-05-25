package dim

import (
	"log"
	"os"
	"testing"
)

func TestExtractMetadata(t *testing.T) {

	rawWithMetadata, err := os.ReadFile("./testData/forgor.jpg")
	if err != nil {
		log.Println("Files not read")
		t.Fatal(err)
	}
	rawNoMetadata, err := os.ReadFile("./testData/Oense.tif")
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
		t.Logf("Test: %v, success", t.Name())
	}
}
