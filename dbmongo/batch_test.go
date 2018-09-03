package main

import (
	"testing"
)

func Test_nextBatchID(t *testing.T) {
	batchID, err := nextBatchID("abcd")

	if err == nil {
		t.Error("Erreur attendue absente. '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test valeur erronée ok")
	}
	batchID, err = nextBatchID("1801")
	if err != nil || batchID != "1802" {
		t.Error("'1802' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test valeur courante ok")
	}
	batchID, err = nextBatchID("1812")
	if err != nil || batchID != "1901" {
		t.Error("'1901' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test passage nouvelle année ok")
	}
}
