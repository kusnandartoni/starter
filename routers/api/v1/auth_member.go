package v1

import (
	"fmt"
	"kusnandartoni/starter/models"
	"kusnandartoni/starter/services/svcmembers"
	"net/http"

	"github.com/segmentio/ksuid"

	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

func registerMember(form interface{}) (int, string) {
	var (
		err            error
		membersService svcmembers.Members
	)

	err = mapstructure.Decode(form, &membersService)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%v", err)
	}

	member, err := membersService.GetByEmail()
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return http.StatusUnprocessableEntity, fmt.Sprintf("%v", err)
	}

	if member.ID > 0 {
		return http.StatusUnprocessableEntity, "Email already exist"
	}

	uuid := ksuid.New()
	membersService.UUID = uuid.String()
	membersService.CreatedBy = uuid.String()
	membersService.ModifiedBy = uuid.String()
	err = membersService.Add()
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%v", err)
	}
	return 0, ""
}

func loginMember(form interface{}) (int, string, models.Members) {
	var (
		err            error
		membersService svcmembers.Members
	)

	err = mapstructure.Decode(form, &membersService)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%v", err), models.Members{}
	}

	member, err := membersService.Identify()
	if err != nil {
		return http.StatusUnprocessableEntity, fmt.Sprintf("%v", err), models.Members{}
	}

	if !member.Verified {
		return http.StatusUnprocessableEntity, fmt.Sprintf("Please Verify your account berfore login"), models.Members{}
	}
	return 0, "", member
}

func verifyMember(email string) (int, string) {
	var (
		err            error
		membersService svcmembers.Members
	)

	membersService.Email = email
	member, err := membersService.GetByEmail()
	if err != nil {
		return http.StatusUnprocessableEntity, fmt.Sprintf("%v", err)
	}

	membersService.ID = member.ID
	membersService.Verified = true

	err = membersService.Edit()
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%v", err)
	}

	return 0, ""
}

func forgotMember(email string) (int, string, string) {
	var (
		err            error
		membersService svcmembers.Members
	)

	membersService.Email = email

	member, err := membersService.GetByEmail()
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%v", err), ""
	}

	return 0, "", member.FullName
}

func resetMember(email, hashPwd string) (int, string) {
	var (
		err            error
		membersService svcmembers.Members
	)

	membersService.Email = email

	member, err := membersService.GetByEmail()
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%v", err)
	}

	membersService.ID = member.ID
	membersService.Password = hashPwd

	err = membersService.Edit()
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%v", err)
	}

	return 0, ""
}
