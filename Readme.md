# Тестовое задание Junior Golang Developer
---
## Для удобного взаимодействия добавлен swagger по адресу localhost:3000/swagger/
## Access token. Тип - JWT
## Refresh token. Тип - произвольный
___
## Получение токенов - POST /create
## Обновление токенов - POST /refresh
## Refresh и Access токены приходят на клиент в json формате. Для обработки скопируйте access и refresh токены и вставьте в заголовок POST запроса с названиями Access и Refresh соответственно
___
## Запуск сервера
```
make run secret=<секртеный ключ>
```
## Запуск сервера через докер 
```
docker-compose up
```