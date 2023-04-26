package lib

import (
	"fmt"
	"html/template"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 //:]+`)

var unprecString []string

// {"\x12", "\x01", ":\x01*", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07", "\x08", "\x09"}

func replAscii(p string) string {
	b := p
	for _, v := range unprecString {
		b = strings.ReplaceAll(b, v, "")
	}

	return nonAlphanumericRegex.ReplaceAllLiteralString(b, "")
}

func combineUnprec() {
	for i := 0; i < 50; i++ {
		x := strconv.Itoa(i)
		if i < 10 {
			x = fmt.Sprintf("0%s", x)
		}
		comb := fmt.Sprintf(`\x%s`, x)
		unprecString = append(unprecString, comb)
	}
	unprecString = append(unprecString, "")
}

func ucFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func unescape(s string) template.HTML {
	return template.HTML(s)
}

func getFirstService(param []Service) Service {
	if len(param) == 0 {
		return Service{}
	}
	return param[0]
}

func strReplaceParam(param string) string {
	return strings.Replace(param, "{id}", "1", -1)
}

func protoFileBaseName(name string) string {
	if ext := path.Ext(name); ext == ".proto" {
		name = name[:len(name)-len(ext)]
	}
	return name
}

func getStringFromOptCode(param string) string {
	switch param {
	case "50056":
		return "httpMode"
	case "72295728":
		return "urlPath"
	case "50060":
		return "requiredField"
	}
	return ""
}

func cleanQuote(param string) string {
	return strings.Replace(param, `"`, "", -1)
}

func stringToOpt(param string) []*Option {
	var newOptions []*Option

	build := param
	countChar := strings.Count(param, ":")
	for i := 1; i <= countChar; i++ {
		idx := strings.Index(build, ":")
		if idx > -1 {
			if build[idx+1] == ' ' && (i%2) == 0 {
				build = build[:idx] + build[idx+1:]
			} else if build[idx+1] == ' ' {
				build = build[:idx] + "=" + build[idx+2:]
			} else {
				build = build[:idx] + "=" + build[idx+1:]
			}

		}
	}
	onOpt := strings.Split(build, " ")

	for _, vOnOpt := range onOpt {
		splitV := strings.Split(vOnOpt, "=")
		if len(splitV) < 2 {
			continue
		}

		newOptions = append(newOptions, &Option{
			Code:  splitV[0],
			Name:  splitV[0],
			Value: cleanQuote(splitV[1]),
		})
	}
	return newOptions
}

func stringExist(param string, slice []string) bool {
	for _, v := range slice {
		if strings.ToLower(v) == strings.ToLower(param) {
			return true
		}
	}
	return false
}

func getHttpUrl(param string) string {
	i := strings.Index(param, `/`)

	if i == -1 {
		return ""
	}

	chars := param[i:]
	return chars

}

func ucDown(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func grpcTypeToTs(param string) string {
	switch param {
	case "TYPE_STRING":
		return "string"
	case "TYPE_INT64":
		return "number"
	case "TYPE_INT32":
		return "number"
	case "TYPE_BOOL":
		return "boolean"
	default:
		return "string"
	}
}

func findMessage(name string, param []Message) (Message, bool) {
	for _, v := range param {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			newMessage := Message(v)
			return newMessage, true
		}
	}
	return Message{}, false
}

func getMessageByName(name string, messages []Message) Message {
	for k, v := range messages {
		if ucFirst(name) == ucFirst(v.Name) {
			return messages[k]
		}
	}
	return Message{}
}

func genEmptyEnum() []*Enum {
	return []*Enum{}
}

func genEmptyMessage() []Message {
	return []Message{}
}

func genLocation(p string, s ...string) string {
	if len(s) == 0 {
		return ""
	}

	res := strings.Count(p, "%")
	genString := ""
	for i := 0; i < res; i++ {
		if i == 0 {
			genString = fmt.Sprintf("%s", s[i])
			continue
		}
		if len(s) > (i + 1) {
			genString = fmt.Sprintf("%s/%s", genString, s[len(s)-1])
			continue
		}
		genString = fmt.Sprintf("%s/%s", genString, s[i])
	}
	genString += "/"
	return genString
}
