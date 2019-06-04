package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dabdada/s3-grep/config"
	"github.com/dabdada/s3-grep/s3"
)

type testObject struct {
	Key      string
	Content  []byte
	NumBytes int64
	Error    error
}

func (o testObject) GetKey() string {
	return o.Key
}

func (o testObject) GetContent(session *config.AWSSession, bucketName string) ([]byte, int64, error) {
	return o.Content, o.NumBytes, o.Error
}

func newTestObject(key string, content []byte, numBytes int64, err error) s3.StoredObject {
	return testObject{Key: key, Content: content, NumBytes: numBytes, Error: err}
}

func TestPartitionS3Objects(t *testing.T) {
	partitionS3ObjectsTestData := []struct {
		name     string
		in       []s3.StoredObject
		num      int
		expected [][]s3.StoredObject
	}{
		{
			name:     "empty list",
			in:       []s3.StoredObject{},
			num:      1,
			expected: [][]s3.StoredObject{},
		},
		{
			name: "one list item divided into one partition",
			in:   []s3.StoredObject{newTestObject("test", []byte{}, 0, nil)},
			num:  1,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test", []byte{}, 0, nil)},
			},
		},
		{
			name: "one list item divided into two partitons",
			in:   []s3.StoredObject{newTestObject("test", []byte{}, 0, nil)},
			num:  2,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test", []byte{}, 0, nil)}, []s3.StoredObject{},
			},
		},
		{
			name: "two list items divided into one partition",
			in:   []s3.StoredObject{newTestObject("test", []byte{}, 0, nil), newTestObject("some", []byte{}, 0, nil)},
			num:  1,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{
					newTestObject("test", []byte{}, 0, nil),
					newTestObject("some", []byte{}, 0, nil),
				},
			},
		},
		{
			name: "two list items divided into two partitions",
			in:   []s3.StoredObject{newTestObject("test", []byte{}, 0, nil), newTestObject("some", []byte{}, 0, nil)},
			num:  2,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test", []byte{}, 0, nil)},
				[]s3.StoredObject{newTestObject("some", []byte{}, 0, nil)},
			},
		},
		{
			name: "two list items divided into three partitions",
			in:   []s3.StoredObject{newTestObject("test", []byte{}, 0, nil), newTestObject("some", []byte{}, 0, nil)},
			num:  3,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test", []byte{}, 0, nil)},
				[]s3.StoredObject{newTestObject("some", []byte{}, 0, nil)},
				[]s3.StoredObject{},
			},
		},
		{
			name: "three list items divided into one partition",
			in: []s3.StoredObject{
				newTestObject("test", []byte{}, 0, nil),
				newTestObject("some", []byte{}, 0, nil),
				newTestObject("objects", []byte{}, 0, nil),
			},
			num: 1,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{
					newTestObject("test", []byte{}, 0, nil),
					newTestObject("some", []byte{}, 0, nil),
					newTestObject("objects", []byte{}, 0, nil),
				},
			},
		},
		{
			name: "three list items divided into two partitions",
			in: []s3.StoredObject{
				newTestObject("test", []byte{}, 0, nil),
				newTestObject("some", []byte{}, 0, nil),
				newTestObject("objects", []byte{}, 0, nil),
			},
			num: 2,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test", []byte{}, 0, nil), newTestObject("some", []byte{}, 0, nil)},
				[]s3.StoredObject{newTestObject("objects", []byte{}, 0, nil)},
			},
		},
		{
			name: "three list items divided into three partitions",
			in: []s3.StoredObject{
				newTestObject("test", []byte{}, 0, nil),
				newTestObject("some", []byte{}, 0, nil),
				newTestObject("objects", []byte{}, 0, nil),
			},
			num: 3,
			expected: [][]s3.StoredObject{
				[]s3.StoredObject{newTestObject("test", []byte{}, 0, nil)},
				[]s3.StoredObject{newTestObject("some", []byte{}, 0, nil)},
				[]s3.StoredObject{newTestObject("objects", []byte{}, 0, nil)},
			},
		},
	}

	for _, tt := range partitionS3ObjectsTestData {
		t.Run(tt.name, func(t *testing.T) {

			actual := partitionS3Objects(tt.in, tt.num)

			for i := range actual {
				for j := 0; j < len(tt.expected[i]); j++ {
					if tt.expected[i][j].GetKey() != actual[i][j].GetKey() {
						t.Errorf(
							"expected[%d][%d] key: %s does not equal actual[%d][%d] key: %s",
							i, j, tt.expected[i][j].GetKey(), i, j, actual[i][j].GetKey())
					}
				}
			}
		})
	}
}

func TestGetContentExcerpt(t *testing.T) {
	testData := []struct {
		name     string
		text     []byte
		query    []byte
		expected []byte
	}{
		{"starts with query", []byte("someThing"), []byte("some"), []byte("someThing")},
		{
			"query in the middle but not enough chars before",
			[]byte("someThing"),
			[]byte("Thing"),
			[]byte("someThing"),
		},
		{
			"query in the middle not enough chars to the left and right",
			[]byte("someThing"),
			[]byte("meT"),
			[]byte("someThing"),
		},
		{
			"more than enough chars right and left of the query",
			[]byte("someThingSuperLongAndWeirdOnlyForTesting"),
			[]byte("Long"),
			[]byte("ThingSuperLongAndWeirdOn"),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {

			actual := getContentExcerpt(tt.text, tt.query)

			if !bytes.Equal(tt.expected, actual) {
				t.Errorf("expected excerpt is '%s' but actual was %s", tt.expected, actual)
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

func TestGrepInObjectContent(t *testing.T) {
	results := make(chan grepResult)
	done := make(chan int)
	testSession, err := config.NewAWSSession("testing")

	if err != nil {
		t.Error("Could not create test aws session")
		return
	}

	input := []s3.StoredObject{
		newTestObject("key0", []byte("This is a test containing the word: Blueberrycheescake"), 55, nil),
		newTestObject("key1", []byte("This is a test not containing the word."), 40, nil),
		newTestObject("key2", []byte{}, 0, nil),
		newTestObject("key3", []byte("This is a test containing the word, but raising an error: Blueberrycheesecake"), 59, errors.New("This is some error")),
	}

	go grepInObjectContent(testSession, "some-bucket", input, "berry", false, results, done)

	finished := 0

	for {
		select {
		case result := <-results:
			if result.Key != "key0" {
				t.Errorf("Key0 was expected, but is %s", result.Key)
			}
		case i := <-done:
			finished += i
		default:
			if finished == 1 {
				close(results)
				close(done)
				return
			}
		}
	}
}
