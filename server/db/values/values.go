package values

import (
	"encoding/binary"
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"server/db/keys"
	"time"
)

type Post struct {
	PostId      uint64 `json:"id"`
	Text        string `json:"text"`
	AuthorId    string `json:"author"`
	TimeCreated string `json:"created"`
	Likes       uint64 `json:"likes"`
}

func (post *Post) Key() []byte {
	return keys.Post(post.PostId)
}

func PostsCountFromBytes(value []byte) uint64 {
	if len(value) != 8 {
		log.Panicln("[values.PostsCountFromBytes] incorrect value length", len(value))
	}
	return binary.BigEndian.Uint64(value)
}

func PostsCount(count uint64) []byte {
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, count)
	return value
}

func getParam[ParamType any](m map[string]any, param string) ParamType {
	valueAny, ok := m[param]
	if !ok {
		log.Panicln("[values.getParam] cannot find param", param)
	}
	value, ok := valueAny.(ParamType)
	if !ok {
		log.Panicf("[values.getParam] param %s has incorrect type %T (%v)", param, valueAny, valueAny)
	}
	return value
}

func PostFromBytes(id uint64, value []byte) *Post {
	m, ok := gjson.ParseBytes(value).Value().(map[string]interface{})
	if !ok {
		log.Panicln("[values.PostFromBytes] incorrect json", string(value))
	}

	return &Post{
		PostId:      id,
		Text:        getParam[string](m, "text"),
		AuthorId:    getParam[string](m, "author"),
		TimeCreated: getParam[string](m, "created"),
	}
}

func NewPost(id uint64, text string, author string) *Post {
	return &Post{
		PostId:      id,
		Text:        text,
		AuthorId:    author,
		TimeCreated: time.Now().Format(time.ANSIC),
		Likes:       0,
	}
}

func (post *Post) Value() ([]byte, error) {
	return json.Marshal(map[string]any{
		"text":    post.Text,
		"author":  post.AuthorId,
		"created": post.TimeCreated,
	})
}

func UsernameFromBytes(bytes []byte) string {
	return string(bytes)
}

func User(value string) []byte {
	return []byte(value)
}

func LikesFromBytes(bytes []byte) uint64 {
	if len(bytes) != 8 {
		log.Panicln("[values.LikesFromBytes] incorrect value length", len(bytes))
	}
	return binary.BigEndian.Uint64(bytes)
}

func Likes(count uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, count)
	return bytes
}
