package markdown

import (
	"regexp"
	"strings"
	"syscall/js"
)

var (
	reCodeInline = regexp.MustCompile("`([^`]+)`")
	reBold       = regexp.MustCompile(`\*\*([^*]+)\*\*|__([^_]+)__`)
	reItalic     = regexp.MustCompile(`\*([^*]+)\*|_([^_]+)_`)
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reHeading    = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	reBlockquote = regexp.MustCompile(`^>\s?(.*)$`)
	reHr         = regexp.MustCompile(`^(?:---|\*\*\*|___)\s*$`)
)

func parseInline(text string) string {
	text = reCodeInline.ReplaceAllString(text, "<code>$1</code>")
	text = reBold.ReplaceAllString(text, "<strong>$1$2</strong>")
	text = reItalic.ReplaceAllString(text, "<em>$1$2</em>")
	text = reLink.ReplaceAllString(text, `<a href="$2">$1</a>`)
	return text
}

func ToHTML(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Error: Missing text argument"
	}

	md := args[0].String()
	lines := strings.Split(md, "\n")
	var result []string
	inParagraph := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if inParagraph {
				result = append(result, "</p>")
				inParagraph = false
			}
			continue
		}

		if reHr.MatchString(trimmed) {
			if inParagraph {
				result = append(result, "</p>")
				inParagraph = false
			}
			result = append(result, "<hr>")
			continue
		}

		if matches := reHeading.FindStringSubmatch(trimmed); matches != nil {
			if inParagraph {
				result = append(result, "</p>")
				inParagraph = false
			}
			level := len(matches[1])
			content := parseInline(matches[2])
			tag := "h" + string(rune('0'+level))
			result = append(result, "<"+tag+">"+content+"</"+tag+">")
			continue
		}

		if matches := reBlockquote.FindStringSubmatch(trimmed); matches != nil {
			if inParagraph {
				result = append(result, "</p>")
				inParagraph = false
			}
			content := parseInline(matches[1])
			result = append(result, "<blockquote>"+content+"</blockquote>")
			continue
		}

		parsedLine := parseInline(line)
		if !inParagraph {
			result = append(result, "<p>"+parsedLine)
			inParagraph = true
		} else {
			result = append(result, "<br>"+parsedLine)
		}
	}

	if inParagraph {
		result = append(result, "</p>")
	}

	return strings.Join(result, "\n")
}