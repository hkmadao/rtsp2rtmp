package entity

import (
	"time"
)

type CameraShare struct {
	Id        string    `orm:"pk;column(id)" json:"id"`
	Name      string    `orm:"column(name)" json:"name"`
	AuthCode  string    `orm:"column(auth_code)" json:"authCode"`
	Enabled   int       `orm:"column(enabled)" json:"enabled"`
	Created   time.Time `orm:"column(created)" json:"created"`
	StartTime time.Time `orm:"column(start_time)" json:"startTime"`
	Deadline  time.Time `orm:"column(deadline)" json:"deadline"`
	CameraId  string    `orm:"column(camera_id)" json:"cameraId"`
	Camera    Camera    `orm:"-" json:"camera"`
}
