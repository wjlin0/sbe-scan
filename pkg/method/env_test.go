package method

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestEnvJson_UnmarshalJSON is a generated function returning the mock function.
func TestEnvJson_UnmarshalJSON(t *testing.T) {
	file, err := os.ReadFile("env.test.json")
	if err != nil {
		t.Error(err)
		return
	}
	e := Configuration{}
	err = json.Unmarshal([]byte(file), &e)
	if err != nil {
		t.Error()
	}

}

// TestConfiguration_GetProperties is a generated function returning the mock function.
func TestConfiguration_GetProperties(t *testing.T) {
	file, err := os.ReadFile("env.test.json")
	if err != nil {
		t.Error(err)
		return
	}
	e := Configuration{}
	err = json.Unmarshal([]byte(file), &e)
	if err != nil {
		t.Error(err)
	}
	properties := e.GetProperties()
	if len(properties) <= 0 {
		t.Error("properties is empty")
	}
	fmt.Println(properties)
}
