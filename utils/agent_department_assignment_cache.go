package utils

import (
	"fmt"
	"time"
)

/* ---------- Keys ---------- */

func ADAByIDKey(id uint) string {
	return fmt.Sprintf("agent-dept-assign:id:%d", id)
}

func ADAByAgentKey(agentID uint) string {
	return fmt.Sprintf("agent-dept-assign:agent:%d", agentID)
}

func ADAByDepartmentKey(deptID uint) string {
	return fmt.Sprintf("agent-dept-assign:department:%d", deptID)
}

func ADAByCompanyKey(companyID uint) string {
	return fmt.Sprintf("agent-dept-assign:company:%d", companyID)
}

func ADAByCompanyAndAgentKey(companyID, agentID uint) string {
	return fmt.Sprintf("agent-dept-assign:company:%d:agent:%d", companyID, agentID)
}

func ADAAgentsByDepartmentKey(companyID, deptID uint) string {
	return fmt.Sprintf("agent-dept-assign:company:%d:department:%d:agents", companyID, deptID)
}

/* ---------- Cache SET ---------- */

func CacheADAByID(id uint, v interface{}) {
	SetToCache(ADAByIDKey(id), v, 5*time.Minute)
}

func CacheADAByAgent(agentID uint, v interface{}) {
	SetToCache(ADAByAgentKey(agentID), v, 3*time.Minute)
}

func CacheADAByDepartment(deptID uint, v interface{}) {
	SetToCache(ADAByDepartmentKey(deptID), v, 3*time.Minute)
}

func CacheADAByCompany(companyID uint, v interface{}) {
	SetToCache(ADAByCompanyKey(companyID), v, 3*time.Minute)
}

func CacheADAByCompanyAndAgent(companyID, agentID uint, v interface{}) {
	SetToCache(
		ADAByCompanyAndAgentKey(companyID, agentID),
		v,
		3*time.Minute,
	)
}

func CacheAgentsByDepartment(companyID, deptID uint, v interface{}) {
	SetToCache(
		ADAAgentsByDepartmentKey(companyID, deptID),
		v,
		3*time.Minute,
	)
}

/* ---------- Cache GET ---------- */

func GetADAByIDFromCache(id uint) (interface{}, bool) {
	return GetFromCache(ADAByIDKey(id))
}

func GetADAByAgentFromCache(agentID uint) (interface{}, bool) {
	return GetFromCache(ADAByAgentKey(agentID))
}

func GetADAByDepartmentFromCache(deptID uint) (interface{}, bool) {
	return GetFromCache(ADAByDepartmentKey(deptID))
}

func GetADAByCompanyFromCache(companyID uint) (interface{}, bool) {
	return GetFromCache(ADAByCompanyKey(companyID))
}

func GetADAByCompanyAndAgentFromCache(companyID, agentID uint) (interface{}, bool) {
	return GetFromCache(ADAByCompanyAndAgentKey(companyID, agentID))
}

func GetAgentsByDepartmentFromCache(companyID, deptID uint) (interface{}, bool) {
	return GetFromCache(ADAAgentsByDepartmentKey(companyID, deptID))
}

/* ---------- Invalidate ---------- */

// invalidación global mínima (segura)
func InvalidateADAGlobal() {
	// No se hace flush total: TTL se encarga del resto
}
