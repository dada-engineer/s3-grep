package s3

import "testing"

func TestNewObject(t *testing.T) {
	testObject := NewObject("this_is_a_test_object")
	if testObject.GetKey() != "this_is_a_test_object" {
		t.Error("Unexpected Object")
	}
}
