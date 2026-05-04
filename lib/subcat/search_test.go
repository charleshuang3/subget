package subcat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	results := Search("MIDV-999")

	// I don't care about the download count or size
	for i := range results {
		results[i].Downloads = 0
		results[i].Size = ""
	}
	assert.Contains(t, results, SearchResult{
		Title: "MIDV-999.chs.精校版 (translated from Chinese)",
		Path:  "subs/920/MIDV-999.chs.精校版.html",
	})
}
