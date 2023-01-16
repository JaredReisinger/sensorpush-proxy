//go:build tools
// +build tools

package tools

import (
	_ "github.com/automation-co/husky"
	// See the comment in Taskfile.yml's "prepare" task.  We can't use task to
	// acquire/setup task itself. :sadpanda:
	// _ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/lintingzhen/commitizen-go"
)
