package jsonsock

import (
	"testing"
)

import (
	"github.com/gando999/jsonsock"
)

func TestFillStruct(t *testing.T) {
	type R struct {
		StringField string
		IntField    int
	}
	r := new(R)
	m := make(map[string]interface{})
	m["StringField"] = "StringValue"
	m["IntField"] = 25

	jsonsock.FillStruct(r, m)
	if r.StringField != "StringValue" {
		t.Errorf("Did not set StringValue correctly, got %s expected %s", r.StringField, "StringValue")
	}
	if r.IntField != 25 {
		t.Errorf("Did not set IntValue correctly, got %d expected %d", r.IntField, 25)
	}

	type C struct {
		RField      R
		Description string
	}

	c := new(C)
	n := make(map[string]interface{})
	n["RField"] = *r
	n["Description"] = "TestMessage"

	jsonsock.FillStruct(c, n)
	if c.Description != "TestMessage" {
		t.Errorf("Did not set Description correctly, got %s expected %s", c.Description, "TestMessage")
	}
	if c.RField != *r {
		t.Errorf("Did not set RField correctly, got %v, expected %v", c.RField, r)
	}
	if c.RField.StringField != "StringValue" {
		t.Errorf("Did not set StringValue correctly, got %s expected %s", c.RField.StringField, "StringValue")
	}
	if c.RField.IntField != 25 {
		t.Errorf("Did not set IntValue correctly, got %d expected %d", c.RField.IntField, 25)
	}
}
