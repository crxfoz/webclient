package webclient

import (
	"io/ioutil"
	"path/filepath"
)

// File Структура описывающая файл для отправки через client.SendFile
type File struct {
	// Название файла. По умолчанию используется имя переданного файла
	Name string
	// Тело файла
	Data []byte
	// Параметр запроса, который отвечает за принятие файла
	Param string
	// ContentType переданного файла. Можно указать свой, иначе будет установлен соответсвующий разрешению файла
	ContentType string
}

// NewFile Создает новую структуру File.
// error возможен только если переданный файл не удалось прочитать
func NewFile(path string, param string) (File, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return File{}, err
	}

	file := File{Data: data, Param: param}
	filename := filepath.Base(path)
	file.Name = filename

	return file, nil
}