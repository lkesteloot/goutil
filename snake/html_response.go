package snake

import (
	"bufio"
	"container/list"
	"html"
	"io"
)

const checkTagPairing = true

type HtmlResponse struct {
	w     *bufio.Writer
	stack *list.List
}

func New(w io.Writer) *HtmlResponse {
	return &HtmlResponse{bufio.NewWriter(w), list.New()}
}

func (h *HtmlResponse) WriteStr(s string) {
	h.w.WriteString(s)
}

func (h *HtmlResponse) Flush() *HtmlResponse {
	h.w.Flush()
	return h
}

func (h *HtmlResponse) openTag(tag string, attrs []attr) *HtmlResponse {
	h.stack.PushBack(tag)

	return h.singleTag(tag, attrs)
}

func (h *HtmlResponse) singleTag(tag string, attrs []attr) *HtmlResponse {
	h.WriteStr("<")
	h.WriteStr(tag)

	for _, a := range attrs {
		h.WriteStr(" ")
		h.WriteStr(a.name)
		h.WriteStr("=\"")
		h.WriteStr(html.EscapeString(a.value))
		h.WriteStr("\"")
	}

	h.WriteStr(">")

	return h
}

func (h *HtmlResponse) closeTag(tag string) *HtmlResponse {
	h.WriteStr("</")
	h.WriteStr(tag)
	h.WriteStr(">")

	back := h.stack.Back()
	if checkTagPairing {
		if back == nil {
			panic("too many closing tags")
		}
		backTag := back.Value.(string)
		if backTag != tag {
			panic("mismatched tag (/" + tag + ", expected /" + backTag + ")")
		}
	}
	h.stack.Remove(back)

	return h
}

func (h *HtmlResponse) autoCloseTag() *HtmlResponse {
	back := h.stack.Back()
	if back == nil {
		panic("too many closing tags")
	}
	backTag := back.Value.(string)

	return h.closeTag(backTag)
}

func (h *HtmlResponse) Doctype() *HtmlResponse {
	h.WriteStr("<!DOCTYPE html>\n")
	return h
}

func (h *HtmlResponse) Html(attrs ...attr) *HtmlResponse {
	return h.openTag("html", attrs)
}

func (h *HtmlResponse) Html_() *HtmlResponse {
	return h.closeTag("html")
}

func (h *HtmlResponse) Head(attrs ...attr) *HtmlResponse {
	return h.openTag("head", attrs)
}

func (h *HtmlResponse) Head_() *HtmlResponse {
	return h.closeTag("head")
}

func (h *HtmlResponse) Body(attrs ...attr) *HtmlResponse {
	return h.openTag("body", attrs)
}

func (h *HtmlResponse) Body_() *HtmlResponse {
	return h.closeTag("body")
}

func (h *HtmlResponse) Title(attrs ...attr) *HtmlResponse {
	return h.openTag("title", attrs)
}

func (h *HtmlResponse) Title_() *HtmlResponse {
	return h.closeTag("title")
}

func (h *HtmlResponse) Script(attrs ...attr) *HtmlResponse {
	return h.openTag("script", attrs)
}

func (h *HtmlResponse) Script_() *HtmlResponse {
	return h.closeTag("script")
}

func (h *HtmlResponse) ScriptLink(url string) *HtmlResponse {
	return h.Script(Src(url)).Script_()
}

func (h *HtmlResponse) Link(attrs ...attr) *HtmlResponse {
	return h.singleTag("link", attrs)
}

func (h *HtmlResponse) StyleLink(href string) *HtmlResponse {
	return h.Link(Rel("stylesheet"), Href(href))
}

func (h *HtmlResponse) Div(attrs ...attr) *HtmlResponse {
	return h.openTag("div", attrs)
}

func (h *HtmlResponse) Div_() *HtmlResponse {
	return h.closeTag("div")
}

func (h *HtmlResponse) P(attrs ...attr) *HtmlResponse {
	return h.openTag("p", attrs)
}

func (h *HtmlResponse) P_() *HtmlResponse {
	return h.closeTag("p")
}

func (h *HtmlResponse) B(attrs ...attr) *HtmlResponse {
	return h.openTag("b", attrs)
}

func (h *HtmlResponse) B_() *HtmlResponse {
	return h.closeTag("b")
}

func (h *HtmlResponse) A(attrs ...attr) *HtmlResponse {
	return h.openTag("a", attrs)
}

func (h *HtmlResponse) A_() *HtmlResponse {
	return h.closeTag("a")
}

func (h *HtmlResponse) Ul(attrs ...attr) *HtmlResponse {
	return h.openTag("ul", attrs)
}

func (h *HtmlResponse) Ul_() *HtmlResponse {
	return h.closeTag("ul")
}

func (h *HtmlResponse) Li(attrs ...attr) *HtmlResponse {
	return h.openTag("li", attrs)
}

func (h *HtmlResponse) Li_() *HtmlResponse {
	return h.closeTag("li")
}

func (h *HtmlResponse) Table(attrs ...attr) *HtmlResponse {
	return h.openTag("table", attrs)
}

func (h *HtmlResponse) Table_() *HtmlResponse {
	return h.closeTag("table")
}

func (h *HtmlResponse) Tr(attrs ...attr) *HtmlResponse {
	return h.openTag("tr", attrs)
}

func (h *HtmlResponse) Tr_() *HtmlResponse {
	return h.closeTag("tr")
}

func (h *HtmlResponse) Td(attrs ...attr) *HtmlResponse {
	return h.openTag("td", attrs)
}

func (h *HtmlResponse) Td_() *HtmlResponse {
	return h.closeTag("td")
}

func (h *HtmlResponse) Img(attrs ...attr) *HtmlResponse {
	return h.singleTag("img", attrs)
}

func (h *HtmlResponse) Meta(attrs ...attr) *HtmlResponse {
	return h.singleTag("meta", attrs)
}

func (h *HtmlResponse) Text(s string) *HtmlResponse {
	return h.RawText(html.EscapeString(s))
}

func (h *HtmlResponse) RawText(s string) *HtmlResponse {
	h.WriteStr(s)
	return h
}

func (h *HtmlResponse) Content(s string) *HtmlResponse {
	h.Text(s)
	return h.autoCloseTag()
}

func (h *HtmlResponse) RawContent(s string) *HtmlResponse {
	h.RawText(s)
	return h.autoCloseTag()
}
