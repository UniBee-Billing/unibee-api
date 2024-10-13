package gateway

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"os"
	"strings"
	"unibee/utility"
)

func ReadEmailHtmlTemplate() (string, error) {
	data, err := os.ReadFile("./resource/email/template.html")
	if err != nil {
		return "", gerror.New(fmt.Sprintf("ReadEmailHtmlTemplate error:%s", err.Error()))
	}
	return string(data), nil
}

func ConvertToHtmlPage(content string) string {
	template, err := ReadEmailHtmlTemplate()
	utility.AssertError(err, "ReadEmailHtmlTemplate")
	template = strings.ReplaceAll(template, "{{content}}", content)
	return template
}

func ConvertPainToHtmlContent(content string) string {
	content = fmt.Sprintf("<p>%s</p>", content)
	content = strings.ReplaceAll(content, "\n", "</p>\n<p>")
	return content
}

func ConvertUniBeeTemplateToPlain(content string) string {
	content = strings.ReplaceAll(content, "&nbsp;", " ")
	content = strings.ReplaceAll(content, "<p>", "")
	content = strings.ReplaceAll(content, "</p>", "\n")
	return content
}
