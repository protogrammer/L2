package db

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/dgraph-io/badger/v3"
	"golang.org/x/crypto/blake2b"
	"server/db/keys"
	"server/db/values"
	"server/words"
)

type User struct {
	id []byte
}

func NewUser() (user *User, secretKey string, err error) {
	secretKeyBytes := make([]byte, 32)
	_, err = rand.Read(secretKeyBytes)
	if err != nil {
		return
	}

	secretKey = base64.RawURLEncoding.EncodeToString(secretKeyBytes)
	public := blake2b.Sum256(secretKeyBytes)
	user = &User{
		id: public[:],
	}

	userKey := user.Key()

	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(userKey)
		return err
	})

	if errors.Is(err, badger.ErrKeyNotFound) {
		err = db.Update(func(txn *badger.Txn) error {
			return txn.Set(userKey, []byte(words.AwesomeUsername()))
		})

		return
	}

	if err != nil {
		return
	}

	return NewUser()
}

func UserBySecretKey(secretKey string) (*User, error) {
	data, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, err
	}
	id := blake2b.Sum256(data)
	return &User{
		id: id[:],
	}, nil
}

func UserByPublicKey(publicKey string) (*User, error) {
	data, err := base64.RawURLEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	return &User{
		id: data,
	}, nil
}

func (user *User) Key() []byte {
	return keys.User(user.id)
}

func (user *User) PublicIdString() string {
	return base64.RawURLEncoding.EncodeToString(user.id)
}

func (user *User) FindUsername() (username string, err error) {
	err = db.Update(func(txn *badger.Txn) error {
		usernameItem, err := txn.Get(user.Key())
		if err != nil {
			return err
		}
		_ = usernameItem.Value(func(val []byte) error {
			username = values.UsernameFromBytes(val)
			return nil
		})
		return err
	})
	return
}

func (user *User) Like(post *Post, value bool) (uint64, bool) {
	likedKey := keys.Liked(user.id, post.id)
	likesCountKey := keys.Likes(post.id)
	likes := uint64(0)
	err := db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(likedKey)

		if errors.Is(err, badger.ErrKeyNotFound) && value {
			err = txn.Set(likedKey, nil)
			if err != nil {
				return err
			}

			likesItem, err := txn.Get(likesCountKey)
			if err == nil {
				_ = likesItem.Value(func(val []byte) error {
					likes = values.LikesFromBytes(val)
					return nil
				})
			}

			likes++

			return txn.Set(likesCountKey, values.Likes(likes))
		}

		if err == nil && !value {
			err = txn.Delete(likedKey)
			if err != nil {
				return err
			}

			likesItem, err := txn.Get(likesCountKey)
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					return nil
				}
				return err
			}

			_ = likesItem.Value(func(val []byte) error {
				likes = values.LikesFromBytes(val)
				return nil
			})

			if likes > 0 {
				likes--
			}

			return txn.Set(likesCountKey, values.Likes(likes))
		}

		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}

		return err
	})

	return likes, err == nil
}

func (user *User) Liked(post *Post) bool {
	likedKey := keys.Liked(user.id, post.id)

	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(likedKey)
		return err
	})

	return err == nil
}
