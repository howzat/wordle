package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var CommitID string
var zlog *zap.Logger

func NewProductionLogger(namespace string) (logr.Logger, error) {

	var log logr.Logger

	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}

	log = zapr.NewLogger(zapLog).WithName(namespace)

	return log, err
}

var log logr.Logger

func init() {
	log, err := NewProductionLogger("admin-build-dictionary")
	if err != nil {
		panic(err)
	}

	log.Info("lambda initialised",
		zap.String("commitId", CommitID),
		zap.String("environment", os.Getenv("ENVIRONMENT")))
}

type IngestResponse struct{}

func main() {
	zlog.Info("lambda started")
	lambda.Start(func(context.Context) (*IngestResponse, error) {
		zlog.Info("lambda called")
		return nil, nil
	})
}
