package utils

import (
	"fmt"
	"time"
)

func CompaniesAllKey() string {
	return "companies:all"
}

func CompanyByIDKey(id uint) string {
	return fmt.Sprintf("companies:id:%d", id)
}

func CacheCompaniesAll(companies interface{}) {
	SetToCache(CompaniesAllKey(), companies, 5*time.Minute)
}

func CacheCompanyByID(id uint, company interface{}) {
	SetToCache(CompanyByIDKey(id), company, 5*time.Minute)
}

func GetCompaniesAllFromCache() (interface{}, bool) {
	return GetFromCache(CompaniesAllKey())
}

func GetCompanyByIDFromCache(id uint) (interface{}, bool) {
	return GetFromCache(CompanyByIDKey(id))
}

func InvalidateCompaniesCache() {
	DeleteCache(CompaniesAllKey())
}

func InvalidateCompanyByID(id uint) {
	DeleteCache(CompanyByIDKey(id))
}

// ---------- By User ----------

func CompaniesByUserIDKey(userID uint) string {
	return fmt.Sprintf("companies:user:%d", userID)
}

func CacheCompaniesByUserID(userID uint, companies interface{}) {
	SetToCache(CompaniesByUserIDKey(userID), companies, 3*time.Minute)
}

func GetCompaniesByUserIDFromCache(userID uint) (interface{}, bool) {
	return GetFromCache(CompaniesByUserIDKey(userID))
}

func InvalidateCompaniesByUserID(userID uint) {
	DeleteCache(CompaniesByUserIDKey(userID))
}
