package wordle

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func NewProductionLogger(namespace string) (*logr.Logger, error) {
	var log logr.Logger
	zapLog, err := zap.NewProduction()
	log = zapr.NewLogger(zapLog).WithName(namespace)
	return &log, err
}
