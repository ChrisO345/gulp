package gulp

import (
	"testing"
)

func TestGulpBuild(t *testing.T) {
	expected := true
	result := Gulp()
	if result != expected {
		t.Errorf("Expected %t but got %t", expected, result)
	}
}
