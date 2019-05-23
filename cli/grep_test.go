package cli

import (
	"testing"

	"github.com/dabdada/s3-grep/config"
	"github.com/dabdada/s3-grep/s3"
)

type testObject struct {
	Key string
}

func (o testObject) GetKey() string {
	return o.Key
}

func (o testObject) GetContent(session *config.AWSSession, bucketName string) ([]byte, int64, error) {
	return []byte{}, 0, nil
}

func newTestObject(key string) s3.StoredObject {
	return testObject{Key: key}
}

func TestPartitionS3Objects(t *testing.T) {
	partitionS3ObjectsTestData := []struct {
		name     string
		in       []s3.StoredObject
		num      int
		expected [][]s3.StoredObject
	}{
		{
			name: "empty list",
			in: []s3.StoredObject{},
			num: 1,
			expected: [][]s3.StoredObject{},
		},
		{
			name: "one list item divided into one partition",
			in: []s3.StoredObject{newTestObject("test")},
			num: 1,
			expected: [][]s3.StoredObject{[]s3.StoredObject{newTestObject("test")}},
		},
		{
			name: "one list item divided into two partitons",
			in: []s3.StoredObject{newTestObject("test")},
			num: 2,
			expected: [][]s3.StoredObject{[]s3.StoredObject{newTestObject("test")}, []s3.StoredObject{}},
		},
		{
			name: "two list items divided into one partition",
			in: []s3.StoredObject{newTestObject("test"), newTestObject("some")},
			num: 1,
			expected: [][]s3.StoredObject{[]s3.StoredObject{newTestObject("test"), newTestObject("some")}},
		},
		{
			name: "two list items divided into two partitions",
			in: []s3.StoredObject{newTestObject("test"), newTestObject("some")},
			num: 2,
			expected: [][]s3.StoredObject{[]s3.StoredObject{newTestObject("test")}, []s3.StoredObject{newTestObject("some")}},
		},
		{
			name: "two list items divided into three partitions",
			in: []s3.StoredObject{newTestObject("test"), newTestObject("some")},
			num: 3,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test")}, []s3.StoredObject{newTestObject("some")}, []s3.StoredObject{}},
		},
		{
			name: "three list items divided into one partition",
			in: []s3.StoredObject{newTestObject("test"), newTestObject("some"), newTestObject("objects")},
			num: 1,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test"), newTestObject("some"), newTestObject("objects")}},
		},
		{
			name: "three list items divided into two partitions",
			in: []s3.StoredObject{newTestObject("test"), newTestObject("some"), newTestObject("objects")},
			num: 2,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test"), newTestObject("some")}, []s3.StoredObject{newTestObject("objects")}},
		},
		{
			name: "three list items divided into three partitions",
			in: []s3.StoredObject{newTestObject("test"), newTestObject("some"), newTestObject("objects")},
			num: 3,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test")},
				[]s3.StoredObject{newTestObject("some")},
				[]s3.StoredObject{newTestObject("objects")}},
		},
	}

	for _, tt := range partitionS3ObjectsTestData {
		t.Run(tt.name, func(t *testing.T) {

			actual := partitionS3Objects(tt.in, tt.num)

			for i := range actual {
				for j := 0; j < len(tt.expected[i]); j++ {
					if tt.expected[i][j] != actual[i][j] {
						t.Errorf(
							"expected[%d][%d]: %s does not equal actual[%d][%d]: %s",
							i, j, tt.expected[i][j], i, j, actual[i][j])
					}
				}
			}
		})
	}
}

func TestCaseAwareContains(t *testing.T) {
	testData := []struct {
		name       string
		in         []byte
		sub        []byte
		ignoreCase bool
		expected   bool
	}{
		{"contains case sensitive", []byte("someThing"), []byte("Thin"), false, true},
		{"contains case insensitive", []byte("someThing"), []byte("Thin"), true, true},
		{"not contains case sensitive", []byte("someThing"), []byte("thin"), false, false},
		{"not contains case insensitive", []byte("someThing"), []byte("bar"), false, false},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {

			actual := caseAwareContains(tt.in, tt.sub, tt.ignoreCase)

			if tt.expected != actual {
				t.Errorf(
					"expected '%s' containing '%s' result is %t, actual was %t, while ignoreCase is %t",
					tt.in, tt.sub, tt.expected, actual, tt.ignoreCase)
			}
		})
	}
}
