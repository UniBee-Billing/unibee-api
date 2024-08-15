package email

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"os"
	"strings"
	"unibee/utility"
)

func readEmailHtmlTemplate() (string, error) {
	data, err := os.ReadFile("./email/template.html")
	if err != nil {
		return "", gerror.New(fmt.Sprintf("readEmailHtmlTemplate error:%s", err.Error()))
	}
	return string(data), nil
}

func ContentToHtmlPage(content string) string {
	template, err := readEmailHtmlTemplate()
	utility.AssertError(err, "ReadEmailHtmlTemplate")
	template = strings.ReplaceAll(template, "{{content}}", content)
	return template
}

func ConvertPainToHtmlContent(content string) string {
	content = fmt.Sprintf("<p>%s</p>", content)
	content = strings.ReplaceAll(content, "\n", "</p><p>")
	return content
}

func ConvertUniBeeTemplateToPlain(content string) string {
	content = strings.ReplaceAll(content, "&nbsp;", " ")
	content = strings.ReplaceAll(content, "<p>", "")
	content = strings.ReplaceAll(content, "</p>", "\n")
	return content
}
