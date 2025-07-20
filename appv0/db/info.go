package db

import (
	"fmt"
	"time"
)

type Info struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func (i *Info) Table() string {
	return "info"
}

func (i *Info) Get(id int) *Info {
	var ret Info
	if err := DB.Table(i.Table()).Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &ret
}

func (i *Info) Save(id int, name string) error {
	return DB.Table(i.Table()).Where("id = ?", id).Update("name", name).Error
}
