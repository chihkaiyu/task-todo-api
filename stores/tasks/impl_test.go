package tasks

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type taskSuite struct {
	suite.Suite
	taskStore *impl
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(taskSuite))
}
