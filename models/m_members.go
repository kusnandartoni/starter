package models

import (
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

// Members :
type Members struct {
	Base `mapstructure:",squash"`

	Email    string `json:"email" gorm:"type:varchar(100);unique_index"`
	Password string `json:"-"`
	FullName string `json:"full_name"`
	PhotoURL string `json:"photo_url"`
	Verified bool   `json:"-"`
	UUID     string `json:"uuid"`
}

// GetMembers :
func GetMembers(pageNum, pageSize int, maps interface{}) ([]Members, error) {
	var (
		members []Members
		err     error
	)

	dbMember := db.Where(maps)
	if pageSize > 0 && pageNum >= 0 {
		dbMember = dbMember.Offset(pageNum).Limit(pageSize)
	}

	err = db.Order("id desc").Find(&members).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return members, nil
}

// AddMembers : insert new member
func AddMembers(data interface{}) error {

	var members Members
	err := mapstructure.Decode(data, &members)
	if err != nil {
		return err
	}

	err = db.Create(&members).Error
	if err != nil {
		return err
	}
	return nil
}

// ExistMembersByEmail : get member data by email
func ExistMembersByEmail(email string) (Members, error) {
	var member Members
	err := db.Model(&Members{}).Where("email = ? AND deleted_on = ?", email, 0).First(&member).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return Members{}, err
	}
	if member.ID > 0 {
		return member, nil
	}
	return Members{}, nil
}

// EditMembers :
func EditMembers(id int64, data interface{}) error {
	if err := db.Model(&Members{}).Where("id = ? and deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
