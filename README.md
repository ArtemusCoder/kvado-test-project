# Тестовое задание

## Запуск через Makefile
Понадобится установленный Docker, Golang

1) <code>make build</code>
2) <code>make run_db</code>
3) <code>make run_server</code> - в новом терминале

## Остановка через Makefile

1) Ctrl^C для терминала с запущенным сервером
2) <code>make stop_db</code> - остановка контейнера

## Очистка через Makefile
<code> make clean </code>

## Unit Tests
Тесты проходят с запущенным Докером <code>make run_db </code> <br />
<code>make tests</code> - запуск тестов

## Проверка работы в ручную 
После запуска через Makefile можно воспользоваться Postman, grpcurl или другим сервисом, позволяющим делать grpc-соединения
![screen](https://github.com/ArtemusCoder/kvado-test-project/assets/33132419/c757b0e9-3cc2-41c1-b71e-47ef0e085b57)

GetBooksByAuthor - поиск книг по автору <br />
GetAuthorsByBook - поиск авторов по книге

Сервис будет доступен по адресу <code>localhost:8080</code>
