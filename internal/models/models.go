package models

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Nominant struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type NominantCategory struct {
	NominantID int64 `json:"nominant_id"`
	CategoryID int64 `json:"category_id"`
}

type User struct {
	TGID int64 `json:"tg_id"`
}

type Vote struct {
	TGUserID   int64 `json:"tg_user_id"`
	NominantID int64 `json:"nominant_id"`
	CategoryID int64 `json:"category_id"`
}

type VoteRequest struct {
	NominantID int64 `json:"nominant_id"`
	CategoryID int64 `json:"category_id"`
}
