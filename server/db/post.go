package db

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
	"server/db/keys"
	"server/db/values"
)

type Post struct {
	id uint64
}

func PostById(id uint64) *Post {
	return &Post{
		id: id,
	}
}

func (post *Post) Id() uint64 {
	return post.id
}

func NewPost(author *User, text string) (post *values.Post, err error) {
	err = db.Update(func(txn *badger.Txn) error {
		postsCountItem, err := txn.Get(keys.PostsCount())
		if err != nil {
			log.Panicln("[db.NewPost] error getting posts count")
		}

		var postId uint64
		_ = postsCountItem.Value(func(val []byte) error {
			postId = values.PostsCountFromBytes(val)
			return nil
		})

		err = txn.Set(keys.PostsCount(), values.PostsCount(postId+1))
		if err != nil {
			log.Panicln("[db.NewPost] error updating posts count")
		}

		post = values.NewPost(postId, text, author.PublicIdString())
		postValue, err := post.Value()
		if err != nil {
			return err
		}
		return txn.Set(post.Key(), postValue)
	})

	return
}

func (post *Post) Find() (*values.Post, bool) {
	var realPost *values.Post
	postKey := keys.Post(post.id)
	err := db.View(func(txn *badger.Txn) error {
		value, err := txn.Get(postKey)
		if err != nil {
			return err
		}
		err = value.Value(func(val []byte) error {
			realPost = values.PostFromBytes(post.id, val)
			return nil
		})
		if err != nil {
			return err
		}
		likesItem, err := txn.Get(keys.Likes(post.id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		return likesItem.Value(func(val []byte) error {
			realPost.Likes = values.LikesFromBytes(val)
			return nil
		})
	})
	if err == nil {
		return realPost, true
	}
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, false
	}
	log.Panicln("[db.FindPost] error", err)
	return nil, false // unreachable, just to compile
}

func (post *Post) Edit(newText string) bool {
	postKey := keys.Post(post.id)
	err := db.Update(func(txn *badger.Txn) error {
		value, err := txn.Get(postKey)
		if err != nil {
			return err
		}
		var newValue []byte
		err = value.Value(func(val []byte) (err error) {
			newValue, err = sjson.SetBytes(val, "text", newText)
			return err
		})
		if err != nil {
			return err
		}
		return txn.Set(postKey, newValue)
	})
	return err == nil
}

func (post *Post) Delete() bool {
	postKey := keys.Post(post.id)
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete(postKey)
	})
	return err == nil
}

func (post *Post) WasCreatedBy(user *User) (result bool) {
	key := keys.Post(post.id)
	_ = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			creator := gjson.GetBytes(val, "author")
			result = creator.String() == user.PublicIdString()
			return nil
		})
	})
	return
}

func PostsCount() (postsCount uint64) {
	_ = db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get(keys.PostsCount())
		_ = item.Value(func(val []byte) error {
			postsCount = values.PostsCountFromBytes(val)
			return nil
		})
		return nil
	})
	return
}
