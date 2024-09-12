package entities

type URL struct {
	ID        int64  `json:"-" gorm:"primaryKey;autoIncrement:true"`
	SortURL   string `json:"sort_url" gorm:"uniqueIndex"`
	LongURL   string `json:"long_url" gorm:"uniqueIndex"`
	CreatedAt int64  `json:"-"`
}
