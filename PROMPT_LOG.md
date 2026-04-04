# Prompt Log
## Задание Средней сложности 1: Создать простое API на Go (Gin) с 2-3 эндпоинтами.
### Промпт 1
**Инструмент:** Auto режим в Cursor.
**Промпт:** "Work in /mid/1.go. Create a simple API on Go (Gin) with 2 endpoints."
**Результат:** Ожидаемый результат, два простых рабочих эндпоинта. Запустил с go run 1.go
### Итого
- Количество промптов: 1
- Что пришлось исправлять вручную: изменил названия эндпоинотов, вывод сообщения
- Время: ~ 5 минут
---
## Задание Средней сложности 2: Реализовать валидацию входных данных в Go.
### Промпт 1
**Инструмент:** Auto режим в Cursor.
**Промпт:** "Now we need to implement a go input validation. Create a separate folder, name it "3". Use code from 1.go, add an email validation, with check of age and name length"
**Результат:** Рабочая валидация, проврил POST - curl -X POST http://localhost:8080/profile -H "Content-Type: application/json" -d '{"display_name": "T", "email": "not-an-email", "age": 200}' {"error":"Key: 'profileBody.DisplayName' Error:Field validation for 'DisplayName' failed on the 'min' tag\nKey: 'profileBody.Email' Error:Field validation for 'Email' failed on the 'email' tag\nKey: 'profileBody.Age' Error:Field validation for 'Age' failed on the 'lte' tag"}% - ожидаемый результат, с правильными данными ошибки не наблюдается.
### Промпт 2
**Инструмент:** Auto режим в Cursor.
**Промпт:** "Now, create a go-test file for 3.go. It should check GET for /hello/:name, POST for /profile."
**Результат:** Рабочий go-test, проверяет все нужные эндпоинты. Запуск с помощью "go test -v" выдает ожидаемый результат.
### Итого
- Количество промптов: 2
- Что пришлось исправлять вручную: Агент криво завершил работу сервера, стандартный порт не освобождался. Пришлось вручную искать процесс по занятому порту и выключать его. Помимо этого, все работает в соответсвии с запросом.
- Время: ~ 15 минут
---
## Задание Средней сложности 3: Передавать сложные структуры данных (JSON) между сервисами.
### Промпт 1
**Инструмент:** Auto режим в Cursor.
**Промпт:** "Create a folder inside /mid, name it "5". Create a FastAPI endpoint that receives a 'Profile' (from 3.go) JSON object and forwards it to a Go service at http://localhost:8080/profile using the httpx library. Use Pydantic for validation on the Python side to match the Go struct."
**Результат:** main.py и requirements.txt. Запустил FastAPI сервис на 8000 порту - uvicorn main:app --reload --port 8000, после чего подаю POST запрос на этот же порт (такой же запрос с email и именем). Получаю сообщение - INFO: 127.0.0.1:54417 - "POST /profile HTTP/1.1" 201 Created от python, а так же в командной строке где запущен go - "GIN] 2026/04/04 - 17:08:47 | 201 |        87.5µs |             ::1 | POST     "/profile"". Полностью рабочая передача JSON.
### Промпт 2
**Инструмент:** Auto режим в Cursor.
**Промпт:** "Now, create a test for main.py. It should send two POST's, one valid and other not."
**Результат:** файл test_main.py. Полностью проходят два теста, один успешный второй с возвращением ошибки 422.
### Итого
- Количество промптов: 
- Что пришлось исправлять вручную: Pip не нашел pytest из-за Homebrew, агент предложил создать .sh скрипт для использования локального .venv проекта. Решил проблему установив pytest глобально.
- Время: ~ 15 минут.
---
## Задание Повышенной сложности 1: Добавить аутентификацию (JWT) в Go-сервисе и проверять токены из Python.
### Промпт 1
**Инструмент:** Auto режим в Cursor.
**Промпт:** "Now, work in /hard. In there, create two folders, call them python-service and go-service. We need to turn task 5 on python and task 3 in go into little services. Go(gin) as a backend and Python as a gateway"
**Результат:** Просто скопировал и переметил файлы для более удобной работы. Вся работоспособность остается такой же
### Промпт 2
**Инструмент:** Auto режим в Cursor.
**Промпт:** "Now only work in folders /hard, /go-service and /python-service. Add a JWT authentification in go service and check tockens from Python"
**Результат:** Рабочая проверка токенов. Использовался запрос из задания 3. Для доступа к эндпоинту POST /profile теперь требуется JWT-токен, который можно получить через эндпоинт /auth/token в Go-сервисе. При непрпавльном запросе или токене, python не позволит сделать запрос на Go сервис.
### Промпт 3
**Инструмент:** Auto режим в Cursor.
**Промпт:** "create a /cmd adn /internal for go-service, sort go service with new file structure. Organize go-service"
**Результат:** Переорганизация файлов. Теперь гораздо удобнее и логичнее ориентироваться в проекте.
### Промпт 4
**Инструмент:** Auto режим в Cursor.
**Промпт:** ""
**Результат:** 
### Итого
- Количество промптов: 3
- Что пришлось исправлять вручную: 
- Время: ~ 20 минут
---
## Задание Повышенной сложности 2: Развернуть оба сервиса в Docker Compose с общей сетью..
### Промпт 1
**Инструмент:** Auto режим в Cursor.
**Промпт:** ""
**Результат:** 
### Итого
- Количество промптов: 
- Что пришлось исправлять вручную:
- Время:
---
