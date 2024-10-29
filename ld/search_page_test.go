package ld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	page, _ := Search("nonosensesense")

	searchPage, _ := page.(SearchPage)

	expected := []string{
		"consensuses",
		"no-nonsense",
		"nonsenses",
		"horse sense",
		"nonsense",
		"common sense",
		"condensers",
		"condenses",
		"consensus",
		"denseness",
	}

	assert.Equal(t, expected, searchPage.Results)
}
