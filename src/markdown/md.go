package markdown

import (
	"html"
	"regexp"
	"strings"
	"syscall/js"
)

var (
	reEscapes       = regexp.MustCompile(`\\([\\*_\#\[\]\(\)\~\-\.\!>])`)
	reCodeInline    = regexp.MustCompile("`([^`]+)`")
	reBold          = regexp.MustCompile(`\*\*([^*]+)\*\*|__([^_]+)__`)
	reItalic        = regexp.MustCompile(`\*([^*]+)\*|_([^_]+)_`)
	reStrikethrough = regexp.MustCompile(`~~([^~]+)~~`)
	reImage         = regexp.MustCompile(`!\[([^\]]+)\]\(([^)]+)\)`)
	reLink          = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reHeading       = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	reBlockquote    = regexp.MustCompile(`^(>+)\s?(.*)$`)
	reHr            = regexp.MustCompile(`^(?:---|\*\*\*|___)\s*$`)
	reUl            = regexp.MustCompile(`^[\*\-]\s+(.+)$`)
	reOl            = regexp.MustCompile(`^\d+\.\s+(.+)$`)
	reTask          = regexp.MustCompile(`^[\*\-]\s+\[([ xX])\]\s+(.+)$`)
)

// Protegim els caràcters escapats amb marcadors temporals
func escapeEscapes(text string) string {
	return reEscapes.ReplaceAllStringFunc(text, func(m string) string {
		char := m[1]
		switch char {
		case '\\': return "__ESC_BS__"
		case '*': return "__ESC_STAR__"
		case '_': return "__ESC_UNDER__"
		case '#': return "__ESC_HASH__"
		case '[': return "__ESC_LB__"
		case ']': return "__ESC_RB__"
		case '(': return "__ESC_LP__"
		case ')': return "__ESC_RP__"
		case '`': return "__ESC_BT__"
		case '~': return "__ESC_TILDE__"
		case '-': return "__ESC_HYPHEN__"
		case '.': return "__ESC_DOT__"
		case '!': return "__ESC_EXCL__"
		case '>': return "__ESC_GT__"
		}
		return m
	})
}

// Restaurem els caràcters originals al final del procés
func restoreEscapes(text string) string {
	text = strings.ReplaceAll(text, "__ESC_BS__", `\`)
	text = strings.ReplaceAll(text, "__ESC_STAR__", "*")
	text = strings.ReplaceAll(text, "__ESC_UNDER__", "_")
	text = strings.ReplaceAll(text, "__ESC_HASH__", "#")
	text = strings.ReplaceAll(text, "__ESC_LB__", "[")
	text = strings.ReplaceAll(text, "__ESC_RB__", "]")
	text = strings.ReplaceAll(text, "__ESC_LP__", "(")
	text = strings.ReplaceAll(text, "__ESC_RP__", ")")
	text = strings.ReplaceAll(text, "__ESC_BT__", "`")
	text = strings.ReplaceAll(text, "__ESC_TILDE__", "~")
	text = strings.ReplaceAll(text, "__ESC_HYPHEN__", "-")
	text = strings.ReplaceAll(text, "__ESC_DOT__", ".")
	text = strings.ReplaceAll(text, "__ESC_EXCL__", "!")
	text = strings.ReplaceAll(text, "__ESC_GT__", ">")
	return text
}

// Il·luminador de sintaxi lleuger i natiu per a blocs de codi
func highlightCode(code string) string {
	safe := html.EscapeString(code)

	reComment := regexp.MustCompile(`(//.*|#.*)`)
	safe = reComment.ReplaceAllString(safe, `<span style="color: #6a737d; font-style: italic;">$1</span>`)

	reString := regexp.MustCompile(`(&quot;[^&]*&quot;|&#39;[^&#]*&#39;)`)
	safe = reString.ReplaceAllString(safe, `<span style="color: #032f62;">$1</span>`)

	keywords := []string{
		"func", "return", "var", "const", "if", "else", "for", "range", 
		"package", "import", "type", "struct", "true", "false", "nil", 
		"int", "string", "bool", "error", "interface", "map", "chan", 
		"go", "defer", "switch", "case", "default", "break", "continue",
		"class", "def", "async", "await", "let", "function",
	}

	for _, kw := range keywords {
		reKw := regexp.MustCompile(`\b` + kw + `\b`)
		safe = reKw.ReplaceAllString(safe, `<span style="color: #d73a49; font-weight: bold;">`+kw+`</span>`)
	}

	return safe
}

func parseInline(text string) string {
	text = reImage.ReplaceAllStringFunc(text, func(m string) string {
		submatches := reImage.FindStringSubmatch(m)
		if len(submatches) == 3 {
			alt := html.EscapeString(submatches[1])
			url := html.EscapeString(submatches[2])
			return `<img src="` + url + `" alt="` + alt + `" style="max-width: 100%;">`
		}
		return m
	})

	text = reLink.ReplaceAllStringFunc(text, func(m string) string {
		submatches := reLink.FindStringSubmatch(m)
		if len(submatches) == 3 {
			linkText := submatches[1]
			linkURL := submatches[2]
			return `<a href="` + html.EscapeString(linkURL) + `">` + linkText + `</a>`
		}
		return m
	})

	text = html.EscapeString(text)

	text = reCodeInline.ReplaceAllStringFunc(text, func(m string) string {
		inner := m[1 : len(m)-1]
		return "<code>" + inner + "</code>"
	})

	text = reBold.ReplaceAllStringFunc(text, func(m string) string {
		var inner string
		if len(m) >= 4 {
			inner = m[2 : len(m)-2]
		}
		return "<strong>" + inner + "</strong>"
	})

	text = reStrikethrough.ReplaceAllStringFunc(text, func(m string) string {
		inner := m[2 : len(m)-2]
		return "<del>" + inner + "</del>"
	})

	text = reItalic.ReplaceAllStringFunc(text, func(m string) string {
		var inner string
		if len(m) >= 2 {
			inner = m[1 : len(m)-1]
		}
		return "<em>" + inner + "</em>"
	})

	return text
}

func parseTableRow(line string, cellTag string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	var cells []string
	for _, p := range parts {
		cellContent := parseInline(strings.TrimSpace(p))
		cells = append(cells, "<"+cellTag+">"+cellContent+"</"+cellTag+">")
	}
	return strings.Join(cells, "\n")
}

func parseTable(lines []string) string {
	if len(lines) < 2 {
		return strings.Join(lines, "<br>")
	}
	var sb strings.Builder
	sb.WriteString(`<table style="border-collapse: collapse; width: 100%;">` + "\n")

	headerCells := parseTableRow(lines[0], "th")
	sb.WriteString("<thead>\n<tr style=\"border-bottom: 2px solid #ccc;\">\n" + headerCells + "\n</tr>\n</thead>\n")

	startIdx := 2
	if !strings.Contains(lines[1], "---") {
		startIdx = 1
	}

	sb.WriteString("<tbody>\n")
	for i := startIdx; i < len(lines); i++ {
		rowCells := parseTableRow(lines[i], "td")
		sb.WriteString("<tr style=\"border-bottom: 1px solid #eee;\">\n" + rowCells + "\n</tr>\n")
	}
	sb.WriteString("</tbody>\n</table>")
	return sb.String()
}

func ToHTML(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Error: Missing text argument"
	}

	md := args[0].String()
	// Protegim els caràcters escapats abans de processar res
	md = escapeEscapes(md)

	lines := strings.Split(md, "\n")
	var result []string

	inParagraph := false
	inCodeBlock := false
	inUl := false
	inOl := false
	var codeBuffer []string
	var codeLang string
	var tableLines []string

	closeAllLists := func() {
		if inUl {
			result = append(result, "</ul>")
			inUl = false
		}
		if inOl {
			result = append(result, "</ol>")
			inOl = false
		}
	}

	closeParagraph := func() {
		if inParagraph {
			result = append(result, "</p>")
			inParagraph = false
		}
	}

	flushTable := func() {
		if len(tableLines) > 0 {
			closeParagraph()
			closeAllLists()
			result = append(result, parseTable(tableLines))
			tableLines = nil
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			closeParagraph()
			closeAllLists()
			flushTable()
			if inCodeBlock {
				fullCode := strings.Join(codeBuffer, "\n")
				highlighted := highlightCode(fullCode)
				if codeLang != "" {
					result = append(result, `<pre><code class="language-`+html.EscapeString(codeLang)+`">`+highlighted+`</code></pre>`)
				} else {
					result = append(result, "<pre><code>"+highlighted+"</code></pre>")
				}
				inCodeBlock = false
				codeBuffer = nil
			} else {
				inCodeBlock = true
				codeLang = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
				codeBuffer = nil
			}
			continue
		}

		if inCodeBlock {
			codeBuffer = append(codeBuffer, line)
			continue
		}

		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			tableLines = append(tableLines, trimmed)
			continue
		} else {
			flushTable()
		}

		if trimmed == "" {
			closeParagraph()
			closeAllLists()
			continue
		}

		if reHr.MatchString(trimmed) {
			closeParagraph()
			closeAllLists()
			result = append(result, "<hr>")
			continue
		}

		if matches := reHeading.FindStringSubmatch(trimmed); matches != nil {
			closeParagraph()
			closeAllLists()
			level := len(matches[1])
			content := parseInline(matches[2])
			tag := "h" + string(rune('0'+level))
			result = append(result, "<"+tag+">"+content+"</"+tag+">")
			continue
		}

		if matches := reBlockquote.FindStringSubmatch(trimmed); matches != nil {
			closeParagraph()
			closeAllLists()
			depth := len(matches[1])
			content := parseInline(matches[2])

			var sb strings.Builder
			for i := 0; i < depth; i++ {
				sb.WriteString("<blockquote>")
			}
			sb.WriteString(content)
			for i := 0; i < depth; i++ {
				sb.WriteString("</blockquote>")
			}
			result = append(result, sb.String())
			continue
		}

		if matches := reTask.FindStringSubmatch(trimmed); matches != nil {
			closeParagraph()
			if inOl {
				result = append(result, "</ol>")
				inOl = false
			}
			if !inUl {
				result = append(result, `<ul style="list-style-type: none; padding-left: 0;">`)
				inUl = true
			}
			checked := strings.ToLower(matches[1]) == "x"
			content := parseInline(matches[2])

			checkboxAttr := `type="checkbox" disabled`
			if checked {
				checkboxAttr += ` checked`
			}
			result = append(result, `<li><input `+checkboxAttr+`> `+content+`</li>`)
			continue
		}

		if matches := reUl.FindStringSubmatch(trimmed); matches != nil {
			closeParagraph()
			if inOl {
				result = append(result, "</ol>")
				inOl = false
			}
			if !inUl {
				result = append(result, "<ul>")
				inUl = true
			}
			content := parseInline(matches[1])
			result = append(result, "<li>"+content+"</li>")
			continue
		}

		if matches := reOl.FindStringSubmatch(trimmed); matches != nil {
			closeParagraph()
			if inUl {
				result = append(result, "</ul>")
				inUl = false
			}
			if !inOl {
				result = append(result, "<ol>")
				inOl = true
			}
			content := parseInline(matches[1])
			result = append(result, "<li>"+content+"</li>")
			continue
		}

		closeAllLists()
		parsedLine := parseInline(line)
		if !inParagraph {
			result = append(result, "<p>"+parsedLine)
			inParagraph = true
		} else {
			result = append(result, "<br>"+parsedLine)
		}
	}

	if inCodeBlock {
		fullCode := strings.Join(codeBuffer, "\n")
		result = append(result, "<pre><code>"+highlightCode(fullCode)+"</code></pre>")
	}
	flushTable()
	closeParagraph()
	closeAllLists()

	// Unim el resultat final i restablim els caràcters escapats originals
	finalHTML := strings.Join(result, "\n")
	return restoreEscapes(finalHTML)
}