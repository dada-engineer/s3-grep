package cli

import "testing"

func TestPartitionS3Objects(t *testing.T) {
	partitionS3ObjectsTestData := []struct {
		name     string
		in       []string
		num      int
		expected [][]string
	}{
		{
			"empty list",
			[]string{}, 1, [][]string{},
		},
		{
			"one list item divided into one partition",
			[]string{"test"}, 1, [][]string{[]string{"test"}},
		},
		{
			"one list item divided into two partitons",
			[]string{"test"}, 2, [][]string{[]string{"test"}, []string{}},
		},
		{
			"two list items divided into one partition",
			[]string{"test", "some"}, 1, [][]string{[]string{"test", "some"}},
		},
		{
			"two list items divided into two partitions",
			[]string{"test", "some"}, 2, [][]string{[]string{"test"}, []string{"some"}},
		},
		{
			"two list items divided into three partitions",
			[]string{"test", "some"}, 3, [][]string{[]string{"test"}, []string{"some"}, []string{}},
		},
		{
			"three list items divided into one partition",
			[]string{"test", "some", "strings"}, 1, [][]string{[]string{"test", "some", "strings"}},
		},
		{
			"three list items divided into two partitions",
			[]string{"test", "some", "strings"}, 2, [][]string{[]string{"test", "some"}, []string{"strings"}},
		},
		{
			"three list items divided into three partitions",
			[]string{"test", "some", "strings"}, 3, [][]string{[]string{"test"}, []string{"some"}, []string{"strings"}},
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
