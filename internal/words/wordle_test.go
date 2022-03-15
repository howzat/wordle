package words

import (
	"context"
	"testing"

	"github.com/howzat/wordle/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileWordList(t *testing.T) {

	log, err := logging.NewProductionLogger("TestCompileWordList")
	require.NoError(t, err)

	files, err := NewWordSourceFiles("../../dictionaries")
	require.NoError(t, err)
	compiled, err := CompileWordList(context.TODO(), log, *files)

	assert.NoError(t, err)
	assert.Equal(t, 31941, compiled.Ingested)
}
