package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/howzat/wordle"
	"github.com/howzat/wordle/internal/wordgen"
	"github.com/pkg/errors"
)

var CommitID string

func main() {

	ctx := context.Background()

	log, err := wordle.NewProductionLogger("admin-build-wordle-dictionary")
	failOnErr(err)

	config, err := wordgen.NewDictionaryConfig(ctx)
	failOnErr(err)

	if config.BaseDir == "" {
		err = errors.New("no location provided for dictionary directory")
		log.Error(err, "%v was empty", wordgen.DictionaryBaseDirKey)
		failOnErr(err)
	}

	log.Info("started ingestion",
		"commitId", CommitID,
		"baseDir", config.BaseDir,
	)

	dictionaryConfig, err := wordgen.NewDictionaryConfig(ctx)
	failOnErr(err)

	wordSource, err := wordgen.NewWordSources(dictionaryConfig)
	failOnErr(err)

	compiled, compileErr := wordSource.LoadWords(ctx, log, wordgen.LowercaseWord, wordgen.WordleCandidate)

	log.Info("complete", "ingested", compiled.Size, "error", compileErr)

	words := map[string]bool{}
	var uniqueWords []string
	for _, word := range compiled.Words {
		if words != nil {
			if _, present := words[word]; !present {
				words[word] = true
				uniqueWords = append(uniqueWords, word)
			}
		}
	}

	log.Info("optimised", "unique", len(uniqueWords))

	dictionaryFile, err := os.OpenFile("cmd/search/dictionary.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	failOnErr(err)

	defer func(f *os.File) {
		_ = f.Close()
	}(dictionaryFile)

	err = dictionaryFile.Truncate(0)
	failOnErr(err)

	_, err = dictionaryFile.Seek(0, 0)
	failOnErr(err)

	sort.Sort(sort.StringSlice(uniqueWords))

	log.Info(fmt.Sprintf("About to write %v lines", len(uniqueWords)))

	var builder strings.Builder
	for _, word := range uniqueWords {
		builder.WriteString(word + "\n")
	}

	_, err = dictionaryFile.WriteString(builder.String())
	failOnErr(err)
}

func failOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
