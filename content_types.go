package webclient

type WContentType string

const (
	TypeHTML       WContentType = "text/html"
	TypeJSON       WContentType = "application/json"
	TypeXML        WContentType = "application/xml"
	TypeText       WContentType = "text/plain"
	TypeForm       WContentType = "application/x-www-form-urlencoded"
	TypeMultipart  WContentType = "multipart/form-data"
)
