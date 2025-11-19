package pdf

import (
	"bytes"
	"strings"

	"github.com/ledongthuc/pdf"
)

// 读取pdf到纯文本文件
func ToString(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		panic(err)
	}
	buf.ReadFrom(b)
	content := buf.String()
	// 将回车换行符统一转换为 \n
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	return content, nil
}
