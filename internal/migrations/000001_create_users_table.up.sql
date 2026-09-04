CREATE TABLE IF NOT EXISTS users (     /*Создаёт таблицу users, если её ещё нет*/
    id SERIAL PRIMARY KEY,             /*Целочисленный ID, автоинкремент, первичный ключ*/
    name TEXT NOT NULL,                /*Имя, не может быть пустым*/
    email TEXT UNIQUE NOT NULL,        /*Email, уникальный, не может быть пустым*/
    password TEXT NOT NULL,            /*Пароль, не может быть пустым */
    created_at TIMESTAMP DEFAULT NOW() /*Дата создания, автоматически заполняется текущим временем*/
);
/*Автоинкремент — это механизм, который автоматически увеличивает значение поля на 1 при каждой новой записи*/