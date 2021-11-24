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
	Responses int `json:"responses,omitempty"`
	StartAt int64 `json:"startAt,omitempty" db:"startAt"`
	EndAt int64 `json:"endAt,omitempty" db:"endAt"`
}

type PollAnswer struct {
	PlayerID int64 `json:"playerID,omitempty" db:"playerID"`
	QuestionID int64 `json:"questionID" db:"questionID"`
	Answer string `json:"answer" db:"answer"`
}