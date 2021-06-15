package stock

import (
	"testing"
)

func TestMakeLine(t *testing.T) {
	var (
		in       = []string{"A", "B", "C"}
		expected = "| A | B | C | "
	)
	actual := MakeLine(in...)
	if actual != expected {
		t.Errorf("MakeLine(%v) = %s; expected %s", in, actual, expected)
	}
}

func TestMakeStockChanges(t *testing.T) {
	var (
		in       = 34.88
		expected = "+34.88%📈"
	)
	actual := MakeStockChanges(in)
	if actual != expected {
		t.Errorf("MakeStockChanges(%f) = %s; expected %s", in, actual, expected)
	}
}
