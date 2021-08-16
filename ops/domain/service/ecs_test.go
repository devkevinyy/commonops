package service

import (
	"testing"
)

func Test_ecs_BatchSshExec(t *testing.T) {
	e := &ecs{}
	e.BatchSshExec([]uint{1, 5, 6}, "hostname")
}
