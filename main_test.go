package heaper

import (
	"testing"
	"time"
)

func TestProcess_Read(t *testing.T) {
	go Run(1, 3)
	time.Sleep(5 * time.Second)
	result := Read()

	if result[0].Time > result[1].Time {
		t.Fail()
	}
}

func TestProcess_Read2(t *testing.T) {
	go Run(1, 3)
	time.Sleep(2 * time.Second)
	result := Read()

	if result[0].Time == 0 {
		t.Fail()
	}
}
