package tests

import (
	"fmt"
	"goto/src/gpath"
	"testing"
)

func BenchmarkCheckRepeatedItems(b *testing.B) {
	// Generate a large list of paths
	count := 1000
	gpaths := make([]gpath.GotoPath, count)
	for i := 0; i < count; i++ {
		gpaths[i] = gpath.GotoPath{
			Path:         fmt.Sprintf("/path/%d", i),
			Abbreviation: fmt.Sprintf("p%d", i),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := gpath.CheckRepeatedItems(gpaths)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
