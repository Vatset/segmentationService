#  (ENG) REST API Service for Dynamic User Segmentation

## Description
This service provides the following capabilities:

- Registration of new users
- Creation, deletion, and updating of segments
- Addition and removal of users from segments
- Adding segments with limited time validity
- Creating auto-adding segments for a certain percentage of registered users
- Retrieving the history of segment additions and deletions on a user's request for a year-month in the form of a link to an automatically downloadable CSV file

### Used Tools and Technologies
- Golang
- Gin Web Framework
- Docker
- PostgreSQL
- Swagger

## Questions That Arise During the Project Implementation:
1. What principle should be used to determine the percentage of users who will be auto-segmented?<br>
 For example, if a value of 10% is provided as input, and segments are assigned to 10% of users starting from the beginning of the table, then the oldest users will always have more segments. To avoid this situation, I decided to use a pseudo-random approach to select users to be added to the segment.

2. Can a segment be duplicated for a user?<br>
 Considering the examples of segments in the requirements, such as "AVITO_DISCOUNT_50," there were doubts, as in real life, a user can have more than one discount. However, logically reasoning and giving a definition to the word "segment," I came to the conclusion that a segment cannot be duplicated. It is similar to a user belonging to the same community at the same time.

3. How should the link to the CSV file be displayed? Should it be forwarded to a service for viewing CSV files, or should it be implemented differently?<br>
 I chose the following method:<br>
The method is called, and if the folder for storing CSV files does not exist, it is created in the root directory with the name "history." Then a CSV file is generated and saved in this folder.<br>
As a result, the user receives a link to the file, for example: "http://localhost:8080/api/user/history/user1_period2023-08_history.csv," and when clicked on this link in any browser, the CSV file is automatically downloaded.

#  (RU) REST API сервис для динамического сегментирования пользователей

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

# Running the Service // Запуск сервиса
1. Clone the repository // Клонируйте репозиторий
```bash
 git clone https://github.com/Vatset/segmentationService
```
2. Create a .env file // Создайте .env файл
  Example // Пример:
```bash
DB_PASSWORD=yourdbpass
```
3. Prepare the database for operation // Подготовка бд к работе<br>
   *Download and run the Docker application beforehand* // *Предварительно скачайте и запустите приложение docker*<br>
Obtain the latest version of PostgreSQL//Получение последней версии postgres
```bash   
docker pull postgres
```
Run a Docker container named "segmentation_service" using the previously downloaded PostgreSQL. // Запуск Docker контейнера с именем "segmentation_service", используя ранее скачанный образ PostgreSQL. 
```bash
docker run --name=segmentation_service -e POSTGRES_PASSWORD="yourdbpass" -p 5436:5432 -d --rm postgres
```
Execute database migrations // Выполнение миграций базы данных
```bash 
migrate -path ./schema  -database 'postgres://postgres:yourdbpass@localhost:5436/postgres?sslmode=disable' up
```
4.Launch the project // Запускаем проект
```bash   
go run cmd/main.go
```

## Examples of Requests and Responses // Примеры запросов и ответов
After launching the project, the Swagger UI for the service will be available at http://localhost:8080/swagger/index.html // После запуска проекта по адресу http://localhost:8080/swagger/index.html будет доступен Swagger UI сервиса

### Create User // Создание пользователя [POST]
```bash   
/api/user/create
```
Request // Запрос
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
### Delete User // Удаление пользователя [DELETE]
```bash   
/api/user/delete
```
Request // Запрос
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
### Display User Segments // Отображение сегментов пользователя [GET]
```bash   
api/user/showSegments/1
```
```bash   
{
    "User Segments": "AVITO_CHAT,AVITO_VOICE_MESSAGES"
}
```
### Get Link to User's Segment History // Получение ссылки на историю сегментов пользователя [GET]
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
By clicking on this link in a web browser, you can download the CSV file of the user's segmentation history for the specified period.
Перейдя по данной ссылке в браузере можно скачать csv файл истории сегментации пользователя за даннный период

### Create Segment // Создание сегмента [POST]
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
### For Auto-Adding a Segment to a Certain Percentage of Users // Для автодобавления сегмента определенному проценту пользователей
Request // Запрос
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
### Delete Segment // Удаление сегмента [DELETE]
```bash   
/api/segment/delete
```
Request // Запрос
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
### Update Segment // Обновление сегмента [PUT]
```bash   
/api/segment/update
```
Request // Запрос
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
### Segmentation // Сегментация [POST]
```bash   
/api/segment/update
```
Request // Запрос
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


