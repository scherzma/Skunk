package c_model

type Message struct {
	id        string
	timestamp string
	content   string
	from      User // should probably use User ID instead of a direct reference
	to        User // should probably use User ID instead of a direct reference
}
