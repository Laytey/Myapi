// Package auth предоставляет JWT-аутентификацию: генерацию и проверку
// токенов, а также middleware для защиты маршрутов.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte

// SetSecret устанавливает секретный ключ для подписи и проверки JWT.
// Должен быть вызван один раз при старте приложения (обычно в main)
// до первого вызова GenerateToken или ParseToken.
//
// Без установленного секрета токены подписываются пустым ключом —
// это небезопасно и подходит только для тестов.
func SetSecret(s string) {
	secret = []byte(s)
}

// Claims описывает полезную нагрузку JWT, используемую в приложении.
//
// Содержит UserID (строка) и стандартные поля jwt.RegisteredClaims:
// exp (время истечения), iat (время выпуска) и другие.
type Claims struct {
	UserID               string `json:"user_id"` // тег: при превращении в JSON назовет это поле user_id
	jwt.RegisteredClaims        // структура из пакета jwt, стандартные поля: exp, iat, iss, aud
}

// GenerateToken создаёт и подписывает JWT для указанного пользователя.
//
// userID — идентификатор пользователя (в строковом виде).
// exp — время жизни токена (например, 24 * time.Hour).
//
// Возвращает подписанную строку токена или ошибку.
// Требует предварительного вызова SetSecret.
func GenerateToken(userID string, exp time.Duration) (string, error) {
	// exp	Expiration Time	Когда токен истечёт
	// exp time.Duration - время жизни токена
	claims := Claims{ // создаём экземпляр структуры Claims и заполняем поля
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)), // время истечения токена
			IssuedAt:  jwt.NewNumericDate(time.Now()),          // запоминает момент появления токена
			// jwt.NewNumericDate фиксирует время с помощью time.Now()
			//
		},
	}
	token := jwt.NewWithClaims( // NewWithClaimsем новый JWT-токен
		jwt.SigningMethodHS256, // алгоритм подписи токена с использованием общего секрета (симметричное шифрование)
		// симмитричное шифрование - для шифровки и расшифровки нужен один и тот же ключ
		claims,
	)
	return token.SignedString(secret)
	// подписывает токен и превращает в строку, в заголовке HTTP-запроса можно передавать только текст
}

// ParseToken проверяет подпись JWT и извлекает из него Claims.
//
// Возвращает ошибку, если токен некорректен, подпись не совпадает,
// или срок действия истёк.
func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid { // булево поле в структуре jwt.Token, становится true, если подпись совпадает, токен не истёк и не был изменён
		// TODO: вынести ошибку в отдельный пакет
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
