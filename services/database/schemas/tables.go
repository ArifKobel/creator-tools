package schemas

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Email  string  `json:"email"`
	Videos []Video `json:"videos"`
	Otp    string  `json:"otp"`
}
type JSONB map[string]interface{}

func (j JSONB) Value() (string, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type Export struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	VideoID uuid.UUID `json:"video_id"`
	URL     string    `json:"url"`
}
type Video struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Filepath      string          `json:"filepath"`
	UserID        uint            `json:"user_id"`
	Filename      string          `json:"filename"`
	Duration      float64         `json:"duration"`
	ThumbnailPath string          `json:"thumbnail_path"`
	CreatedAt     string          `json:"created_at"`
	Subtitles     json.RawMessage `json:"subtitles" gorm:"type:jsonb"`
	Exports       []Export        `json:"exports"`
}
