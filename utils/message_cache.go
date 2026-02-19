package utils

import (
	"fmt"
	"time"
)

/* -------------------------------
   OPEN CASES (LISTADOS)
--------------------------------*/

func OpenCasesKey(companyID, deptID uint, page, limit int) string {
	return fmt.Sprintf(
		"cases:open:company:%d:dept:%d:page:%d:limit:%d",
		companyID, deptID, page, limit,
	)
}

func CacheOpenCases(companyID, deptID uint, page, limit int, data interface{}) {
	SetToCache(OpenCasesKey(companyID, deptID, page, limit), data, 45*time.Second)
}

func GetOpenCasesFromCache(companyID, deptID uint, page, limit int) (interface{}, bool) {
	return GetFromCache(OpenCasesKey(companyID, deptID, page, limit))
}

/* -------------------------------
   OPEN CASES STATS
--------------------------------*/

func OpenCasesStatsKey(companyID, deptID uint) string {
	return fmt.Sprintf(
		"cases:open:stats:company:%d:dept:%d",
		companyID, deptID,
	)
}

func CacheOpenCasesStats(companyID, deptID uint, data interface{}) {
	SetToCache(OpenCasesStatsKey(companyID, deptID), data, 20*time.Second)
}

func GetOpenCasesStatsFromCache(companyID, deptID uint) (interface{}, bool) {
	return GetFromCache(OpenCasesStatsKey(companyID, deptID))
}

/* -------------------------------
   ACTIVE CASES BY AGENT
--------------------------------*/

func ActiveCasesByAgentKey(agentID uint) string {
	return fmt.Sprintf("cases:active:agent:%d", agentID)
}

func CacheActiveCasesByAgent(agentID uint, data interface{}) {
	SetToCache(ActiveCasesByAgentKey(agentID), data, 30*time.Second)
}

func GetActiveCasesByAgentFromCache(agentID uint) (interface{}, bool) {
	return GetFromCache(ActiveCasesByAgentKey(agentID))
}

func InvalidateActiveCasesByAgent(agentID uint) {
	DeleteCache(ActiveCasesByAgentKey(agentID))
}

/* -------------------------------
   FIRST PAGE MESSAGES (UX)
--------------------------------*/

func FirstMessagesByCaseKey(caseID uint) string {
	return fmt.Sprintf("messages:case:%d:first", caseID)
}

func CacheFirstMessagesByCase(caseID uint, data interface{}) {
	SetToCache(FirstMessagesByCaseKey(caseID), data, 15*time.Second)
}

func GetFirstMessagesByCaseFromCache(caseID uint) (interface{}, bool) {
	return GetFromCache(FirstMessagesByCaseKey(caseID))
}

func InvalidateFirstMessagesByCase(caseID uint) {
	DeleteCache(FirstMessagesByCaseKey(caseID))
}
