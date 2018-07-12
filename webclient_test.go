package webclient

// todo: Добавить больше тестов для специфичиских кейсов
// todo: проверять на urlencode, пробелы, пустые параметры

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMapToUrlValues(t *testing.T) {
	data := map[string][]string{
		"foo": {"bar"},
		"b":   {"1", "2"},
		"4":   {"6"},
	}

	res := mapToUrlValues(data)
	if res.Get("foo") != "bar" {
		t.Errorf("Expected %s, got %v", "bar", res.Get("foo"))
	}

	if res.Get("4") != "6" {
		t.Errorf("Expected %s, got %v", "6", res.Get("4"))
	}

	if res.Get("b") != "1" {
		t.Errorf("Expected %s, got %v", "1", res.Get("b"))
	}

	if res.Encode() != "4=6&b=1&b=2&foo=bar" {
		t.Errorf("Encode failed, expected '%s', got '%v'", "4=6&b=1&b=2&foo=bar", res.Encode())
	}
}

func TestTimeoutWithout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
	}))

	client := Config{}.New()
	_, _, err := client.Get(ts.URL).Do()
	if err != nil {
		t.Errorf("Got an unexpected error: %v", err)
	}

	defer ts.Close()
}

func TestTimeoutSet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
	}))

	defer ts.Close()

	client := Config{Timeout: 2}.New()
	_, _, err := client.Get(ts.URL).Do()
	if err == nil {
		t.Errorf("Timeout doesn't work, request should be failed")
	}

	if err, ok := err.(net.Error); !ok || !err.Timeout() {
		t.Errorf("Expecting to get Timeout error, got: %v", err)
	}
}

func TestFollowFalse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://httpbin.org/absolute-redirect/2")
		w.WriteHeader(302)
	}))

	defer ts.Close()

	client := Config{FollowRedirect: false}.New()
	resp, _, err := client.Get(ts.URL).Do()
	if err != nil {
		t.Errorf("Got an unexpected error (probably test isn't failed but service is gone): %v", err)
	}

	if len(resp.Header.Get("Location")) == 0 {
		t.Errorf("Didn't get Location header as expected")
	}

	if resp.Request.URL.String() != ts.URL {
		t.Errorf("Got wrong URL, expected: %s, got: %s", ts.URL, resp.Request.URL.String())
	}
}

func TestFollowTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://httpbin.org/absolute-redirect/2")
		w.WriteHeader(302)
	}))

	defer ts.Close()

	client := Config{FollowRedirect: true}.New()
	resp, _, err := client.Get(ts.URL).Do()
	if err != nil {
		t.Errorf("Got an unexpected error (probably test isn't failed but service is gone): %v", err)
	}

	if ts.URL == resp.Request.URL.String() {
		t.Errorf("Urls are same. But we should be redirected")
	}
}

func TestMethods(t *testing.T) {
	methods := []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS",
		"HEAD",
		"PATCH",
	}

	for _, method := range methods {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				t.Errorf("Expected method %s, got: %s", method, r.Method)
			}
		}))

		client := Config{}.New()
		switch method {
		case "GET":
			client.Get(ts.URL).Do()
		case "POST":
			client.Post(ts.URL).Do()
		case "PUT":
			client.Put(ts.URL).Do()
		case "DELETE":
			client.Delete(ts.URL).Do()
		case "OPTIONS":
			client.Options(ts.URL).Do()
		case "HEAD":
			client.Head(ts.URL).Do()
		case "PATCH":
			client.Patch(ts.URL).Do()
		default:
			t.Errorf("Got unexpected method: %s", method)
		}

		ts.Close()
	}
}

func TestHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "Opera 1.0" {
			t.Errorf("Didn't get expected header User-Agent, expected: %s, got: %s", "Opera 1.0", r.Header.Get("User-Agent"))
		}

		if r.Header.Get("Referer") != "https://google.com/" {
			t.Errorf("Didn't get expected header Referer, expected: %s, got: %s", "https://google.com/", r.Header.Get("Referer"))
		}

		if r.Header.Get("Accept") != "text/html" {
			t.Errorf("Didn't get expected header Accept, expected: %s, got: %s", "text/html", r.Header.Get("Accept"))
		}

	}))
	defer ts.Close()

	client := Config{}.New()
	client.Get(ts.URL).
		UserAgent("Opera 1.0").
		Referer("https://google.com/").
		SetHeader("Accept", "text/html").
		Do()
}

func TestQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/path" {
			t.Errorf("Expected path: %s, got: %s", "/path", r.URL.Path)
		}
		params := r.URL.Query()
		if params.Get("foo") != "bar" {
			t.Errorf("Expected query param. Expected: %s, got: %s", "bar", params.Get("foo"))
		}

		if params.Get("foz") != "boz" {
			t.Errorf("Expected query param. Expected: %s, got: %s", "boz", params.Get("foz"))
		}

		if params.Get("l[]") != "123" {
			t.Errorf("Expected query param. Expected: %s, got: %s", "123", params.Get("l[]"))
		}

	}))

	defer ts.Close()

	client := Config{}.New()
	client.Get(ts.URL+"/path").Query("foo=bar&foz=boz").QueryParam("l[]", "123").Do()
}

func TestPostParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/path2" {
			t.Errorf("Expected path: %s, got: %s", "/path", r.URL.Path)
		}

		r.ParseForm()

		params := r.Form
		if params.Get("foo") != "bar" {
			t.Errorf("Expected query param. Expected: %s, got: %s", "bar", params.Get("foo"))
		}

		if params.Get("foz") != "boz" {
			t.Errorf("Expected query param. Expected: %s, got: %s", "boz", params.Get("foz"))
		}

		if params.Get("var") != "val" {
			t.Errorf("Expected query param. Expected: %s, got: %s", "val", params.Get("var"))
		}

		encoded := r.Form.Encode()
		if strings.Count(encoded, "foo=bar") != 2 {
			t.Errorf("Two same parameters should be acceptable: %v", encoded)
		}

	}))

	defer ts.Close()

	client := Config{}.New()
	client.Post(ts.URL+"/path2").Send("foo=bar&foz=boz").SendParam("var", "val").SendParam("foo", "bar").Do()
}

func TestMultipartFile(t *testing.T) {
	const case01_default_file = "/default_file"
	const case02_custom_file = "/custom_file"
	const case03_no_file = "/no_file"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got: %s", r.Method)
		}

		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("Expected content-type: %s, got: %s", "multipart/form-data", r.Header.Get("Content-Type"))
		}
		r.ParseMultipartForm(4096)

		query := r.URL.Query()
		if query.Get("q1") != "a" || query.Get("q2") != "b" {
			t.Errorf("Unexpected Query parameters")
		}

		switch r.URL.Path {
		case case01_default_file:
			if r.MultipartForm.Value["s1"][0] != "a" || r.MultipartForm.Value["s2"][0] != "b" {
				t.Errorf("Unexpected body of request")
			}

			ct := r.MultipartForm.File["userfile"][0].Header.Get("Content-Type")
			if ct != "application/octet-stream" {
				t.Errorf("Expected content type of file: %s, got: %s", "application/octet-stream", ct)
			}
		case case02_custom_file:
			f := r.MultipartForm.File["userfile"][0]
			if f.Header.Get("Content-Type") != "image/png" {
				t.Errorf("Expected content type of file: %s, got: %s", "image/png", f.Header.Get("Content-Type"))
			}

			if f.Filename != "newname.txt" {
				t.Errorf("Expected name of file: %s, got: %s", "newname.txt", f.Filename)

			}
		case case03_no_file:
			if r.MultipartForm.Value["s1"][0] != "a" || r.MultipartForm.Value["s2"][0] != "b" {
				t.Errorf("Unexpected body of request")
			}

		default:
			t.Error("Unexpected path")
		}
	}))

	defer ts.Close()

	f, err := NewFile("./README.md", "userfile")
	if err != nil {
		t.Errorf("Cant open file: %v", err)
	}

	Config{}.New().
		Post(ts.URL+case01_default_file).
		QueryParam("q1", "a").
		QueryParam("q2", "b").
		SendParam("s1", "a").
		SendParam("s2", "b").
		SendFile(f).
		Do()

	f.Name = "newname.txt"
	f.ContentType = "image/png"

	Config{}.New().
		Post(ts.URL+case02_custom_file).
		QueryParam("q1", "a").
		QueryParam("q2", "b").
		SendParam("s1", "a").
		SendParam("s2", "b").
		SendFile(f).
		Do()

	Config{}.New().
		Post(ts.URL+case03_no_file).
		ContentType(TypeMultipart).
		QueryParam("q1", "a").
		QueryParam("q2", "b").
		SendParam("s1", "a").
		SendParam("s2", "b").
		Do()
}

func TestCustomContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != string(TypeJSON) {
			t.Errorf("Expect header content-type: %s, got: %s", TypeJSON, r.Header.Get("Content-Type"))
		}

		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Got unexpected error: %v", err)
		}

		if string(d) != "foo=bar" {
			t.Errorf("Expected body: %s, got: %s", "foo=bar", string(d))
		}

	}))

	Config{}.New().Post(ts.URL).ContentType(TypeJSON).SendParam("foo", "bar").Do()
}

func TestRequest_SendJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != string(TypeJSON) {
			t.Errorf("Got unexpected content-tyoe header. Expect: %s, got: %s",
				string(TypeJSON),
				r.Header.Get("Content-Type"))
		}

		type Case01Struct struct {
			Name  string `json:"name"`
			Films []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"films"`
		}

		var ss Case01Struct
		data, _ := ioutil.ReadAll(r.Body)

		err := json.Unmarshal(data, &ss)
		if err != nil {
			t.Errorf("Got unexpected error: %v", err)
		} else {
			if ss.Name != "gleb" {
				t.Errorf("Expected: %s, got: %s", "gleb", ss.Name)
			}
		}
	}))

	Config{}.New().Post(ts.URL).SendJSON(`{"name": "gleb", "films": [{"id": 1, "name": "foo"}, {"id": 2, "name": "bar"}]}`).Do()
}

func TestRequest_SendXML(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Case01Struct struct {
			Food []struct {
				Name     string `xml:"name"`
				Calories string `xml:"calories"`
			} `xml:"food"`
		}

		var ss Case01Struct
		data, _ := ioutil.ReadAll(r.Body)

		err := xml.Unmarshal(data, &ss)
		if err != nil {
			t.Errorf("Got unexpected error: %v", err)
		}

		if ss.Food[0].Name != "foo" {
			t.Errorf("Expected food[0].Name: %s, got: %s", "foo", ss.Food[0].Name)
		}

	}))

	Config{}.New().Post(ts.URL).SendXML(
		`<?xml version="1.0" encoding="UTF-8"?><menu><food><name>foo</name><calories>650</calories></food>
		<food><name>bar</name><calories>777</calories></food></menu>`).
		Do()
}

type Person struct {
	Name string `json:"name" xml:"name"`
	Pets []Pet  `json:"pets" xml:"pets"`
}

type Pet struct {
	ID  int `json:"id" xml:"id"`
	Age int `json:"age" xml:"age"`
}

func TestRequest_SendStruct(t *testing.T) {
	const case01_json = "/json"
	const case02_xml = "/xml"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, _ := ioutil.ReadAll(r.Body)

		switch r.URL.Path {
		case case01_json:
			if r.Header.Get("Content-Type") != string(TypeJSON) {
				t.Errorf("Unexpected content-type header. Expect: %s, got: %s", TypeJSON, r.Header.Get("Content-Type"))
			}

			if string(b) != `[{"name":"foo","pets":[{"id":1,"age":13},{"id":2,"age":14}]},{"name":"bar","pets":[{"id":3,"age":6}]}]` {
				t.Errorf("Unexpected body")
			}
		case case02_xml:
			if r.Header.Get("Content-Type") != string(TypeXML) {
				t.Errorf("Unexpected content-type header. Expect: %s, got: %s", TypeXML, r.Header.Get("Content-Type"))
			}

			if string(b) != `<Person><name>foo</name><pets><id>1</id><age>13</age></pets><pets><id>2</id><age>14</age></pets></Person><Person><name>bar</name><pets><id>3</id><age>6</age></pets></Person>` {
				t.Errorf("Unexpected body")
			}

		default:
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
	}))

	b := []Person{
		{
			Name: "foo",
			Pets: []Pet{
				{1, 13},
				{2, 14}},
		},
		{
			Name: "bar",
			Pets: []Pet{
				{3, 6}},
		},
	}

	Config{}.New().Post(ts.URL + case01_json).SendStruct(&b).ContentType(TypeJSON).Do()
	Config{}.New().Post(ts.URL + case02_xml).SendStruct(&b).ContentType(TypeXML).Do()
	Config{}.New().Post(ts.URL + case01_json).SendStruct(&b).Do()
}

func TestRequest_SetHeaders(t *testing.T) {
	m := map[string]string{
		"User-Agent":      "Mozilla/5.0",
		"Accept":          "text/html",
		"Accept-Language": "en-EN",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "Mozilla/5.0" {
			t.Errorf("Expected header User-Agent: %s, got: %s", "Mozilla/5.0", r.Header.Get("User-Agent"))
		}

		if r.Header.Get("Accept") != "text/html" {
			t.Errorf("Expected header Accept: %s, got: %s", "text/html", r.Header.Get("Accept"))
		}

		if r.Header.Get("Accept-Language") != "en-EN" {
			t.Errorf("Expected header Accept-Language: %s, got: %s", "en-EN", r.Header.Get("User-Agent"))
		}

		if r.Header.Get("Referer") != "https://google.com" {
			t.Errorf("Expected header Referer: %s, got: %s", "https://google.com", r.Header.Get("Referer"))
		}
	}))

	Config{}.New().Get(ts.URL).SetHeaders(m).SetHeader("Referer", "https://google.com").Do()

}
