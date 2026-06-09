package utils

import (
	"fmt"
	"time"
)

/* ---------- Keys ---------- */

func DepartmentsAllKey() string {
	return "departments:all"
}

func DepartmentByIDKey(id uint) string {
	return fmt.Sprintf("departments:id:%d", id)
}

func DepartmentsByCompanyKey(companyID uint) string {
	return fmt.Sprintf("departments:company:%d", companyID)
}

func DepartmentsByCompanyAndUserKey(companyID, userID uint) string {
	return fmt.Sprintf("departments:company:%d:user:%d", companyID, userID)
}

/* ---------- Cache Set ---------- */

func CacheDepartmentsAll(depts interface{}) {
	SetToCache(DepartmentsAllKey(), depts, 5*time.Minute)
}

func CacheDepartmentByID(id uint, dept interface{}) {
	SetToCache(DepartmentByIDKey(id), dept, 5*time.Minute)
}

func CacheDepartmentsByCompany(companyID uint, depts interface{}) {
	SetToCache(DepartmentsByCompanyKey(companyID), depts, 5*time.Minute)
}

func CacheDepartmentsByCompanyAndUser(companyID, userID uint, depts interface{}) {
	SetToCache(
		DepartmentsByCompanyAndUserKey(companyID, userID),
		depts,
		3*time.Minute,
	)
}

/* ---------- Cache Get ---------- */

func GetDepartmentsAllFromCache() (interface{}, bool) {
	return GetFromCache(DepartmentsAllKey())
}

func GetDepartmentByIDFromCache(id uint) (interface{}, bool) {
	return GetFromCache(DepartmentByIDKey(id))
}

func GetDepartmentsByCompanyFromCache(companyID uint) (interface{}, bool) {
	return GetFromCache(DepartmentsByCompanyKey(companyID))
}

func GetDepartmentsByCompanyAndUserFromCache(companyID, userID uint) (interface{}, bool) {
	return GetFromCache(DepartmentsByCompanyAndUserKey(companyID, userID))
}

/* ---------- Invalidate ---------- */

func InvalidateDepartmentsCache() {
	DeleteCache(DepartmentsAllKey())
}

func InvalidateDepartmentByID(id uint) {
	DeleteCache(DepartmentByIDKey(id))
}

func InvalidateDepartmentsByCompany(companyID uint) {
	DeleteCache(DepartmentsByCompanyKey(companyID))
}
