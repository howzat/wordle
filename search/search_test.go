package search

import (
	"errors"
	"testing"

	"github.com/howzat/wordle/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetterMatches(t *testing.T) {

	tests := []struct {
		name       string
		search     string
		results    []string
		dictionary []string
		knowledge  []MatchType
		err        error
	}{
		{
			name:      "providing no knowledge of the wordle is an error",
			search:    "brink",
			knowledge: nil,
			err:       errors.New("searching without knowledge will match the entire dictionary"),
		}, {
			name:      "providing knowledge that is empty is an error",
			search:    "brink",
			knowledge: []MatchType{None, None, None, None, None},
			err:       errors.New("searching without knowledge will match the entire dictionary"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			log, err := logging.NewProductionLogger(tt.name)
			require.NoError(t, err)

			db, err := NewIndexedDB(*log, tt.dictionary, UseXXHashID)
			require.NoError(t, err)

			store := NewSearchEngine(db)
			search, err := NewWordle(tt.search)
			require.NoError(t, err)

			matchResult, err := store.Search(*search)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.Equal(t, tt.knowledge, matchResult.LetterMatches)
			}
		})
	}
}
