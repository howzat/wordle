package words

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLineSeperatedDictionary(t *testing.T) {
	var contents = `AARDVARK
Brioche
123
Ma.ch
"/"
/
&
☞
😂
cafés
april
    April
APRIL`

	file, tidyFn := createTempFile(t, contents)

	defer tidyFn()

	wordSelect := ParseLineSeperatedDictionary(file.Name())
	words, err := wordSelect(NormaliseWord, WordleCandidate)
	require.NoError(t, err)
	assert.EqualValues(t, []string{"april", "april", "april"}, words)
}

func TestParseWordsetDictionary(t *testing.T) {
	var contents = `{"attentively": {}, "april": {}, "after hours": {}, "APRIL": {}, "April": {}}`

	file, tidyFn := createTempFile(t, contents)

	defer tidyFn()

	wordSelect := ParseWordsetDictionary(file.Name())
	words, err := wordSelect(NormaliseWord, WordleCandidate)
	require.NoError(t, err)
	assert.EqualValues(t, []string{"april", "april", "april"}, words)
}

func createTempFile(t *testing.T, contents string) (*os.File, func()) {
	file, err := ioutil.TempFile(".", "tmp")

	if _, err = file.Write([]byte(contents)); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	require.NoError(t, err)

	return file, func() {
		_ = os.Remove(file.Name())
	}
}
