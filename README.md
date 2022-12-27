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
В качестве СУБД используется **BadgerDB**
Под идентификатором пользователя обозначается хэшсумма `BLAKE2B` его секретного ключа

### Post
По идентификатру поста (который имеет тип `uint64`) возвращается его текст, дата создания и идентификатор автора

### User
По идентификатору пользователя возвращается его юзернейм

### Likes
По идентификатору поста возвращается счётчик лайков

### Liked
Содержит идентификаторы пользователя и поста.


# Значимые фрагменты кода

## 1. Создание поста
```go
router.POST("/api/new", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			 return
		}

		textBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			 c.AbortWithStatus(http.StatusInternalServerError)
			 log.Println("[server] POST /api/new, error reading request body:", err)
			 return
		}

		post, err := db.NewPost(user, string(textBytes))
		if err != nil {
			 c.AbortWithStatus(http.StatusInternalServerError)
			 log.Println("[server] POST /api/new, error creating new post:", err)
			 return
		}

		c.JSON(http.StatusOK, gin.H{
			 "id":      post.PostId,
			 "text":    post.Text,
 			"author":  nil,
			 "created": post.TimeCreated,
	 		"likes":   post.Likes,
	 		"liked":   false,
		})
})
```

## Редактирование поста
```go
router.POST("/api/edit", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			 return
		}

		postIdString := c.Request.URL.Query().Get("id")
		postId, err := strconv.ParseUint(postIdString, 10, 64)
		if err != nil {
	 		c.AbortWithStatus(http.StatusExpectationFailed)
 			log.Println("[server] POST /api/edit, non-integer postId:", err)
 			return
		}

		textBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
	 		c.AbortWithStatus(http.StatusInternalServerError)
	 		log.Println("[server] POST /api/edit, error reading request body:", err)
	 		return
		}

		post := db.PostById(postId)

		if !post.WasCreatedBy(user) {
 			c.AbortWithStatus(http.StatusExpectationFailed)
 			log.Println("[server] POST /api/edit, trying to edit the post created by someone else")
	 		return
		}

		if !post.Edit(string(textBytes)) {
		 	c.AbortWithStatus(http.StatusExpectationFailed)
		 	log.Println("[server] POST /api/edit, edit failed")
		 	return
		}

		c.Data(http.StatusOK, "text/plain", textBytes)
})
```

## Удаление поста
```go
router.POST("/api/delete", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
		 	return
		}

		postIdString := c.Request.URL.Query().Get("id")
		postId, err := strconv.ParseUint(postIdString, 10, 64)
		if err != nil {
			 c.AbortWithStatus(http.StatusExpectationFailed)
		 	log.Println("[server] POST /api/delete, non-integer postId:", err)
		 	return
		}

		post := db.PostById(postId)

		if !post.WasCreatedBy(user) {
		 	c.String(http.StatusOK, "false")
			 log.Println("[server] POST /api/delete, trying to edit the post created by someone else")
			 return
		}

		ok = post.Delete()

		c.String(http.StatusOK, strconv.FormatBool(ok))
})
```

## Получить все посты
```go
router.GET("/api/get", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			 return
		}

		var jsonTags []any

		f := func(postId uint64) {
		 	post := db.PostById(postId)
		 	realPost, ok := post.Find()
		 	if !ok {
	 	 		return
	 		}
 
		 	other, _ := db.UserByPublicKey(realPost.AuthorId)
 			var username any = nil
 			if user.PublicIdString() != other.PublicIdString() {
	  			username, _ = other.FindUsername()
		 	}

	 		jsonTags = append(jsonTags, gin.H{
		  		"id":      realPost.PostId,
		  		"text":    realPost.Text,
		  		"author":  username,
		  		"created": realPost.TimeCreated,
		  		"likes":   realPost.Likes,
		  		"liked":   user.Liked(post),
 			})
		}

		id := db.PostsCount()

		if id == 0 {
 			c.JSON(http.StatusOK, jsonTags)
	 		return
		}

		for {
	 		id--

	 		f(id)

 			if id == 0 {
	  			break
	 		}
		}

		c.JSON(http.StatusOK, jsonTags)
})
```

## Получить свой юзернейм
```go
router.GET("/api/me", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
 			return
		}

		username, err := user.FindUsername()
		if err != nil {
	 		c.AbortWithStatus(http.StatusInternalServerError)
	 		log.Println("[server] GET /api/me, error finding username:", err)
	 		return
		}

		c.String(http.StatusOK, username)
})
```

## Поставить/снять лайк
```go
router.POST("/api/like", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			 return
		}

		postIdString := c.Request.URL.Query().Get("id")
		postId, err := strconv.ParseUint(postIdString, 10, 64)
		if err != nil {
 			c.AbortWithStatus(http.StatusExpectationFailed)
	 		log.Println("[server] POST /api/like, non-integer postId:", err)
		 	return
		}

		valueString := c.Request.URL.Query().Get("value")
		value, err := strconv.ParseBool(valueString)
		if err != nil {
 			c.AbortWithStatus(http.StatusExpectationFailed)
	 		log.Println("[server] POST /api/like, non-boolean value:", err)
		 	return
		}

		post := db.PostById(postId)

		count, ok := user.Like(post, value)
		if !ok {
 			c.AbortWithStatus(http.StatusInternalServerError)
	 		log.Println("[server] POST /api/like, like error:", err)
		 	return
		}

		c.String(http.StatusOK, strconv.FormatUint(count, 10))
})
```

## Создание нового пользователя
```go
func updateCookie(c *gin.Context) (user *db.User, ok bool) {
 	ok = true
	 secretKey, err := c.Cookie("secretKey")
	 if err != nil {
  		if !errors.Is(err, http.ErrNoCookie) {
	   		log.Panicln("[server.updateCookie] error getting cookie:", err)
  		}
  		user, secretKey, err = db.NewUser()
	  	if err != nil {
	   		log.Panicln("[server.updateCookie] error creating user:", err)
	  	}

	  	c.SetCookie("secretKey", secretKey, 0, "/", env.Domain(), false, true)
  		return
  }
 
	 user, err = db.UserBySecretKey(secretKey)
	 if err != nil {
	  	c.AbortWithStatus(http.StatusBadRequest)
  		return nil, false
 	}

	 return
}
```
