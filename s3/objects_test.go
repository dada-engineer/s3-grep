package s3

import "testing"

func TestNewObject(t *testing.T) {
	testObject := NewObject("this_is_a_test_object")
	if (
		testObject.Key != "this_is_a_test_object" &&
		testObject.Content != nil &&
		testObject.NumBytes != 0 &&
		testObject.Error != nil) {
		t.Error("Unexpected Object")
	}
}
