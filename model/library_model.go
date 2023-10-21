// model/library_model.go

package model

import (
	"Reference/database"
	"fmt"
)

type Library struct {
	Id     uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name   string `json:"name"`
	UserId uint   `gorm:"not null" json:"user_id"`
}

func (l *Library) Save() (*Library, error) {
	err := database.DB.Create(&l).Error
	if err != nil {
		return &Library{}, err
	}
	return l, err
}

func GetLibrary(id int) (Library, error) {
	var library Library
	err := database.DB.Debug().Where("id = ?", id).First(&library).Error
	if err != nil {
		return Library{}, err
	}
	fmt.Println("================", library)
	return library, nil
}

func FindLibraryByName(name string) (Library, error) {
	var library Library
	err := database.DB.Where("name = ?", name).First(&library).Error
	if err != nil {
		return Library{}, err
	}
	return library, nil
}

func UpdateLibrary(library *Library) (err error) {
	err = database.DB.Save(library).Error
	if err != nil {
		return err
	}
	return nil
}
