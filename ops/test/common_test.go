package test

import (
	"fmt"
	"testing"
	"time"
)

func TestInitK8sClient(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	tt, err := time.Parse(layout, "2021-01-07 09:19:45")
	fmt.Println(tt, err)
}
