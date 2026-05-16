package report

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"timetrack/internal/database"
	"timetrack/internal/models"
	"timetrack/pkg/utils"
	"time"

	"github.com/spf13/cobra"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "统计报表",
	Long:  "按日、周、月查看工作时长和收入统计",
}

type TagStatSorted struct {
	TagName string
	Hours   float64
}

var dailyReportCmd = &cobra.Command{
	Use:   "daily",
	Short: "日报表",
	Run: func(cmd *cobra.Command, args []string) {
		dateStr, _ := cmd.Flags().GetString("date")
		tagFilter, _ := cmd.Flags().GetString("tag")
		var reportDate time.Time

		if dateStr != "" {
			parsed, err := utils.ParseDate(dateStr)
			if err != nil {
				fmt.Printf("日期格式错误: %v\n", err)
				return
			}
			reportDate = parsed
		} else {
			reportDate = time.Now()
		}

		generateDailyReport(reportDate, tagFilter)
	},
}

var weeklyReportCmd = &cobra.Command{
	Use:   "weekly",
	Short: "周报表",
	Run: func(cmd *cobra.Command, args []string) {
		dateStr, _ := cmd.Flags().GetString("date")
		tagFilter, _ := cmd.Flags().GetString("tag")
		var reportDate time.Time

		if dateStr != "" {
			parsed, err := utils.ParseDate(dateStr)
			if err != nil {
				fmt.Printf("日期格式错误: %v\n", err)
				return
			}
			reportDate = parsed
		} else {
			reportDate = time.Now()
		}

		generateWeeklyReport(reportDate, tagFilter)
	},
}

var monthlyReportCmd = &cobra.Command{
	Use:   "monthly",
	Short: "月报表",
	Run: func(cmd *cobra.Command, args []string) {
		dateStr, _ := cmd.Flags().GetString("date")
		tagFilter, _ := cmd.Flags().GetString("tag")
		var reportDate time.Time

		if dateStr != "" {
			parsed, err := utils.ParseDate(dateStr)
			if err != nil {
				fmt.Printf("日期格式错误: %v\n", err)
				return
			}
			reportDate = parsed
		} else {
			reportDate = time.Now()
		}

		generateMonthlyReport(reportDate, tagFilter)
	},
}

func generateDailyReport(date time.Time, tagFilter string) {
	startDate := utils.StartOfDay(date)
	endDate := utils.EndOfDay(date)

	var tagPtr *string
	if tagFilter != "" {
		tagPtr = &tagFilter
	}

	entries, err := database.ListTimeEntries(startDate, endDate, nil, tagPtr)
	if err != nil {
		fmt.Printf("获取时间记录失败: %v\n", err)
		return
	}

	title := fmt.Sprintf("日报表 - %s", utils.FormatDate(date))
	if tagFilter != "" {
		title += fmt.Sprintf(" (标签: %s)", tagFilter)
	}
	fmt.Println(title)
	fmt.Println(strings.Repeat("=", 50))

	if len(entries) == 0 {
		fmt.Println("该日期没有时间记录")
		return
	}

	projectStats := make(map[string]*models.ProjectStat)
	tagStats := make(map[string]float64)
	var totalHours float64
	var totalIncome float64
	currency := "CNY"

	for _, entry := range entries {
		if entry.EndTime == nil {
			continue
		}

		project, _ := database.GetProjectByID(entry.ProjectID)
		if project == nil {
			continue
		}
		currency = project.Currency

		duration := database.GetElapsedDuration(&entry.TimeEntry)
		hours := duration.Hours()
		income := hours * project.HourlyRate

		if _, ok := projectStats[project.Name]; !ok {
			projectStats[project.Name] = &models.ProjectStat{ProjectName: project.Name}
		}
		projectStats[project.Name].Hours += hours
		projectStats[project.Name].Income += income

		for _, tag := range entry.Tags {
			tagStats[tag] += hours
		}

		totalHours += hours
		totalIncome += income
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "项目\t时长(小时)\t收入")
	fmt.Fprintln(w, "----\t----------\t----")

	for _, ps := range projectStats {
		fmt.Fprintf(w, "%s\t%.2f\t%.2f\n", ps.ProjectName, ps.Hours, ps.Income)
	}
	fmt.Fprintf(w, "\n合计\t%.2f\t%.2f %s\n", totalHours, totalIncome, currency)
	w.Flush()

	if len(tagStats) > 0 {
		fmt.Println()
		fmt.Println("标签统计:")
		sortedTags := sortTagsByHours(tagStats)
		for _, ts := range sortedTags {
			fmt.Printf("  %s: %.2f小时\n", ts.TagName, ts.Hours)
		}
	}
}

func generateWeeklyReport(date time.Time, tagFilter string) {
	startDate := utils.StartOfWeek(date)
	endDate := utils.EndOfWeek(date)

	var tagPtr *string
	if tagFilter != "" {
		tagPtr = &tagFilter
	}

	entries, err := database.ListTimeEntries(startDate, endDate, nil, tagPtr)
	if err != nil {
		fmt.Printf("获取时间记录失败: %v\n", err)
		return
	}

	title := fmt.Sprintf("周报表 - %s 至 %s", utils.FormatDate(startDate), utils.FormatDate(endDate))
	if tagFilter != "" {
		title += fmt.Sprintf(" (标签: %s)", tagFilter)
	}
	fmt.Println(title)
	fmt.Println(strings.Repeat("=", 60))

	if len(entries) == 0 {
		fmt.Println("该周没有时间记录")
		return
	}

	projectStats := make(map[string]*models.ProjectStat)
	tagStats := make(map[string]float64)
	dailyStats := make(map[string]float64)
	var totalHours float64
	var totalIncome float64
	currency := "CNY"

	for _, entry := range entries {
		if entry.EndTime == nil {
			continue
		}

		project, _ := database.GetProjectByID(entry.ProjectID)
		if project == nil {
			continue
		}
		currency = project.Currency

		duration := database.GetElapsedDuration(&entry.TimeEntry)
		hours := duration.Hours()
		income := hours * project.HourlyRate
		dayKey := utils.FormatDate(entry.StartTime)

		if _, ok := projectStats[project.Name]; !ok {
			projectStats[project.Name] = &models.ProjectStat{ProjectName: project.Name}
		}
		projectStats[project.Name].Hours += hours
		projectStats[project.Name].Income += income

		for _, tag := range entry.Tags {
			tagStats[tag] += hours
		}

		dailyStats[dayKey] += hours
		totalHours += hours
		totalIncome += income
	}

	days := 7
	avgDaily := totalHours / float64(days)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "项目\t时长(小时)\t收入")
	fmt.Fprintln(w, "----\t----------\t----")

	for _, ps := range projectStats {
		fmt.Fprintf(w, "%s\t%.2f\t%.2f\n", ps.ProjectName, ps.Hours, ps.Income)
	}
	fmt.Fprintf(w, "\n合计\t%.2f\t%.2f %s\n", totalHours, totalIncome, currency)
	fmt.Fprintf(w, "日均\t%.2f\t-\n", avgDaily)
	w.Flush()

	fmt.Println()
	fmt.Println("每日时长:")
	for day := startDate; !day.After(endDate); day = day.AddDate(0, 0, 1) {
		dayKey := utils.FormatDate(day)
		hours := dailyStats[dayKey]
		bar := strings.Repeat("█", int(hours*2))
		fmt.Printf("  %s: %5.2fh |%s\n", dayKey, hours, bar)
	}

	if len(tagStats) > 0 {
		fmt.Println()
		fmt.Println("最活跃标签:")
		sortedTags := sortTagsByHours(tagStats)
		for i, ts := range sortedTags {
			if i >= 5 {
				break
			}
			fmt.Printf("  %d. %s: %.2f小时\n", i+1, ts.TagName, ts.Hours)
		}
	}
}

func generateMonthlyReport(date time.Time, tagFilter string) {
	startDate := utils.StartOfMonth(date)
	endDate := utils.EndOfMonth(date)

	var tagPtr *string
	if tagFilter != "" {
		tagPtr = &tagFilter
	}

	entries, err := database.ListTimeEntries(startDate, endDate, nil, tagPtr)
	if err != nil {
		fmt.Printf("获取时间记录失败: %v\n", err)
		return
	}

	title := fmt.Sprintf("月报表 - %s", startDate.Format("2006年1月"))
	if tagFilter != "" {
		title += fmt.Sprintf(" (标签: %s)", tagFilter)
	}
	fmt.Println(title)
	fmt.Println(strings.Repeat("=", 60))

	if len(entries) == 0 {
		fmt.Println("该月没有时间记录")
		return
	}

	projectStats := make(map[string]*models.ProjectStat)
	tagStats := make(map[string]float64)
	weeklyStats := make(map[int]float64)
	var totalHours float64
	var totalIncome float64
	currency := "CNY"
	activeDays := make(map[string]bool)

	for _, entry := range entries {
		if entry.EndTime == nil {
			continue
		}

		project, _ := database.GetProjectByID(entry.ProjectID)
		if project == nil {
			continue
		}
		currency = project.Currency

		duration := database.GetElapsedDuration(&entry.TimeEntry)
		hours := duration.Hours()
		income := hours * project.HourlyRate
		dayKey := utils.FormatDate(entry.StartTime)
		weekNum := getWeekOfMonth(entry.StartTime)

		if _, ok := projectStats[project.Name]; !ok {
			projectStats[project.Name] = &models.ProjectStat{ProjectName: project.Name}
		}
		projectStats[project.Name].Hours += hours
		projectStats[project.Name].Income += income

		for _, tag := range entry.Tags {
			tagStats[tag] += hours
		}

		weeklyStats[weekNum] += hours
		activeDays[dayKey] = true
		totalHours += hours
		totalIncome += income
	}

	numDays := endDate.Day()
	avgDaily := totalHours / float64(numDays)
	avgWorkDay := totalHours / float64(len(activeDays))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "项目\t时长(小时)\t收入")
	fmt.Fprintln(w, "----\t----------\t----")

	for _, ps := range projectStats {
		fmt.Fprintf(w, "%s\t%.2f\t%.2f\n", ps.ProjectName, ps.Hours, ps.Income)
	}
	fmt.Fprintf(w, "\n合计\t%.2f\t%.2f %s\n", totalHours, totalIncome, currency)
	fmt.Fprintf(w, "日均(全部)\t%.2f\t-\n", avgDaily)
	fmt.Fprintf(w, "日均(工作日)\t%.2f\t-\n", avgWorkDay)
	fmt.Fprintf(w, "工作天数\t%d\t-\n", len(activeDays))
	w.Flush()

	fmt.Println()
	fmt.Println("每周时长:")
	for week := 1; week <= 5; week++ {
		hours := weeklyStats[week]
		bar := strings.Repeat("█", int(hours))
		fmt.Printf("  第%d周: %6.2fh |%s\n", week, hours, bar)
	}

	if len(tagStats) > 0 {
		fmt.Println()
		fmt.Println("最活跃标签:")
		sortedTags := sortTagsByHours(tagStats)
		for i, ts := range sortedTags {
			if i >= 5 {
				break
			}
			fmt.Printf("  %d. %s: %.2f小时\n", i+1, ts.TagName, ts.Hours)
		}
	}
}

func getWeekOfMonth(t time.Time) int {
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	_, firstWeek := firstDay.ISOWeek()
	_, currentWeek := t.ISOWeek()
	return currentWeek - firstWeek + 1
}

func sortTagsByHours(tagStats map[string]float64) []TagStatSorted {
	var result []TagStatSorted
	for tag, hours := range tagStats {
		result = append(result, TagStatSorted{TagName: tag, Hours: hours})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Hours > result[j].Hours
	})
	return result
}

func init() {
	dailyReportCmd.Flags().String("date", "", "指定日期 (YYYY-MM-DD)")
	dailyReportCmd.Flags().String("tag", "", "按标签过滤")

	weeklyReportCmd.Flags().String("date", "", "指定周内任意日期 (YYYY-MM-DD)")
	weeklyReportCmd.Flags().String("tag", "", "按标签过滤")

	monthlyReportCmd.Flags().String("date", "", "指定月内任意日期 (YYYY-MM-DD)")
	monthlyReportCmd.Flags().String("tag", "", "按标签过滤")

	ReportCmd.AddCommand(dailyReportCmd)
	ReportCmd.AddCommand(weeklyReportCmd)
	ReportCmd.AddCommand(monthlyReportCmd)
}
