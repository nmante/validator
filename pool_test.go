package validator

import (
	"testing"
)

func TestSetNumWorkers(t *testing.T) {
	numWorkers := 20
	p, _ := NewWorkerPool(numWorkers, []Job{})

	if p.numWorkers != numWorkers {
		t.Errorf("Num workers should be %d", numWorkers)
	}

	p.setNumWorkers(100000000)

	if p.numWorkers != MaxWorkers {
		t.Errorf("Num workers should be %d", MaxWorkers)
	}

	p.setNumWorkers(-1)

	if p.numWorkers != MinWorkers {
		t.Errorf("Num workers should be %d", MinWorkers)
	}
}

func TestNewWorkerPool(t *testing.T) {
	numWorkers := 20
	p, err := NewWorkerPool(numWorkers, nil)

	if err != nil {
		t.Errorf("There should be no error. %s", err.Error())
	}

	if p.numWorkers != numWorkers {
		t.Errorf("There should %d workers in the pool.", numWorkers)
	}
}
