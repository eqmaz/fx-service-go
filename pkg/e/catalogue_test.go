package e

import "testing"

func TestSetCatalogue(t *testing.T) {
	c := ErrorMap{
		"e12345": "This is an example error",
	}
	SetCatalogue(c)

	if catalogue["e12345"] != "This is an example error" {
		t.Errorf("expected 'This is an example error', got %s", catalogue["e12345"])
	}
}

func TestThrowErrorFromCatalogue(t *testing.T) {
	c := ErrorMap{
		"e12345": "This is an example error",
	}
	SetCatalogue(c)

	ex := FromCode("e12345")
	if ex.GetMessage() != "This is an example error" {
		t.Errorf("expected 'This is an example error', got %s", ex.GetMessage())
	}
}
