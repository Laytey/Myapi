package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("jwt-secret-word")

// глобальная переменная для пакета,
// []byte превращает строку в слайс байтов
// "jwt-secret-word" ключ
type Claims struct {
	UserID               string `json:"user_id"` // тег: при превращении в JSON назовет это поле user_id
	jwt.RegisteredClaims        // структура из пакета jwt, стандартные поля: exp, iat, iss, aud
}

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
		jwt.SigningMethodHS256, // алгоритм подписи токена с использованием общего секрета (симметричное шифрование) хз что это
		// симмитричное шифрование - для шифровки и расшифровки нужен один и тот же ключ
		claims,
	)
	return token.SignedString(secret)
	// подписывает токен и превращает в строку, в заголовке HTTP-запроса можно передавать только текст
}

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
