package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var (
// 	ErrRequiredArgMissing = errors.New("required arg missing")
// )

func TestGetArgOrFile(t *testing.T) {
	argFile := "testdata/argfile.txt"
	argFileMissing := "testdata/argfile-MISSING.txt"
	fileVal := "value with whitespace"

	for i, c := range []struct {
		val         string
		file        string
		required    bool
		expected    string
		expectedErr error
	}{
		{"", "", false, "", nil},
		{"yes", "", false, "yes", nil},

		{"", argFile, false, fileVal, nil},
		{"yes", argFile, false, "yes", nil},

		{"", argFileMissing, false, "", assert.AnError},
		{"yes", argFileMissing, false, "yes", nil},

		{"", "", true, "", ErrRequiredArgMissing},
		{"yes", "", true, "yes", nil},

		{"", argFile, true, fileVal, nil},
		{"yes", argFile, true, "yes", nil},

		{"", argFileMissing, true, "", assert.AnError},
		{"yes", argFileMissing, true, "yes", nil},
	} {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			t.Parallel()
			actual, actualErr := GetArgOrFile(c.val, c.file, c.required)
			assert.Equal(t, c.expected, actual)

			switch c.expectedErr {
			case nil:
				assert.Nil(t, actualErr)
			case assert.AnError:
				assert.Error(t, actualErr)
			default:
				assert.ErrorIs(t, actualErr, c.expectedErr)
			}
		})
	}
}
