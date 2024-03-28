package schemas

type User struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Email  string  `json:"email"`
	Videos []Video `json:"videos"`
	Otp    string  `json:"otp"`
}

type Video struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	FilePath string `json:"file_path"`
}
