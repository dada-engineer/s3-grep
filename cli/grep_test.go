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

func TestGetContentExcerpt(t *testing.T) {
	testData := []struct {
		name     string
		text     []byte
		query    []byte
		expected []byte
	}{
		{"starts with query", []byte("someThing"), []byte("some"), []byte("someThing")},
		{
			"query in the middle but not more than MAX_EXCERPT_LENGTH chars in text",
			[]byte("Bounty tackle nipper red ensign execution dock Sail ho spirits hail-shot scourge"),
			[]byte("dock"),
			[]byte("Bounty tackle nipper red ensign execution dock Sail ho spirits hail-shot scourge"),
		},
		{
			"query in the middle not enough chars to the left",
			[]byte("Bounty tackle nipper red ensign execution dock Sail ho spirits hail-shot scourge of the seven seas barkadeer booty keel hands provost loaded to the gunwalls"),
			[]byte("nipper"),
			[]byte("Bounty tackle nipper red ensign execution dock Sail ho spirits hail-shot scourge"),
		},
		{
			"query in the middle not enough chars to the right",
			[]byte("Bounty tackle nipper red ensign execution dock Sail ho spirits hail-shot scourge of the seven seas barkadeer booty keel hands provost loaded to the gunwalls"),
			[]byte("barkadeer"),
			[]byte("Sail ho spirits hail-shot scourge of the seven seas barkadeer booty keel hands provost loaded to the gunwalls"),
		},
		{
			"more than enough chars right and left of the query",
			[]byte("Bounty tackle nipper red ensign execution dock Sail ho spirits hail-shot scourge of the seven seas barkadeer booty keel hands provost loaded to the gunwalls"),
			[]byte("shot"),
			[]byte("nipper red ensign execution dock Sail ho spirits hail-shot scourge of the seven seas barkadeer booty keel hands provost"),
		},
		{
			"more than enough chars right and left of the query, find a space in from index",
			[]byte("Bounty tackle nipper red ensign execution dock Sail ho spirits hail-shot scourge of the seven seas barkadeer booty keel hands provost loaded to the gunwalls"),
			[]byte("even"),
			[]byte("execution dock Sail ho spirits hail-shot scourge of the seven seas barkadeer booty keel hands provost loaded to the gunwalls"),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {

			actual := getContentExcerpt(tt.text, tt.query)

			if !bytes.Equal(tt.expected, actual) {
				t.Errorf("expected excerpt is '%s' but actual was '%s'", tt.expected, actual)
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
	results := make(chan *grepResult)
	done := make(chan bool)
	testSession, err := config.NewAWSSession("testing")

	if err != nil {
		t.Error("Could not create test aws session")
		return
	}

	input := []s3.StoredObject{
		newTestObject("key0", []byte("This is a test containing the word: Blueberrycheescake"), 54, nil),
		newTestObject("key1", []byte("This is a test not containing the word."), 39, nil),
		newTestObject("key2", []byte{}, 0, nil),
		newTestObject("key3", []byte("This is a test containing the word, but raising an error: Blueberrycheesecake"), 77, errors.New("This is some error")),
	}

	for _, object := range input {
		go grepInObjectContent(testSession, "some-bucket", object, "berry", false, results, done)
	}

	finished := 0

	for {
		select {
		case result := <-results:
			if result.Key != "key0" {
				t.Errorf("Key0 was expected, but is %s", result.Key)
			}
		case <-done:
			finished++
		default:
			if finished == len(input) {
				close(results)
				close(done)
				return
			}
		}
	}
}
