package models

type PollQuestion struct {
	QuestionID int64 `json:"questionID" db:"questionID"`
	Question string `json:"question"`
	Options string `json:"options"`
	Required bool `json:"required"`
}

type Poll struct {
	ID int64 `json:"id" db:"ID"`
	Name string `json:"name"`
	Description string `json:"description"`
	Questions []PollQuestion `json:"questions"`
}

type PollAnswer struct {
	QuestionID int64 `json:"questionID" db:"questionID"`
	Answer string `json:"answer" db:"answer"`
}