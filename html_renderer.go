package main

import (
	"bytes"
	"fmt"
	"io"
	"unicode"

	"github.com/russross/blackfriday/v2"
)

const (
	HeadingStartTag = `
<h%d>
  <a name="%s" class="anchor" href="#%s" rel="nofollow" aria-hidden="true">
    <span class="octicon octicon-link">
    </span>
  </a>`

	HeadingEndTag = `</h%d>`
)

type renderer struct {
	delegate *blackfriday.HTMLRenderer
}

func (r *renderer) RenderHeader(w io.Writer, ast *blackfriday.Node) { r.delegate.RenderHeader(w, ast) }
func (r *renderer) RenderFooter(w io.Writer, ast *blackfriday.Node) { r.delegate.RenderFooter(w, ast) }

// RenderNode overrides blackfriday renderer.
// This overriden method applies anchor link if the node is Heading.
// Else, it fallbacks to blackfriday standard method.
func (r *renderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Type != blackfriday.Heading {
		return r.delegate.RenderNode(w, node, entering)
	}

	if entering {
		// node.FirstChild is text of the heading. Use it as the anchor link name
		anchor := sanitizeAnchorName(string(node.FirstChild.Literal))
		w.Write([]byte(fmt.Sprintf(HeadingStartTag, node.HeadingData.Level, anchor, anchor)))
	} else {
		w.Write([]byte(fmt.Sprintf(HeadingEndTag, node.HeadingData.Level) + "\n"))
	}

	return blackfriday.GoToNext
}

func sanitizeAnchorName(s string) string {
	var anchorName []rune
	var futureDash = false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}

func escape(src []byte) []byte {
	escapeChars := map[byte]string{'"': "&quot;", '&': "&amp;", '<': "&lt;", '>': "&gt;"}
	buff := []byte{}
	out := bytes.NewBuffer(buff)
	org := 0
	for i, ch := range src {
		if entity, ok := escapeChars[ch]; ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}

	return out.Bytes()
}
