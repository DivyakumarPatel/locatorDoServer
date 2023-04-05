package myStructs

import "time"

type User struct {
	Id        int `json:"id" gorm:"primaryKey"`
	Email         string `json:"email" binding:"required"`
	First_name    string `json:"first_name" binding:"required"`
	Middle_name   string `json:"middle_name" binding:"required"`
	Phone_number  string `json:"phone_number" binding:"required"`
	Firebase_id   string `json:"firebase_id,omitempty" `
	Password      string `json:"-" `
	Profile_photo string `json:"profile_photo" default:"https://www.pngitem.com/pimgs/m/30-307416_profile-icon-png-image-free-download-searchpng-employee.png"`
	Is_admin      bool   `json:"is_admin" `
}

type LoginUser struct {
	Email         string `json:"email" binding:"required"`

	Password      string `json:"-" `
	
}



type LoginData struct {
	Email       string `json:"email" `
	Password    string `json:"Password" `
	Firebase_id string `json:"firebase_id" `
}

type FcmModel struct {
	FirstName          string
	LastName           string
	IsNotificationSent bool
}
type MyLocation struct {
	OriginLatitude    string    `json:"origin_latitude" `
	OriginLongitude   string    `json:"origin_longitude"`
	MaxDistance       string    `json:"max_distance" `
}

type LocationUpdate struct {
	User_distance     float32   `json:"user_distance"`
	Notification_sent string    `json:"notification_sent" `
	CurrentLatitude   string    `json:"current_latitude" `
	CurrentLongitude  string    `json:"current_longitude" `
	UserId            string    `json:"user_id" `
	OriginLatitude    string    `json:"origin_latitude" `
	OriginLongitude   string    `json:"origin_longitude"`
	MaxDistance       float32    `json:"max_distance" `
	LastUpdate        time.Time `json:"time" `
	Email             string    `json:"email" `
	FirstName         string    `json:"first_name" `
	MiddleName        string    `json:"middle_name" `
	PhoneNumber       string    `json:"phone_number"`
}
