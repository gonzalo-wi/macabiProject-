package expensesdto

import (
	expensesports "macabi-back/internal/expenses/application/ports"
)

type SummaryResponse struct {
	TotalApproved string                `json:"total_approved"`
	ByMonth       []MonthlyBucketResponse `json:"by_month"`
}

type MonthlyBucketResponse struct {
	Month string `json:"month"`
	Total string `json:"total"`
}

type AnalyticsResponse struct {
	TotalApproved string                    `json:"total_approved"`
	TotalCount    int64                     `json:"total_count"`
	PendingCount  int64                     `json:"pending_count"`
	ApprovedCount int64                     `json:"approved_count"`
	RejectedCount int64                     `json:"rejected_count"`
	Granularity   string                    `json:"granularity"`
	ByProject     []AnalyticsProjectResponse `json:"by_project"`
	ByBucket      []AnalyticsBucketResponse  `json:"by_bucket"`
}

type AnalyticsProjectResponse struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Total       string `json:"total"`
}

type AnalyticsBucketResponse struct {
	Bucket string `json:"bucket"`
	Total  string `json:"total"`
}

func AnalyticsToResp(a *expensesports.ExpenseAnalyticsResult) AnalyticsResponse {
	byProject := make([]AnalyticsProjectResponse, len(a.ByProject))
	for i, p := range a.ByProject {
		byProject[i] = AnalyticsProjectResponse{ProjectID: p.ProjectID, ProjectName: p.ProjectName, Total: p.Total.String()}
	}
	byBucket := make([]AnalyticsBucketResponse, len(a.ByBucket))
	for i, b := range a.ByBucket {
		byBucket[i] = AnalyticsBucketResponse{Bucket: b.Bucket, Total: b.Total.String()}
	}
	return AnalyticsResponse{
		TotalApproved: a.TotalApproved.String(),
		TotalCount:    a.TotalCount,
		PendingCount:  a.PendingCount,
		ApprovedCount: a.ApprovedCount,
		RejectedCount: a.RejectedCount,
		Granularity:   a.Granularity,
		ByProject:     byProject,
		ByBucket:      byBucket,
	}
}

func SummaryToResp(s *expensesports.ProjectExpenseSummary) SummaryResponse {
	by := make([]MonthlyBucketResponse, len(s.ByMonth))
	for i, b := range s.ByMonth {
		by[i] = MonthlyBucketResponse{Month: b.MonthYYYYMM, Total: b.Total.String()}
	}
	return SummaryResponse{TotalApproved: s.TotalApproved.String(), ByMonth: by}
}
