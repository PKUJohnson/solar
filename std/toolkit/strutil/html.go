package strutil

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	gohtml "golang.org/x/net/html"
)

// HTML strips html tags, replace common entities, and escapes <>&;'" in the result.
// Note the returned text may contain entities as it is escaped by HTMLEscapeString, and most entities are not translated.
func RemoveHtmlTag(s string) (output string) {

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {

		// First remove line breaks etc as these have no meaning outside html tags (except pre)
		// this means pre sections will lose formatting... but will result in less unintentional paras.
		s = strings.Replace(s, "\n", "", -1)

		// Then replace line breaks with newlines, to preserve that formatting
		s = strings.Replace(s, "</p>", "\n", -1)
		s = strings.Replace(s, "<br>", "\n", -1)
		s = strings.Replace(s, "</br>", "\n", -1)
		s = strings.Replace(s, "<br/>", "\n", -1)
		s = strings.Replace(s, "<br />", "\n", -1)

		// Walk through the string removing all tags
		b := bytes.NewBufferString("")
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}

	// Remove a few common harmless entities, to arrive at something more like plain text
	output = strings.Replace(output, "&#8216;", "'", -1)
	output = strings.Replace(output, "&#8217;", "'", -1)
	output = strings.Replace(output, "&#8220;", "\"", -1)
	output = strings.Replace(output, "&#8221;", "\"", -1)
	output = strings.Replace(output, "&nbsp;", " ", -1)
	output = strings.Replace(output, "&quot;", "\"", -1)
	output = strings.Replace(output, "&apos;", "'", -1)

	// Translate some entities into their plain text equivalent (for example accents, if encoded as entities)
	output = html.UnescapeString(output)

	// In case we have missed any tags above, escape the text - removes <, >, &, ' and ".
	output = template.HTMLEscapeString(output)

	// After processing, remove some harmless entities &, ' and " which are encoded by HTMLEscapeString
	output = strings.Replace(output, "&#34;", "\"", -1)
	output = strings.Replace(output, "&#39;", "'", -1)
	output = strings.Replace(output, "&amp;", "&", -1)       // NB space after
	output = strings.Replace(output, "&amp; ", "& ", -1)     // NB space after
	output = strings.Replace(output, "&amp;amp; ", "& ", -1) // NB space after

	return output
}

func RemoveHtmlTagExceptBlank(s string) (output string) {
	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {

		// First remove line breaks etc as these have no meaning outside html tags (except pre)
		// this means pre sections will lose formatting... but will result in less unintentional paras.

		// Then replace line breaks with newlines, to preserve that formatting
		s = strings.Replace(s, "</p>", "\n", -1)
		s = strings.Replace(s, "<br>", "\n", -1)
		s = strings.Replace(s, "</br>", "\n", -1)
		s = strings.Replace(s, "<br/>", "\n", -1)
		s = strings.Replace(s, "<br />", "\n", -1)

		// Walk through the string removing all tags
		b := bytes.NewBufferString("")
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}

	// Remove a few common harmless entities, to arrive at something more like plain text
	output = strings.Replace(output, "&#8216;", "'", -1)
	output = strings.Replace(output, "&#8217;", "'", -1)
	output = strings.Replace(output, "&#8220;", "\"", -1)
	output = strings.Replace(output, "&#8221;", "\"", -1)
	output = strings.Replace(output, "&nbsp;", " ", -1)
	output = strings.Replace(output, "&quot;", "\"", -1)
	output = strings.Replace(output, "&apos;", "'", -1)

	// Translate some entities into their plain text equivalent (for example accents, if encoded as entities)
	output = html.UnescapeString(output)

	// In case we have missed any tags above, escape the text - removes <, >, &, ' and ".
	output = template.HTMLEscapeString(output)

	// After processing, remove some harmless entities &, ' and " which are encoded by HTMLEscapeString
	output = strings.Replace(output, "&#34;", "\"", -1)
	output = strings.Replace(output, "&#39;", "'", -1)
	output = strings.Replace(output, "&amp; ", "& ", -1)     // NB space after
	output = strings.Replace(output, "&amp;amp; ", "& ", -1) // NB space after

	return output
}

func RemoveSpecialCharacters(s string) string {
	s = strings.Replace(s, "\u2028", "", -1)
	s = strings.Replace(s, "\u2029", "", -1)

	return s
}

func RemoveEmptyPTag(s string) string {
	s = strings.Replace(s, "<p>&nbsp;</p>", "", -1)

	return s
}

func Html2(s *goquery.Selection) (ret string, e error) {
	// Since there is no .innerHtml, the HTML content must be re-created from
	// the nodes using html.Render.
	var buf bytes.Buffer

	if len(s.Nodes) > 0 {
		for c := s.Nodes[0]; c != nil; c = c.NextSibling {
			e = gohtml.Render(&buf, c)
			if e != nil {
				return
			}
		}
		ret = buf.String()
	}

	return
}

func RemoveStylesOfHtmlTag(html string, keptStyles ...string) string {
	styleTagR := regexp.MustCompile(`style="([^=>]*;?)*"`)
	keptStyleR := regexp.MustCompile(fmt.Sprintf("^%s", strings.Join(keptStyles, "|")))
	hasKeptStyles := len(keptStyles) != 0

	matches := styleTagR.FindAllStringSubmatch(html, -1)
	var toCleanStyles = make(map[string]struct{})

	for _, groups := range matches {
		if len(groups)!=2{
			continue
		}

		styles := strings.Split(groups[1], ";")
		for _, style := range styles{
			trimStr := strings.ToLower(strings.TrimSpace(style))
			if trimStr == "" || (hasKeptStyles && keptStyleR.MatchString(trimStr)) {
				continue
			}

			toCleanStyles[style] = struct{}{}
		}
	}

	cleanedHtml := html
	for cleanStyle := range toCleanStyles{
		escapedStyle := strings.Replace(cleanStyle, "(", `\(`, -1)
		escapedStyle = strings.Replace(escapedStyle, ")", `\)`, -1)
		cleanR := regexp.MustCompile(fmt.Sprintf("%s[;]?", escapedStyle))
		cleanedHtml = cleanR.ReplaceAllString(cleanedHtml, "")
	}

	return cleanedHtml
}