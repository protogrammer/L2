# L2. CRUD

## Задание
Реализовать анонимный форум, на котором пользователь может писать, редактировать и удалять свои сообщения и смотреть чужие.

### Базовые возможности
- Добавление записей в общую ленту
- Добавление и удаление лайков

### Дополнительные реализованные возможности 
- Редактирование и удаление своих записей


# Ход работы

## Пользовательский интерфейс

## Пользовательские сценарии работы

### /
Главная страница. Содержит записи пользователй и кнопку для добавления записи

## Описание API сервера, хореографии

При любом запросе пользователю в `cookie` передаётся секретный ключ, обозначающий его идентификатор. Если ключ уже есть, он не обновляется. В базе он хранится в виде хэш-суммы `BLAKE2B`. Также при создании ключа создаётся случайное имя пользователя

### GET /api/get
Возвращает список постов в формате **JSON**
Каждый пост содержит следующие поля
 - **id** &ndash; идентификатор поста, число в десятичной системе счисления
 - **text** &ndash; текст сообщения
 - **author** &ndash; автор сообщения. Если пользователь является автором, то `null`
 - **created** &ndash; дата создания поста
 - **likes** &ndash; количество лайков
 - **liked** &ndash; поставил ли пользователь лайк, `true` или `false`
 
### GET /api/me
Возвращает юзернейм пользователя

### POST /api/new
Создание нового поста. В тело запроса передаётся текст
Возвращается сам пост в том формате, который описан в `GET /api/get`

### POST /api/edit
Редактирование существующего поста. В `URL` передаётся параметр **id**, в тело &ndash; текст. Редактировать чужой пост невозможно

### POST /api/delete
Удаление поста. В `URL` передаётся параметр **id**. Удалить чужой пост невозможно

## Структура бвзы данных


# Значимые фрагменты кода
