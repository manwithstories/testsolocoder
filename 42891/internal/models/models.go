package models

import "time"

type Project struct {
	ID        int64
	Name      string
	HourlyRate float64
	Currency  string
	Archived  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TimeEntry struct {
	ID                 int64
	ProjectID          int64
	StartTime          time.Time
	EndTime            *time.Time
	Paused             bool
	PausedAt           *time.Time
	TotalPausedSeconds  int64
	Tags               []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Tag struct {
	ID   int64
	Name string
}

type TimeEntryWithTags struct {
	TimeEntry
	ProjectName string
	Tags        []string
}

type DailyReport struct {
	Date         time.Time
	TotalHours   float64
	TotalIncome  float64
	ProjectStats []ProjectStat
	TopTags      []TagStat
}

type ProjectStat struct {
	ProjectName string
	Hours       float64
	Income      float64
}

type TagStat struct {
	TagName string
	Hours   float64
}
