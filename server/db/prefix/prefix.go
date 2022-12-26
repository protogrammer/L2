package prefix

const (
	PostsCount byte = iota // uint64
	Post                   // id uint64 => json { text: string, created: string, authorId: string }
	User                   // secretKeyBlake [32]byte => username string
	Likes                  // id uint64 => likeCount uint64
	Liked                  // UserId & PostId => ()
)
