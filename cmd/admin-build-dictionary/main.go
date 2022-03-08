package main

import (
	"context"
	"os"

	"github.com/howzat/wordle/internal/logging"
	"github.com/howzat/wordle/internal/words"
)

var CommitID string

func main() {
	log, err := logging.NewProductionLogger("admin-build-wordle-dictionary")
	if err != nil {
		panic(err)
	}

	log.Info("started ingestion",
		"commitId", CommitID,
		"environment", os.Getenv("ENVIRONMENT"),
	)

	ctx := context.Background()
	compiled, compileErr := words.CompileWordList(ctx, log, words.CompileConfig{
		Outfile:             "wordle-words.json",
		WordsetDataDir:      "dictionaries/wordset-dictionary/data",
		EnglishWordsDataDir: "dictionaries/english-words",
	})

	log.Info("complete", "ingested", compiled, "error", compileErr)
}
