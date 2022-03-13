package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/howzat/wordle/internal/logging"
	"github.com/howzat/wordle/internal/words"
	"github.com/howzat/wordle/search"
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
	files, err := words.NewWordSourceFiles("dictionaries")
	if err != nil {
		panic(err)
	}

	compiled, compileErr := words.CompileWordList(ctx, log, *files)

	log.Info("complete", "ingested", compiled.Ingested, "error", compileErr)

	db, err := search.NewIndexedDB(*log, compiled.Words, search.UseSeaHashID)
	if err != nil {
		panic(err)
	}

	ws := search.NewSearchEngine(db)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")

	for {
		fmt.Print("guess> ")
		choice, _ := reader.ReadString('\n')
		choice = strings.Replace(choice, "\n", "", -1)

		s, err := search.NewWordle(choice)
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
