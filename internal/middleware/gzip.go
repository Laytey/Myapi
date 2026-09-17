package middleware

import (
	"compress/gzip"
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

func GzipMiddleware(c *gin.Context) {
	//распаковка запроса (Content-Encoding)
	if c.GetHeader("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid gzip body"})
			return
		}
		c.Request.Body = gzipReader
		defer gzipReader.Close() // закроется после завершения всей цепочки
	}

	// сжатие ответа (Accept-Encoding)
	// проверяем, готов ли клиент принять gzip-ответ
	if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
		c.Next() // клиент не поддерживает gzip, просто пропускаем запрос
		return
	}

	// создаем gzip-обертку для ответа
	gz := gzip.NewWriter(c.Writer)
	defer gz.Close() // важно: Close() дописывает остатки сжатых данных в поток

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
