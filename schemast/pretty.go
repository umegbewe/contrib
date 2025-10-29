package schemast

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	reEdges  = regexp.MustCompile(`return\s+\[\]ent\.Edge\{\s*([^}]*)\}`)
)

func prettyCompositeLiterals(src []byte) []byte {
	out := src

	out = reEdges.ReplaceAllFunc(out, func(m []byte) []byte {
		return prettyReturnBlock(m, "ent.Edge", "edge.")
	})
	return out
}
func prettyReturnBlock(m []byte, elemType, splitPrefix string) []byte {
	open := "return []" + elemType + "{"
	close := "}"
	s := string(m)

	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end < 0 || end <= start {
		return m
	}
	inner := strings.TrimSpace(s[start+1 : end])

	if inner == "" {
		return m
	}

	parts := splitCalls(inner, splitPrefix)

	var buf bytes.Buffer
	buf.WriteString(open)
	buf.WriteString("\n")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		p = strings.TrimSuffix(p, ",")
		buf.WriteString("\t\t")
		buf.WriteString(p)
		buf.WriteString(",\n")
	}
	buf.WriteString("\t")
	buf.WriteString(close)
	return []byte(buf.String())
}

func splitCalls(inner, splitPrefix string) []string {
	var parts []string
	var curr strings.Builder

	depth := 0
	i := 0
	for i < len(inner) {
		ch := inner[i]

		switch ch {
		case '(':
			depth++
			curr.WriteByte(ch)
			i++

		case ')':
			if depth > 0 {
				depth--
			}
			curr.WriteByte(ch)
			i++

		case ',':
			if depth == 0 {
				j := i + 1
				for j < len(inner) && (inner[j] == ' ' || inner[j] == '\t' || inner[j] == '\n' || inner[j] == '\r') {
					j++
				}
				if strings.HasPrefix(inner[j:], splitPrefix) {
					parts = append(parts, strings.TrimSpace(curr.String()))
					curr.Reset()
					i++
					continue
				}
			}
			curr.WriteByte(ch)
			i++

		default:
			curr.WriteByte(ch)
			i++
		}
	}
	if curr.Len() > 0 {
		parts = append(parts, strings.TrimSpace(curr.String()))
	}
	return parts
}
