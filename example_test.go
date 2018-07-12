package webclient

import (
	"fmt"
)

func ExampleConfig_New() {
	client := Config{Timeout: 10, FollowRedirect: false}.New()
	resp, body, err := client.Get("http://github.com/").Referer("https://google.com").Do()
	if err != nil {
		// ...
	}

	if resp.StatusCode != 200 {
		// ...
	}

	fmt.Println(body)
}

func ExampleRequest_ContentType() {
	client := Config{}.New()
	rawJSON := `{"raw": "json"}`

	client.Post("https://example.com/").ContentType(TypeJSON).SendPlain(rawJSON).Do()

	// custom content-type
	client.Post("https://example.com/").ContentType(WContentType("application/json")).SendPlain(rawJSON).Do()
}

func ExampleRequest_SendFile() {
	f, _ := NewFile("/tmp/1.txt", "userfile")
	client := Config{UseKeepAlive: false}.New()
	client.Post("http://example.com/uploader").SendFile(f).Do()
}

func ExampleRequest_SendFiles() {
	names := []string{"apple", "orange", "horse"}
	files := make([]File, len(names))

	for _, name := range names {
		f, _ := NewFile(fmt.Sprintf("/tmp/%s.txt", name), "userfile")
		files = append(files, f)
	}
	client := Config{UseKeepAlive: false}.New()
	client.Post("http://example.com/uploader").SendFiles(files...).Do()
}

func ExampleRequest_Query() {
	// Пример запроса: /search?foo=bar&foz=baz
	client := Config{UseKeepAlive: false}.New()
	client.Get("http://example.com/search").Query("foo=bar&foz=baz").Do()
}

func ExampleRequest_QueryParam() {
	// Пример запроса: /search?foo=bar&foz=baz
	client := Config{UseKeepAlive: false}.New()
	client.Get("http://example.com/search").
		QueryParam("foo", "bar").
		QueryParam("foz", "baz").
		Do()
}

func ExampleRequest_Send() {
	// Пример запроса: POST http://example.com/submit
	// data: foo=bar&foz=baz

	client := Config{UseKeepAlive: false}.New()
	client.Post("http://example.com/submit").Send("foo=bar&foz=baz").Do()
	// or
	client.Post("http://example.com/submit").
		Send("foo=bar").
		Send("foz=baz").
		Do()
}

func ExampleRequest_SendParam() {
	// Пример запроса: POST http://example.com/submit
	// data: foo=bar&foz=baz

	client := Config{UseKeepAlive: false}.New()
	client.Post("http://example.com/submit").
		SendParam("foo", "bar").
		SendParam("foz", "baz").
		Do()
}

func ExampleRequest_SendStruct() {
	// Пример запроса: POST http://example.com/submit
	// data: {"id":"32131","browser":"Opera","browser_ver":"10"}

	client := Config{UseKeepAlive: false}.New()

	uifno := struct {
		ID         string `json:"id"`
		Browser    string `json:"browser"`
		BrowserVer string `json:"browser_ver"`
	}{
		ID:         "32131",
		Browser:    "Opera",
		BrowserVer: "10",
	}

	client.Post("http://example.com/track").SendStruct(uifno).Do()
}

func ExampleRequest_SendJSON() {
	// Пример запроса: POST http://example.com/submit
	// data: {"id":"32131","browser":"Opera","browser_ver":"10"}

	data := `{"id":"32131","browser":"Opera","browser_ver":"10"}`
	client := Config{UseKeepAlive: false}.New()

	client.Post("http://example.com/track").SendJSON(data).Do()
}

func ExampleRequest_SendXML() {
	// Пример запроса: POST http://example.com/submit
	// data: <id>32131</id><browser>Opera</browser>

	data := `<id>32131</id><browser>Opera</browser>`
	client := Config{UseKeepAlive: false}.New()

	client.Post("http://example.com/track").SendXML(data).Do()
}

func ExampleRequest_Cookie() {
	client := Config{UseKeepAlive: false}.New()

	client.Post("http://example.com/myprofile").
		Cookie("sessionid", "...").
		Do()
}

func ExampleRequest_SendPlain() {
	client := Config{}.New()
	rawJSON := `{"raw": "json"}`

	client.Post("https://example.com/").
		ContentType(TypeJSON).
		SendPlain(rawJSON).
		Do()
}
