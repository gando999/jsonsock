package jsonsock

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

import (
	"github.com/gando999/jsonsock"
)

func TestMarshalRequest(t *testing.T) {
	expectedIntArgs := `{"id":123,"jsonrpc":"2.0","method":"hello.world","params":[2,3]}`
	if resp, err := jsonsock.MarshalRequest(
		"hello.world",
		[]interface{}{2, 3}); err == nil {
		if resp != expectedIntArgs {
			t.Errorf("Failed, got %s expected %s", resp, expectedIntArgs)
		}
	}

	expectedStringArgs := `{"id":123,"jsonrpc":"2.0","method":"hello.world","params":["2","3"]}`
	if resp, err := jsonsock.MarshalRequest(
		"hello.world",
		[]interface{}{"2", "3"}); err == nil {
		if resp != expectedStringArgs {
			t.Errorf("Failed, got %s expected %s", resp, expectedStringArgs)
		}
	}

	expectedMixedArgs := `{"id":123,"jsonrpc":"2.0","method":"hello.world","params":[2,"3"]}`
	if resp, err := jsonsock.MarshalRequest(
		"hello.world",
		[]interface{}{2, "3"}); err == nil {
		if resp != expectedMixedArgs {
			t.Errorf("Failed, got %s expected %s", resp, expectedMixedArgs)
		}
	}

	expectedObjectArgs := `{"id":123,"jsonrpc":"2.0","method":"hello.world","params":[2,{"test":"me"}]}`
	if resp, err := jsonsock.MarshalRequest(
		"hello.world",
		[]interface{}{2, map[string]string{"test": "me"}}); err == nil {
		if resp != expectedObjectArgs {
			t.Errorf("Failed, got %s expected %s", resp, expectedObjectArgs)
		}
	}

	expectedMixedObjectArgs := `{"id":123,"jsonrpc":"2.0","method":"hello.world","params":[2,{"test":"me","test2":2}]}`
	if resp, err := jsonsock.MarshalRequest(
		"hello.world",
		[]interface{}{2, map[string]interface{}{"test": "me", "test2": 2}}); err == nil {
		if resp != expectedMixedObjectArgs {
			t.Errorf("Failed, got %s expected %s", resp, expectedMixedObjectArgs)
		}
	}
}

func TestUnmarshalRequest(t *testing.T) {
	request := `{"id":123,"jsonrpc":"2.0","method":"hello.world","params":[2,3]}`
	if resp, err := jsonsock.UnmarshalRequest(string(request)); err == nil {
		if resp.Id != 123 {
			t.Errorf("Failed id, got %s expected %s", resp.Id, 123)
		}
		if resp.Jsonrpc != "2.0" {
			t.Errorf("Failed JsonRpc, got %s expected %s", resp.Jsonrpc, "2.0")
		}
		if resp.Method != "hello.world" {
			t.Errorf("Failed Method, got %s expected %s", resp.Method, "hello.world")
		}
		if cmp.Equal(resp.Params, []interface{}{2, 3}) {
			t.Errorf("Failed params, got %s expected %s", resp.Params, []interface{}{2, 3})
		}
	}
}

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
