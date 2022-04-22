package db

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/howzat/wordle"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetterMatches(t *testing.T) {

	tests := []struct {
		name       string
		results    []string
		dictionary []string
		search     string
		knowledge  []Knowlege
		err        error
	}{
		{
			name:      "providing no knowledge of the wordle is an error",
			search:    "blink",
			knowledge: nil,
			err:       errors.New("searching without search will match the entire dictionary"),
		}, {
			name:      "providing empty knowledge is an error",
			search:    "blink",
			knowledge: []Knowlege{None, None, None, None, None},
			err:       errors.New("searching without search will match the entire dictionary"),
		}, {
			name:       "providing 1 part knowledge should return all words containing that letter",
			search:     "blink",
			knowledge:  []Knowlege{Present, None, None, None, None},
			dictionary: []string{"beast", "crank", "dense", "sober"},
			results:    []string{"beast", "sober"},
		}, {
			name:       "providing 1 part knowledge matching 2 words should return all words containing those letters",
			search:     "blink",
			knowledge:  []Knowlege{Present, None, None, None, None},
			dictionary: []string{"beast", "crank", "dense", "sober"},
			results:    []string{"beast", "sober"},
		}, {
			name:       "providing 2 part knowledge should return all words containing those letters",
			search:     "blink", // assume the wordle is [b][l][oom <--no knowledge]
			knowledge:  []Knowlege{Present, Present, None, None, None},
			dictionary: []string{"beast", "crank", "dense", "slides"},
			results:    []string{"beast", "slides"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			log, err := wordle.NewProductionLogger(tt.name)
			require.NoError(t, err)

			db, err := NewIndex(*log, tt.dictionary, UseXXHashID)
			require.NoError(t, err)

			wordleDb := NewSearchEngine(db)
			search, err := NewWordleSearch(tt.search, tt.knowledge)
			require.NoError(t, err)

			matchResult, err := wordleDb.Search(*search)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.Equal(t, tt.results, matchResult.Items)
			}
		})
	}
}

func TestLetterMatchProps(t *testing.T) {

	log, err := wordle.NewProductionLogger("TestLetterMatchProps")
	require.NoError(t, err)

	words, err := loadWords("../cmd/search/dictionary.txt")
	require.NoError(t, err)

	db, err := NewIndex(*log, words, UseXXHashID)

	for i := 0; i < 1000; i++ {
		word := db.PickRandomWord()
		wdl, err := db.CandidateGuess(word)
		require.NoError(t, err)
		assert.NotEmpty(t, wdl)

		searchResults, err := db.Search(*wdl)
		allKnownLetter := wdl.KnownLetters()
		for _, result := range searchResults.Items {
			mustContainAllKnownLetters(t, allKnownLetter, result)
			if len(wdl.FullyKnownLetters()) > 0 {
				mustPreserveFullLetterMatches(t, *wdl, result)
			}
		}
	}
}

func loadWords(filepath string) ([]string, error) {
	fileReader, err := os.Open(filepath)
	if err != nil {
		return nil, wordle.WrapErr(err, "error reading file contents [%v]", filepath)
	}

	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)

	words := map[string]bool{}
	var uwords []string
	for scanner.Scan() {
		word := scanner.Text()
		if words != nil {
			if _, present := words[word]; !present {
				words[word] = true
				uwords = append(uwords, word)
			}
		}
	}

	sort.Strings(uwords)
	return uwords, nil
}

func mustPreserveFullLetterMatches(t *testing.T, wordle Wordle, result string) {
	t.Helper()
	for i, k := range wordle.knowledge {
		if k == Full {
			expected := string(wordle.letters[i])
			actual := string(result[i])
			if expected != actual {
				t.Fatalf(fmt.Sprintf("[%s] char [%s] at position[%d] was not in a matching position in [%s]", wordle.letters, expected, i, result))
			}
		}
	}
}

func mustContainAllKnownLetters(t *testing.T, chars []string, item string) {
	t.Helper()
	for _, ch := range chars {
		assert.Contains(t, item, ch)
	}
}

func TestBuildKnowledge(t *testing.T) {
	k := BuildKnowledgeForGuess("stick", "cider")
	assert.Equal(t, []Knowlege{Present, Present, None, None, None}, k)

	k1 := BuildKnowledgeForGuess("stick", "stick")
	assert.Equal(t, []Knowlege{Full, Full, Full, Full, Full}, k1)

	k2 := BuildKnowledgeForGuess("perch", "audio")
	assert.Equal(t, []Knowlege{None, None, None, None, None}, k2)
}
