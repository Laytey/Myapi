// Package middleware содержит HTTP-middleware для Gin.
package middleware

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter              // даёт zipWriter все методы gin.ResponseWriter
	gz                 *gzip.Writer // сюда будем писать
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	w.Header().Del("Content-Length")
	return w.gz.Write(data)
}

// GzipMiddleware включает поддержку gzip: распаковывает тело запроса
// (если задан Content-Encoding: gzip) и сжимает ответ (если клиент
// прислал Accept-Encoding: gzip).
//
// При распаковке запроса:
//   - Если заголовок Content-Encoding равен "gzip" — тело оборачивается
//     в gzip.Reader, который распаковывает данные по мере чтения.
//   - Если тело повреждено — возвращается 400 Bad Request.
//
// При сжатии ответа:
//   - Если клиент не поддерживает gzip — запрос пропускается без изменений.
//   - Иначе c.Writer оборачивается в gzipWriter, заголовки
//     Content-Encoding: gzip и Vary: Accept-Encoding устанавливаются.
//
// Middleware прозрачен для хендлеров: они работают с обычными данными.
func GzipMiddleware(c *gin.Context) {
	//распаковка запроса (Content-Encoding)
	if c.GetHeader("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid gzip body"})
			return
		}
		c.Request.Body = gzipReader
		defer func() {
			if err := gzipReader.Close(); err != nil {
				log.Printf("failed to close gzip reader: %v", err)
			}
		}() // закроется после завершения всей цепочки
	}

	// сжатие ответа (Accept-Encoding)
	// проверяем, готов ли клиент принять gzip-ответ
	if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
		c.Next() // клиент не поддерживает gzip, просто пропускаем запрос
		return
	}

	// создаем gzip-обертку для ответа
	gz := gzip.NewWriter(c.Writer)
	defer func() {
		if err := gz.Close(); err != nil {
			log.Printf("failed to close gzip writer: %v", err)
		}
	}() // важно: Close() дописывает остатки сжатых данных в поток

	// подменяем Writer в контексте
	c.Writer = &gzipWriter{
		ResponseWriter: c.Writer,
		gz:             gz,
	}

	// устанавливаем заголовки ответа
	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding") // Говорим кэшам, что ответ зависит от Accept-Encoding

	c.Next()
}
