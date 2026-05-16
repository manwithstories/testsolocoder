package timecmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"timetrack/internal/database"
	"timetrack/internal/models"
	"timetrack/pkg/utils"
	"time"

	"github.com/spf13/cobra"
)

var TimeCmd = &cobra.Command{
	Use:   "time",
	Short: "时间记录管理",
	Long:  "开始、停止、暂停、恢复计时，查看当前计时状态",
}

var startCmd = &cobra.Command{
	Use:   "start [project]",
	Short: "开始计时",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		tags, _ := cmd.Flags().GetStringSlice("tag")

		project, err := database.GetProjectByName(projectName)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if project == nil {
			fmt.Printf("项目 '%s' 不存在\n", projectName)
			return
		}
		if project.Archived {
			fmt.Printf("项目 '%s' 已归档，无法记时\n", projectName)
			return
		}

		activeEntry, err := database.GetActiveTimeEntry()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		if activeEntry != nil {
			conflict, truncateAt, err := database.CheckForCrossDayConflict(activeEntry)
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}

			if conflict {
				fmt.Printf("检测到跨天计时冲突！\n")
				activeProject, _ := database.GetProjectByID(activeEntry.ProjectID)
				fmt.Printf("当前有未停止的计时: 项目 '%s', 开始于 %s\n",
					activeProject.Name, utils.FormatDateTime(activeEntry.StartTime))
				fmt.Printf("该计时已跨天，建议在 %s 截断\n", utils.FormatDateTime(truncateAt))
				fmt.Printf("请选择处理方式:\n")
				fmt.Printf("  1) 截断 - 以 %s 作为结束时间\n", utils.FormatDateTime(truncateAt))
				fmt.Printf("  2) 丢弃 - 删除该计时记录\n")
				fmt.Printf("  3) 取消 - 不做处理\n")
				fmt.Printf("请输入选项 (1/2/3): ")

				reader := bufio.NewReader(os.Stdin)
				choice, _ := reader.ReadString('\n')
				choice = strings.TrimSpace(choice)

				switch choice {
				case "1":
					if err := database.TruncateTimeEntry(activeEntry.ID, truncateAt); err != nil {
						fmt.Printf("截断失败: %v\n", err)
						return
					}
					fmt.Println("已截断跨天计时")
				case "2":
					if err := deleteTimeEntry(activeEntry.ID); err != nil {
						fmt.Printf("删除失败: %v\n", err)
						return
					}
					fmt.Println("已丢弃跨天计时")
				default:
					fmt.Println("已取消操作")
					return
				}
			} else {
				fmt.Printf("已有正在进行的计时，请先停止\n")
				return
			}
		}

		entry, err := database.StartTimeEntry(project.ID, time.Now(), tags)
		if err != nil {
			fmt.Printf("开始计时失败: %v\n", err)
			return
		}

		fmt.Printf("开始计时: 项目 '%s'\n", projectName)
		fmt.Printf("  开始时间: %s\n", utils.FormatDateTime(entry.StartTime))
		if len(tags) > 0 {
			fmt.Printf("  标签: %s\n", utils.JoinTags(tags))
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止当前计时",
	Run: func(cmd *cobra.Command, args []string) {
		activeEntry, err := database.GetActiveTimeEntry()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if activeEntry == nil {
			fmt.Println("当前没有进行中的计时")
			return
		}

		if activeEntry.Paused && activeEntry.PausedAt != nil {
			if err := database.ResumeTimeEntry(activeEntry.ID); err != nil {
				fmt.Printf("恢复计时失败: %v\n", err)
				return
			}
		}
		endTime := time.Now()

		conflict, truncateAt, err := database.CheckForCrossDayConflict(activeEntry)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		if conflict {
			fmt.Printf("警告: 该计时已跨天（开始于 %s）\n", utils.FormatDateTime(activeEntry.StartTime))
			fmt.Printf("请选择处理方式:\n")
			fmt.Printf("  1) 截断 - 以 %s 作为结束时间\n", utils.FormatDateTime(truncateAt))
			fmt.Printf("  2) 正常停止 - 以当前时间结束\n")
			fmt.Printf("请输入选项 (1/2): ")

			reader := bufio.NewReader(os.Stdin)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			if choice == "1" {
				endTime = truncateAt
			}
		}

		if err := database.StopTimeEntry(activeEntry.ID, endTime); err != nil {
			fmt.Printf("停止计时失败: %v\n", err)
			return
		}

		stoppedEntry, _ := database.GetTimeEntryByID(activeEntry.ID)
		project, _ := database.GetProjectByID(activeEntry.ProjectID)
		duration := database.GetElapsedDuration(stoppedEntry)
		hours := duration.Hours()
		income := hours * project.HourlyRate

		fmt.Printf("计时已停止\n")
		fmt.Printf("  项目: %s\n", project.Name)
		fmt.Printf("  时长: %s (%.2f小时)\n", utils.FormatDuration(duration), hours)
		fmt.Printf("  收入: %.2f %s\n", income, project.Currency)
		if stoppedEntry.TotalPausedSeconds > 0 {
			pausedDuration := time.Duration(stoppedEntry.TotalPausedSeconds) * time.Second
			fmt.Printf("  累计暂停: %s\n", utils.FormatDuration(pausedDuration))
		}
		if len(activeEntry.Tags) > 0 {
			fmt.Printf("  标签: %s\n", utils.JoinTags(activeEntry.Tags))
		}
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "暂停当前计时",
	Run: func(cmd *cobra.Command, args []string) {
		activeEntry, err := database.GetActiveTimeEntry()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if activeEntry == nil {
			fmt.Println("当前没有进行中的计时")
			return
		}
		if activeEntry.Paused {
			fmt.Println("计时已经暂停")
			return
		}

		if err := database.PauseTimeEntry(activeEntry.ID); err != nil {
			fmt.Printf("暂停失败: %v\n", err)
			return
		}

		project, _ := database.GetProjectByID(activeEntry.ProjectID)
		duration := time.Since(activeEntry.StartTime)

		fmt.Printf("计时已暂停\n")
		fmt.Printf("  项目: %s\n", project.Name)
		fmt.Printf("  已计时: %s\n", utils.FormatDuration(duration))
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "恢复已暂停的计时",
	Run: func(cmd *cobra.Command, args []string) {
		activeEntry, err := database.GetActiveTimeEntry()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if activeEntry == nil {
			fmt.Println("当前没有进行中的计时")
			return
		}
		if !activeEntry.Paused {
			fmt.Println("计时未暂停")
			return
		}

		if err := database.ResumeTimeEntry(activeEntry.ID); err != nil {
			fmt.Printf("恢复失败: %v\n", err)
			return
		}

		project, _ := database.GetProjectByID(activeEntry.ProjectID)
		fmt.Printf("计时已恢复\n")
		fmt.Printf("  项目: %s\n", project.Name)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看当前计时状态",
	Run: func(cmd *cobra.Command, args []string) {
		activeEntry, err := database.GetActiveTimeEntry()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if activeEntry == nil {
			fmt.Println("当前没有进行中的计时")
			return
		}

		project, err := database.GetProjectByID(activeEntry.ProjectID)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		elapsed := database.GetElapsedDuration(activeEntry)
		hours := elapsed.Hours()
		income := hours * project.HourlyRate

		status := "进行中"
		if activeEntry.Paused {
			status = "已暂停"
		}

		fmt.Printf("当前计时状态: %s\n", status)
		fmt.Printf("  项目: %s\n", project.Name)
		fmt.Printf("  开始时间: %s\n", utils.FormatDateTime(activeEntry.StartTime))
		fmt.Printf("  已计时: %s (%.2f小时)\n", utils.FormatDuration(elapsed), hours)
		fmt.Printf("  当前收入: %.2f %s\n", income, project.Currency)
		if activeEntry.TotalPausedSeconds > 0 {
			pausedDuration := time.Duration(activeEntry.TotalPausedSeconds) * time.Second
			if activeEntry.Paused && activeEntry.PausedAt != nil {
				pausedDuration += time.Since(*activeEntry.PausedAt)
			}
			fmt.Printf("  累计暂停: %s\n", utils.FormatDuration(pausedDuration))
		}
		if len(activeEntry.Tags) > 0 {
			fmt.Printf("  标签: %s\n", utils.JoinTags(activeEntry.Tags))
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出最近的时间记录",
	Run: func(cmd *cobra.Command, args []string) {
		days, _ := cmd.Flags().GetInt("days")
		projectFilter, _ := cmd.Flags().GetString("project")
		tagFilter, _ := cmd.Flags().GetString("tag")

		endDate := utils.EndOfDay(time.Now())
		startDate := utils.StartOfDay(time.Now().AddDate(0, 0, -(days - 1)))

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
			fmt.Println("暂无时间记录")
			return
		}

		filterDesc := fmt.Sprintf("最近 %d 天", days)
		if tagFilter != "" {
			filterDesc += fmt.Sprintf(", 标签: %s", tagFilter)
		}
		fmt.Printf("时间记录 (%s):\n", filterDesc)
		fmt.Println(strings.Repeat("-", 80))

		for _, e := range entries {
			status := "已完成"
			duration := database.GetElapsedDuration(&e.TimeEntry)

			if e.EndTime == nil {
				status = "进行中"
				if e.Paused {
					status = "已暂停"
				}
			}

			fmt.Printf("%s | %s | %s\n",
				utils.FormatDateTime(e.StartTime),
				e.ProjectName,
				status)
			fmt.Printf("  时长: %s (%.2f小时) | 标签: %s\n",
				utils.FormatDuration(duration),
				duration.Hours(),
				utils.JoinTags(e.Tags))
			fmt.Println()
		}
	},
}

func deleteTimeEntry(entryID int64) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM time_entry_tags WHERE time_entry_id = ?`, entryID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM time_entries WHERE id = ?`, entryID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func init() {
	startCmd.Flags().StringSliceP("tag", "t", []string{}, "为计时添加标签，可多次使用")

	listCmd.Flags().Int("days", 7, "显示最近多少天的记录")
	listCmd.Flags().String("project", "", "按项目过滤")
	listCmd.Flags().String("tag", "", "按标签过滤")

	TimeCmd.AddCommand(startCmd)
	TimeCmd.AddCommand(stopCmd)
	TimeCmd.AddCommand(pauseCmd)
	TimeCmd.AddCommand(resumeCmd)
	TimeCmd.AddCommand(statusCmd)
	TimeCmd.AddCommand(listCmd)
}

func formatEntry(e *models.TimeEntry, projectName string) string {
	return fmt.Sprintf("项目: %s, 开始: %s, 标签: %s",
		projectName, utils.FormatDateTime(e.StartTime), utils.JoinTags(e.Tags))
}
