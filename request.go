package webclient

// todo: Добавить методы SendJSON и SendXML
// todo: AddHeader в дополнение к SetHeader?

import (
	"net/url"
	"log"
	"net/http"
	"io/ioutil"
	"io"
	"bytes"
	"mime/multipart"
	"net/textproto"
	"fmt"
)

// Request Структура содержащая все состовные части для запроса
type Request struct {
	client *http.Client

	url       string
	ctype     WContentType
	customCType WContentType
	method    string
	headers   map[string]string
	cookies   map[string]string
	queryData map[string][]string
	formData  map[string][]string
	files []File
}

// NewRequest Создает новый Request
func NewRequest(client *http.Client, transport *http.Transport, targetURL string, method string) *Request {
	client.Transport = transport

	return &Request{
		client: client,
		url: targetURL,
		method: method,
		headers: make(map[string]string),
		cookies: make(map[string]string),
		files: make([]File, 0),
		queryData: make(map[string][]string),
		formData:  make(map[string][]string),
	}
}

// Cookie Добавляет куку
func (r *Request) Cookie(name string, value string) *Request {
	r.cookies[name] = value
	return r
}

// ContentType Устанавливает хидер Content-Type
func (r *Request) ContentType(name WContentType) *Request {
	r.customCType = name
	return r
}

// SetHeader Устанавливает заголовок для запроса
func (r *Request) SetHeader(header string, data string) *Request {
	r.headers[header] = data
	return r
}

// UserAgent Устанавливает заголовок User-Agent для запроса
func (r *Request) UserAgent(data string) *Request {
	r.headers["User-Agent"] = data
	return r
}

// Referer Устанавливает заголовок Referer для запроса
func (r *Request) Referer(data string) *Request {
	r.headers["Referer"] = data
	return r
}


// Query Устанавливает Query данные для запроса
//
// Пример запроса: /search?foo=bar&foz=baz
//      client.
//          Get("/page").
//          QueryString("foo=bar&foz=baz")
func (r *Request) Query(data string) *Request {
	parsed, err := url.ParseQuery(data)
	if err != nil {
		log.Panic(err)
		return r
	}

	for key, values := range parsed {
		for _, v := range values {
			r.queryData[key] = append(r.queryData[key], v)
		}
	}

	return r
}

// QueryParam Устанавливает Query данные для запроса
//
// Пример запроса: /search?foo=bar&foz=baz
//      client.
//          Get("/page").
//          Query("foo", "bar").
//          Query("foz", "baz")
func (r *Request) QueryParam(key string, value string) *Request {
	r.queryData[key] = append(r.queryData[key], value)
	return r
}

// Send Добавить данные для POSTDATA
//
// Пример запроса: foo=bar&foz=baz
//      client.
//          Post("/page").
//          Send("foo=bar&foz=baz")
//
//      client.
//          Post("/page").
//          Send("foo=bar").
//          Send("foz=baz")
func (r *Request) Send(data string) *Request {
	parsed, err := url.ParseQuery(data)
	if err != nil {
		log.Panic(err)
		return r
	}

	for key, values := range parsed {
		for _, v := range values {
			r.queryData[key] = append(r.queryData[key], v)
		}
	}

	return r
}

// SendParam Устанавливает данные для POSTDATA
//
// Пример запроса: foo=bar&foz=baz
//      client.
//          Post("/page").
//          Send("foo", "bar").
//          Send("foz", "baz")
func (r *Request) SendParam(key string, value string) *Request {
	r.formData[key] = append(r.formData[key], value)

	return r
}

// SendFile Добавить отправку файла к запросу. В этом случае будет отправлен multipart-запрос
//
// Пример:
// 		f, _ := NewFile("/tmp/1.txt", "userfile")
// 		client.Post('/upload').
//		SendFile(f)
func (r *Request) SendFile(file File) *Request {
	r.files = append(r.files, file)

	return r
}

// newRequest Собирает воедино http.Request
func (r *Request) newRequest() (*http.Request, error) {
	var (
		data io.Reader
		req  *http.Request
		err  error
	)

	// Если имеются добавленные файлы -> используем multipart запрос
	if len(r.files) > 0 || r.customCType == TypeMultipart {
		buf := &bytes.Buffer{}
		multipartWriter := multipart.NewWriter(buf)

		// К multipart запросу добавляются данные из formData
		for key, values := range r.formData {
			for _, v := range values {
				fw, _ := multipartWriter.CreateFormField(key)
				fw.Write([]byte(v))
			}
		}

		for _, file := range r.files {
			var fw io.Writer
			// Если указан конкретный Content-Type для файла -> Используем его
			// Иначе -> Используем Content-Type по умолчанию (application/octet-stream)
			if len(file.ContentType) > 0 {
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, file.Param, file.Name))
				h.Set("Content-Type", file.ContentType)
				fw, _ = multipartWriter.CreatePart(h)
			} else {
				fw, _ = multipartWriter.CreateFormFile(file.Param, file.Name)
			}

			fw.Write(file.Data)
		}

		data = buf

		multipartWriter.Close()

		// Указывает правильный Content-Type для multipart запроса (включая boundary)
		r.ctype = WContentType(multipartWriter.FormDataContentType())

	} else {
		if len(r.formData) > 0 {
			b := []byte(mapToUrlValues(r.formData).Encode())
			data = bytes.NewReader(b)

			r.ctype = TypeForm
		}
	}

	if req, err = http.NewRequest(r.method, r.url, data); err != nil {
		return nil, err
	}

	// Если установлен кастомный Content-Type -> Используем его
	// Иначе используется тот, что определила либа (или пустой)
	if len(r.customCType) > 0 && r.customCType != TypeMultipart {
		// TypeMultipart Как кастомный игнорируется
		// в данном блоке, т.к. обрабатывается выше в коде
		req.Header.Set("Content-Type", string(r.customCType))
	} else {
		if len(r.ctype) > 0 {
			req.Header.Set("Content-Type", string(r.ctype))
		}
	}

	// Энкодим Query-часть запроса
	req.URL.RawQuery = mapToUrlValues(r.queryData).Encode()

	// Устанавливаем хидеры
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	// Добавляем кукисы
	for k, v := range r.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	return req, nil
}

// Do Выполняет запрос
// warning: Не считывать тело запроса с resp.Body, для получения контента используется второй возвращаемый параметр
func (r *Request) Do() (*http.Response, string, error) {
	req, err := r.newRequest()
	if err != nil {
		return nil, "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return resp, string(body), nil
}
