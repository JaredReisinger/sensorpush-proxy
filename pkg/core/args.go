package core

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrRequiredArgMissing = errors.New("required arg missing")
)

func GetArgOrFile(argValue string, argFile string, required bool) (string, error) {
	// slog.Debug("GetArgOrFile", "argValue", argValue, "argFile", argFile, "required", required)
	if argValue != "" {
		return argValue, nil
	}

	if argFile != "" {
		b, err := os.ReadFile(argFile)
		if err != nil {
			err = fmt.Errorf("unable to read from %q: %w", argFile, err)
			return "", err
		}
		s := strings.TrimSpace(string(b))
		return s, nil
	}

	if required {
		return "", ErrRequiredArgMissing
	}

	return "", nil
}
