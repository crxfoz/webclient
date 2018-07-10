package webclient

import (
	"log"
	"net/http"
	"net/url"
)

// Webclient Базовя структура, содержащая http.client и http.Transport (в том числе их инициализация)
type Webclient struct {
	transport *http.Transport
	client    *http.Client
}

// Get Отправить запрос методом GET
func (w *Webclient) Get(url string) *Request {
	return NewRequest(w.client, w.transport, url, http.MethodGet)
}

// Post Отправить запрос методом POST
func (w *Webclient) Post(url string) *Request {
	return NewRequest(w.client, w.transport, url, http.MethodPost)
}

// Head Отправить запрос методом HEAD
func (w *Webclient) Head(url string) *Request {
	return NewRequest(w.client, w.transport, url, http.MethodHead)
}

// Put Отправить запрос методом PUT
func (w *Webclient) Put(url string) *Request {
	return NewRequest(w.client, w.transport, url, http.MethodPut)
}

// Delete Отправить запрос методом DELETE
func (w *Webclient) Delete(url string) *Request {
	return NewRequest(w.client, w.transport, url, http.MethodDelete)
}

// Patch Отправить запрос методом PATCH
func (w *Webclient) Patch(url string) *Request {
	return NewRequest(w.client, w.transport, url, http.MethodPatch)
}

// Options Отправить запрос методом OPTIONS
func (w *Webclient) Options(url string) *Request {
	return NewRequest(w.client, w.transport, url, http.MethodOptions)
}

// Proxy Установить прокси для запросов
func (w *Webclient) Proxy(proxyURL string) *Webclient {
	p, err := url.Parse(proxyURL)
	if err != nil {
		log.Println("Cant parse proxy URL")
		return w
	}

	w.transport.Proxy = http.ProxyURL(p)

	return w
}

