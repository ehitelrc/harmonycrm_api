package repository

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/models"
)

type LoginRepository struct{}

func NewLoginRepository() *LoginRepository {
	return &LoginRepository{}
}

func (lr *LoginRepository) Login(email string, password string, companyID int) (*models.VWUsersCompanies, error) {

	var validate struct {
		HasPermission bool `gorm:"column:has_permission"`
		UserID        *int `gorm:"column:user_id"`
	}

	if err := config.DB.
		Raw(`SELECT * FROM public.fn_validate_user_company(?, ?, ?)`,
			companyID, email, password).
		Scan(&validate).Error; err != nil {
		return nil, err
	}

	if !validate.HasPermission {
		return nil, fmt.Errorf("usuario no autorizado para esta compañía")
	}

	var user models.VWUsersCompanies
	if err := config.DB.
		Where("user_id = ? AND company_id = ?", validate.UserID, companyID).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	return &user, nil

	// var userRolePermissions []models.UserRolePermission
	// if err := config.DB.
	// 	Raw(`SELECT * FROM vw_user_role_permissions WHERE user_id = ? AND company_id = ?`,
	// 		validate.UserID, companyID).
	// 	Scan(&userRolePermissions).Error; err != nil {
	// 	return nil, err
	// }

	// return &userRolePermissions, nil
}

func (lr *LoginRepository) GetUserPermissions(userID int, companyID int) (*[]models.UserEffectivePermission, error) {
	var permissions []models.UserEffectivePermission

	err := config.DB.Where("user_id = ? AND company_id = ?", userID, companyID).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	var data struct {
		RowsAffected int64 `gorm:"column:rows_affected"`
	}
	data.RowsAffected = int64(len(permissions))
	if data.RowsAffected == 0 {
		return nil, nil
	}

	return &permissions, nil

}
