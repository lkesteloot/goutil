package snake

type attr struct {
	name  string
	value string
}

func Attr(name string, value string) attr {
	return attr{name, value}
}

func Class(class string) attr {
	return Attr("class", class)
}

func Id(id string) attr {
	return Attr("id", id)
}

func Rel(rel string) attr {
	return Attr("rel", rel)
}

func Href(href string) attr {
	return Attr("href", href)
}

func Src(src string) attr {
	return Attr("src", src)
}

func Title(title string) attr {
	return Attr("title", title)
}

func TabIndex(tabindex string) attr {
	return Attr("tabindex", tabindex)
}

func Charset(tabindex string) attr {
	return Attr("charset", tabindex)
}

func Property(tabindex string) attr {
	return Attr("property", tabindex)
}

func Content(tabindex string) attr {
	return Attr("content", tabindex)
}

func Name(tabindex string) attr {
	return Attr("name", tabindex)
}
