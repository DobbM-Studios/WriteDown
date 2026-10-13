package dobbmmarkup

import (
	"strings"
	"syscall/js"
)

func parseInline(text string) string {
	var result strings.Builder
	runes := []rune(text)
	i := 0
	n := len(runes)

	for i < n {
		if (runes[i] == ':' || runes[i] == '|') && i > 0 {
			
			if runes[i] == ':' && i+2 < n && runes[i+1] == '/' && runes[i+2] == '/' {
				result.WriteRune(runes[i])
				i++
				continue
			}

			tagStart := i - 1
			for tagStart >= 0 && ((runes[tagStart] >= 'a' && runes[tagStart] <= 'z') || 
				(runes[tagStart] >= 'A' && runes[tagStart] <= 'Z') || 
				(runes[tagStart] >= '0' && runes[tagStart] <= '9')) {
				tagStart--
			}
			tagStart++

			if tagStart < i {
				tag := string(runes[tagStart:i])
				var attrs string
				colonIdx := -1

				if runes[i] == '|' {
					attrStart := i + 1
					for j := attrStart; j < n; j++ {
						if runes[j] == ':' {
							if j+2 < n && runes[j+1] == '/' && runes[j+2] == '/' {
								continue
							}
							colonIdx = j
							break
						}
					}
					if colonIdx != -1 {
						attrs = " " + strings.TrimSpace(string(runes[attrStart:colonIdx]))
					}
				} else {
					colonIdx = i
					attrs = ""
				}

				if colonIdx != -1 {
					currentStr := result.String()
					result.Reset()
					if len(currentStr) >= (i - tagStart) {
						result.WriteString(currentStr[:len(currentStr)-(i-tagStart)])
					}

					contentIdx := colonIdx + 1

					for contentIdx < n && (runes[contentIdx] == ' ' || runes[contentIdx] == '\t' || runes[contentIdx] == '\n' || runes[contentIdx] == '\r') {
						contentIdx++
					}

					if contentIdx < n && runes[contentIdx] == '[' {
						balance := 1
						closeIdx := -1
						for k := contentIdx + 1; k < n; k++ {
							if runes[k] == '[' {
								balance++
							} else if runes[k] == ']' {
								balance--
								if balance == 0 {
									closeIdx = k
									break
								}
							}
						}

						if closeIdx != -1 {
							innerContent := string(runes[contentIdx+1 : closeIdx])
							parsedInner := parseInline(innerContent)

							result.WriteByte('<')
							result.WriteString(tag)
							result.WriteString(attrs)
							result.WriteByte('>')
							result.WriteString(parsedInner)
							result.WriteString("</")
							result.WriteString(tag)
							result.WriteByte('>')
							i = closeIdx + 1
							continue
						}
					}

					endOfLine := contentIdx
					for endOfLine < n && runes[endOfLine] != '\n' && runes[endOfLine] != '\r' {
						endOfLine++
					}
					
					plainContent := string(runes[contentIdx:endOfLine])
					result.WriteByte('<')
					result.WriteString(tag)
					result.WriteString(attrs)
					result.WriteByte('>')
					result.WriteString(plainContent)
					result.WriteString("</")
					result.WriteString(tag)
					result.WriteByte('>')
					i = endOfLine
					continue
				}
			}
		}

		result.WriteRune(runes[i])
		i++
	}

	return result.String()
}

func ToHTML(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Error: Missing text argument"
	}

	md := args[0].String()
	
	htmlParsed := parseInline(md)

	lines := strings.Split(htmlParsed, "\n")
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

		isBlock := strings.HasPrefix(trimmed, "<div") || strings.HasPrefix(trimmed, "</div") ||
			strings.HasPrefix(trimmed, "<ul") || strings.HasPrefix(trimmed, "</ul") ||
			strings.HasPrefix(trimmed, "<li") || strings.HasPrefix(trimmed, "</li") ||
			strings.HasPrefix(trimmed, "<h") || strings.HasPrefix(trimmed, "</h")

		if isBlock {
			if inParagraph {
				result = append(result, "</p>")
				inParagraph = false
			}
			result = append(result, line)
		} else {
			if !inParagraph {
				result = append(result, "<p>"+line)
				inParagraph = true
			} else {
				result = append(result, "<br>"+line)
			}
		}
	}

	if inParagraph {
		result = append(result, "</p>")
	}

	return strings.Join(result, "\n")
}
