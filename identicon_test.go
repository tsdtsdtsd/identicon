package identicon

import (
	"testing"
)

func TestNew(t *testing.T) {
	id := "test-string"

	ic, err := New(id)
	if err != nil {
		t.Error("Error creating new Identicon struct ", err)
	}

	if ic.ID != id {
		t.Error("ID error: expected", id, ", got ", ic.ID)
	}

	if ic.HashString() != "661f8009fa8e56a9d0e94a0a644397d7" {
		t.Error("MD5 error: expected 661f8009fa8e56a9d0e94a0a644397d7, got ", ic.HashString())
	}
}
