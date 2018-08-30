package main

import (
	"testing"
)

func Test_nextBatchID(t *testing.T) {
	batchID, err := nextBatchID("1814")
	if err == nil {
		t.Error("Erreur attendue absente. '" + batchID + "' obtenu. err = " + err.Error())
	}
	batchID, err = nextBatchID("1801")
	if err != nil || batchID != "1802" {
		t.Error("'1802' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	}
	batchID, err = nextBatchID("1812")
	if err != nil || batchID != "1901" {
		t.Error("'1901' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	}
}
