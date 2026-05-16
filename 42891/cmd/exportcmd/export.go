package exportcmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"timetrack/internal/database"
	"timetrack/internal/models"
	"timetrack/pkg/utils"
	"time"

	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "数据导出",
	Long:  "导出时间记录为CSV或JSON格式，支持按日期范围和项目过滤",
}

var exportCSVCmd = &cobra.Command{
	Use:   "csv",
	Short: "导出为CSV格式",
	Run: func(cmd *cobra.Command, args []string) {
		exportData("csv", cmd)
	},
}

var exportJSONCmd = &cobra.Command{
	Use:   "json",
	Short: "导出为JSON格式",
	Run: func(cmd *cobra.Command, args []string) {
		exportData("json", cmd)
	},
}

func exportData(format string, cmd *cobra.Command) {
	startStr, _ := cmd.Flags().GetString("start")
	endStr, _ := cmd.Flags().GetString("end")
	projectFilter, _ := cmd.Flags().GetString("project")
	tagFilter, _ := cmd.Flags().GetString("tag")
	outputFile, _ := cmd.Flags().GetString("output")

	var startDate, endDate time.Time
	var err error

	if startStr != "" {
		startDate, err = utils.ParseDate(startStr)
		if err != nil {
			fmt.Printf("开始日期格式错误: %v\n", err)
			return
		}
		startDate = utils.StartOfDay(startDate)
	} else {
		startDate = utils.StartOfDay(time.Now().AddDate(0, 0, -30))
	}

	if endStr != "" {
		endDate, err = utils.ParseDate(endStr)
		if err != nil {
			fmt.Printf("结束日期格式错误: %v\n", err)
			return
		}
		endDate = utils.EndOfDay(endDate)
	} else {
		endDate = utils.EndOfDay(time.Now())
	}

	var projectID *int64
	if projectFilter != "" {
		project, err := database.GetProjectByName(projectFilter)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if project == nil {
			fmt.Printf("项目 '%s' 不存在\n", projectFilter)
			return
		}
		projectID = &project.ID
	}

	var tagPtr *string
	if tagFilter != "" {
		tagPtr = &tagFilter
	}

	entries, err := database.ListTimeEntries(startDate, endDate, projectID, tagPtr)
	if err != nil {
		fmt.Printf("获取时间记录失败: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("没有找到符合条件的时间记录")
		return
	}

	var output string
	if outputFile == "" {
		outputFile = fmt.Sprintf("timetrack_export_%s.%s", time.Now().Format("20060102_150405"), format)
	}

	if format == "csv" {
		output, err = exportToCSV(entries)
	} else {
		output, err = exportToJSON(entries)
	}

	if err != nil {
		fmt.Printf("导出失败: %v\n", err)
		return
	}

	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("导出成功！共导出 %d 条记录\n", len(entries))
	fmt.Printf("文件已保存到: %s\n", outputFile)
}

func exportToCSV(entries []models.TimeEntryWithTags) (string, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	headers := []string{"ID", "项目", "开始时间", "结束时间", "时长(小时)", "标签", "状态"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	for _, e := range entries {
		status := "已完成"
		var endTimeStr string
		var hours float64

		if e.EndTime == nil {
			status = "进行中"
			if e.Paused {
				status = "已暂停"
			}
			endTimeStr = ""
			hours = database.GetElapsedDuration(&e.TimeEntry).Hours()
		} else {
			endTimeStr = utils.FormatDateTime(*e.EndTime)
			hours = database.GetElapsedDuration(&e.TimeEntry).Hours()
		}

		row := []string{
			fmt.Sprintf("%d", e.ID),
			e.ProjectName,
			utils.FormatDateTime(e.StartTime),
			endTimeStr,
			fmt.Sprintf("%.2f", hours),
			strings.Join(e.Tags, "|"),
			status,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type JSONExportEntry struct {
	ID        int64    `json:"id"`
	Project   string   `json:"project"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time,omitempty"`
	Hours     float64  `json:"hours"`
	Tags      []string `json:"tags"`
	Status    string   `json:"status"`
}

type JSONExport struct {
	ExportTime string            `json:"export_time"`
	Count      int               `json:"count"`
	Entries    []JSONExportEntry `json:"entries"`
}

func exportToJSON(entries []models.TimeEntryWithTags) (string, error) {
	jsonEntries := make([]JSONExportEntry, 0, len(entries))

	for _, e := range entries {
		status := "已完成"
		var hours float64
		var endTimeStr string

		if e.EndTime == nil {
			status = "进行中"
			if e.Paused {
				status = "已暂停"
			}
			hours = database.GetElapsedDuration(&e.TimeEntry).Hours()
		} else {
			endTimeStr = e.EndTime.Format(time.RFC3339)
			hours = database.GetElapsedDuration(&e.TimeEntry).Hours()
		}

		jsonEntries = append(jsonEntries, JSONExportEntry{
			ID:        e.ID,
			Project:   e.ProjectName,
			StartTime: e.StartTime.Format(time.RFC3339),
			EndTime:   endTimeStr,
			Hours:     hours,
			Tags:      e.Tags,
			Status:    status,
		})
	}

	export := JSONExport{
		ExportTime: time.Now().Format(time.RFC3339),
		Count:      len(entries),
		Entries:    jsonEntries,
	}

	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func init() {
	for _, cmd := range []*cobra.Command{exportCSVCmd, exportJSONCmd} {
		cmd.Flags().String("start", "", "开始日期 (YYYY-MM-DD), 默认30天前")
		cmd.Flags().String("end", "", "结束日期 (YYYY-MM-DD), 默认今天")
		cmd.Flags().String("project", "", "按项目名称过滤")
		cmd.Flags().String("tag", "", "按标签过滤")
		cmd.Flags().StringP("output", "o", "", "输出文件路径")
	}

	ExportCmd.AddCommand(exportCSVCmd)
	ExportCmd.AddCommand(exportJSONCmd)
}
