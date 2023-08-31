# REST API сервис для динамического сегментирования пользователей

## Описание

Данный сервис предоставляет возможность:

- регистрации новых пользователей
- создания, удаления и обновления сегментов
- добавление и удаление пользователей в сегмент
- добавление сегментов с ограниченным временем действия
- создание автодобавляющихся сегментов для определенного процента зарегистрированных пользователей
- получение истории добавления и удаления сегментов по запросу год-месяц для пользователя в виде ссылки на автоматически скачиваемый csv файл

### Используемые инструменты и технологии

* Golang
* Gin Web Framework
* Docker
* PostgreSQL
* Swagger

## Запуск сервиса
1. Клонируйте репозиторий
```bash
 git clone https://github.com/Vatset/segmentationService
```
2. Создайте .env файл
   Пример:
```bash
DB_PASSWORD=yourdbpass
```
3. Подготовка бд к работе<br>
   *Предварительно скачайте и запустите приложение docker*<br>
Получение последней версии postgres
```bash   
docker pull postgres
```
Запуск Docker контейнера с именем "segmentation_service", используя ранее скачанный образ PostgreSQL. 
```bash
docker run --name=segmentation_service -e POSTGRES_PASSWORD="yourdbpass" -p 5436:5432 -d --rm postgres
```
Выполнение миграций базы данных
```bash 
migrate -path ./schema  -database 'postgres://postgres:yourdbpass@localhost:5436/postgres?sslmode=disable' up
```
4.Запускаем проект
```bash   
go run cmd/main.go
```

## Примеры запросов и ответов
После запуска проекта по адресу http://localhost:8080/swagger/index.html будет доступен Swagger UI сервиса

### Создание пользователя [POST]
```bash   
/api/user/create
```
Запрос
```bash   
{
    "username":"Joe"
}
```
```bash   
{
  "id": 1
}
```
### Удаление пользователя [DELETE]
```bash   
/api/user/delete
```
Запрос
```bash   
{
    "username":"Joe"
}
```
```bash   
{
    "status": "user was successful deleted"
}
```
### Отображение сегментов пользователя [GET]
```bash   
api/user/showSegments/1
```
```bash   
{
    "User Segments": "AVITO_CHAT,AVITO_VOICE_MESSAGES"
}
```
### Получение ссылки на историю сегментов пользователя [GET]
```bash   
api/user/historyLink/1
```
```bash   
{
  "timestamp": "2023-08"
}
```
```bash   
{
    "Link": "http://localhost:8080/api/user/history/user1_period2023-08_history.csv"
}
```
Перейдя по данной ссылке в браузере можно скачать csv файл истории сегментации пользователя за даннный период

### Создание сегмента [POST]
```bash   
/api/segment/create
```
Запрос
```bash   
{
    "segment":"AVITO_VOICE_MESSAGES"
}
```
```bash   
{
    "id": 1
}
```
### Для автодобавления сегмента определенному проценту пользователей
Запрос
```bash   
 {   
    "segment":"AVITO_CHAT",
    "percent":50
 }
```
```bash   
{
    "id": 2
}
```
### Удаление сегмента [DELETE]
```bash   
/api/segment/delete
```
Запрос
```bash   
 {   
    "segment":"AVITO_CHAT"
 }
 ```
```bash   
{
    "status": "segment was successful deleted"
}
```
### Обновление сегмента [PUT]
```bash   
/api/segment/update
```
Запрос
```bash   
 {   
    "last_name":"AVITO_VOICE_MESSAGES",
    "new_name":"AVITO_VIDEO_MESSAGES"
 }
```
```bash   
{
    "status": "segment was successful updated"
}
```
### Сегментация [POST]
```bash   
/api/segment/update
```
Запрос
```bash   
 {   
    "last_name":"AVITO_VOICE_MESSAGES",
    "new_name":"AVITO_VIDEO_MESSAGES"
 }
```
```bash   
{
    "status": "segment was successful updated"
}
```
### Вопросы, которые появились в процессе реализации проекта:
1. По какому принципу определенный процент пользователей будет автосегментироваться?<br>
   Например, если на вход поступает значение 10%, и при этом сегменты присваиваются для 10% пользователей начиная с начала таблицы, то у самых старых пользователей всегда будет больше сегментов. Для избежания этой ситуации я решшила использовать псевдослучайный подход для выбора пользователей, которые будут добавляться в сегмент
   
2. Может ли сегмент у пользователя дублироваться?<br>
   Рассматривая примеры сегментов в ТЗ, и видя такие сегменты как "AVITO_DISCOUNT_50", возникли сомнения, так как в реальной жизни у пользователя может быть более одной скидки.  Но рассуждая логически и давая определение слову сегмент я пришла к выводу,  что сегмент не может дублироваться. Это тоже самое если бы пользователь относился к одному и тому же сообществу одновременно.

3. Как должна отображаться ссылка на CSV файл? Должна ли она пересылать на какой-то сервис для просмотра CSV файлов или это должно быть реализовано иначе?<br>
   Я выбрала следующий способ:<br>
   Вызывается метод,если папки для хранения CSV файлов не существует,то  она создается в корневой директории с именем "history".Затем генерируется CSV файл и сохраняется в этой папке.<br>
   В результате пользователь получает ссылку на файл, например: "http://localhost:8080/api/user/history/user1_period2023-08_history.csv", при переходе по которой в любом браузере данный CSV файл автоматически скачивается.

