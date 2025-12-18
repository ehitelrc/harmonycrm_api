package models

type CaseStatsSummary struct {
	TotalCases      int64 `json:"total_cases"`
	OpenCases       int64 `json:"open_cases"`
	ClosedCases     int64 `json:"closed_cases"`
	UnassignedCases int64 `json:"unassigned_cases"`
	UnreadCases     int64 `json:"unread_cases"`
}

type CaseStatsGroup struct {
	Label string `json:"label"`
	Total int64  `json:"total"`
}

type CaseStatsResponse struct {
	Summary       CaseStatsSummary `json:"summary"`
	ByChannel     []CaseStatsGroup `json:"by_channel"`
	ByAgent       []CaseStatsGroup `json:"by_agent"`
	ByFunnelStage []CaseStatsGroup `json:"by_funnel_stage"`
}
