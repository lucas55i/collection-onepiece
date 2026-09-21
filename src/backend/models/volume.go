package models

import "time"

// Volume representa um volume do mangá One Piece
type Volume struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	VolumeNumber int        `json:"volume_number" gorm:"uniqueIndex;not null"`
	Title        string     `json:"title"`
	CoverImage   string     `json:"cover_image"`
	Chapters     string     `json:"chapters"` // ex: "1-7"
	Collected    bool       `json:"collected" gorm:"default:false"`
	AcquiredAt   *time.Time `json:"acquired_at" gorm:"default:null"` // data/hora em que foi adquirido
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
