package model

type Property struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Type      string `json:"type"`
}
