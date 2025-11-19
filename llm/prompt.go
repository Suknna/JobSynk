package llm

import (
	_ "embed"
)

// 使用embed注入提示词文件
//
//go:embed prompt.txt
var Prompt string
