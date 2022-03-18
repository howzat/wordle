package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"

	"github.com/howzat/wordle/db"
	"github.com/howzat/wordle/internal/logging"
	"github.com/howzat/wordle/internal/words"
)

var CommitID string

const DictionaryBaseDirKey = "DICTIONARY_DIR"

type WordleConfig struct {
	BaseDir     string `env:"DICTIONARY_DIR"`
	Environment string `env:"ENVIRONMENT"`
}

func main() {

	log, err := logging.NewProductionLogger("admin-build-wordle-dictionary")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	var c WordleConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		panic(err)
	}

	if c.BaseDir == "" {
		panic("baseDir for reading dictionaries was empty,")
	}

	log.Info("started ingestion",
		"commitId", CommitID,
		"environment", c.Environment,
		"baseDir", c.BaseDir,
	)

	wordSource, err := words.NewWordSources(c.BaseDir)
	if err != nil {
		panic(err)
	}

	compiled, compileErr := wordSource.LoadWords(ctx, log)

	log.Info("complete", "ingested", compiled.Size, "error", compileErr)

	index, err := db.NewIndex(*log, compiled.Words, db.UseXXHashID)

	if err != nil {
		panic(err)
	}

	ws := db.NewSearchEngine(index)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")

	wordle := index.PickRandomWord()

	for {
		fmt.Print("guess> ")
		choice, _ := reader.ReadString('\n')
		choice = strings.Replace(choice, "\n", "", -1)

		s, err := db.NewWordleSearch(choice, db.BuildKnowledgeForGuess(wordle, choice))
		if err != nil {
			fmt.Println(err.Error())
		}

		match, err := ws.Search(*s)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("match result>guess %q\nmatch result>matches %q ", match.Guess, match.Items))
		fmt.Print("")
	}
}
