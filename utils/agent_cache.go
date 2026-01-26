package utils

import (
	"fmt"
	"time"
)

/* ---------- Keys ---------- */

func AgentsAllKey() string {
	return "agents:all"
}

func AgentByUserIDKey(userID uint) string {
	return fmt.Sprintf("agents:user:%d", userID)
}

func AgentsWithUserInfoAllKey() string {
	return "agents:with-user-info:all"
}

func AgentWithUserInfoByUserIDKey(userID uint) string {
	return fmt.Sprintf("agents:with-user-info:user:%d", userID)
}

func AgentsByCompanyWithUserInfoKey(companyID uint) string {
	return fmt.Sprintf("agents:company:%d:with-user-info", companyID)
}

func AgentsByCompanyAndDepartmentWithUserInfoKey(companyID, departmentID uint) string {
	return fmt.Sprintf("agents:company:%d:department:%d:with-user-info", companyID, departmentID)
}

func NonAgentsAllKey() string {
	return "agents:non-agents:all"
}

/* ---------- Cache Set ---------- */

func CacheAgentsAll(v interface{}) {
	SetToCache(AgentsAllKey(), v, 5*time.Minute)
}

func CacheAgentByUserID(userID uint, v interface{}) {
	SetToCache(AgentByUserIDKey(userID), v, 5*time.Minute)
}

func CacheAgentsWithUserInfoAll(v interface{}) {
	SetToCache(AgentsWithUserInfoAllKey(), v, 5*time.Minute)
}

func CacheAgentWithUserInfoByUserID(userID uint, v interface{}) {
	SetToCache(AgentWithUserInfoByUserIDKey(userID), v, 5*time.Minute)
}

func CacheAgentsByCompanyWithUserInfo(companyID uint, v interface{}) {
	SetToCache(AgentsByCompanyWithUserInfoKey(companyID), v, 3*time.Minute)
}

func CacheAgentsByCompanyAndDepartmentWithUserInfo(companyID, departmentID uint, v interface{}) {
	SetToCache(
		AgentsByCompanyAndDepartmentWithUserInfoKey(companyID, departmentID),
		v,
		3*time.Minute,
	)
}

func CacheNonAgentsAll(v interface{}) {
	SetToCache(NonAgentsAllKey(), v, 5*time.Minute)
}

/* ---------- Cache Get ---------- */

func GetAgentsAllFromCache() (interface{}, bool) {
	return GetFromCache(AgentsAllKey())
}

func GetAgentByUserIDFromCache(userID uint) (interface{}, bool) {
	return GetFromCache(AgentByUserIDKey(userID))
}

func GetAgentsWithUserInfoAllFromCache() (interface{}, bool) {
	return GetFromCache(AgentsWithUserInfoAllKey())
}

func GetAgentWithUserInfoByUserIDFromCache(userID uint) (interface{}, bool) {
	return GetFromCache(AgentWithUserInfoByUserIDKey(userID))
}

func GetAgentsByCompanyWithUserInfoFromCache(companyID uint) (interface{}, bool) {
	return GetFromCache(AgentsByCompanyWithUserInfoKey(companyID))
}

func GetAgentsByCompanyAndDepartmentWithUserInfoFromCache(companyID, departmentID uint) (interface{}, bool) {
	return GetFromCache(AgentsByCompanyAndDepartmentWithUserInfoKey(companyID, departmentID))
}

func GetNonAgentsAllFromCache() (interface{}, bool) {
	return GetFromCache(NonAgentsAllKey())
}

/* ---------- Invalidate ---------- */

func InvalidateAgentsCache() {
	DeleteCache(AgentsAllKey())
	DeleteCache(AgentsWithUserInfoAllKey())
	DeleteCache(NonAgentsAllKey())
}
