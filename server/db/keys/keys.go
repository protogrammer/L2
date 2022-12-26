package keys

import (
	"encoding/binary"
	"server/db/prefix"
)

func PostsCount() []byte {
	return []byte{prefix.PostsCount}
}

func Post(id uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, id)
	return append([]byte{prefix.Post}, key...)
}

func User(userPublicKey []byte) []byte {
	return append([]byte{prefix.User}, userPublicKey...)
}

func Likes(id uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, id)
	return append([]byte{prefix.Likes}, key...)
}

func Liked(userPublicId []byte, postId uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, postId)
	return append(append([]byte{prefix.Liked}, userPublicId...), Post(postId)...)
}
