package project

import (
	"fmt"
	"os"
	"text/tabwriter"
	"timetrack/internal/database"
	"timetrack/internal/models"
	"timetrack/pkg/utils"

	"github.com/spf13/cobra"
)

var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "项目管理",
	Long:  "创建、列出、归档、删除项目，以及设置项目时薪",
}

var createProjectCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "创建新项目",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		hourlyRate, _ := cmd.Flags().GetFloat64("rate")
		currency, _ := cmd.Flags().GetString("currency")

		existing, err := database.GetProjectByName(name)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if existing != nil {
			fmt.Printf("项目 '%s' 已存在\n", name)
			return
		}

		project, err := database.CreateProject(name, hourlyRate, currency)
		if err != nil {
			fmt.Printf("创建项目失败: %v\n", err)
			return
		}

		fmt.Printf("项目创建成功!\n")
		fmt.Printf("  ID: %d\n", project.ID)
		fmt.Printf("  名称: %s\n", project.Name)
		fmt.Printf("  时薪: %.2f %s\n", project.HourlyRate, project.Currency)
	},
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有项目",
	Run: func(cmd *cobra.Command, args []string) {
		includeArchived, _ := cmd.Flags().GetBool("all")

		projects, err := database.ListProjects(includeArchived)
		if err != nil {
			fmt.Printf("获取项目列表失败: %v\n", err)
			return
		}

		if len(projects) == 0 {
			fmt.Println("暂无项目")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\t名称\t时薪\t币种\t状态\t创建时间")
		fmt.Fprintln(w, "--\t----\t----\t----\t----\t--------")

		for _, p := range projects {
			status := "活跃"
			if p.Archived {
				status = "已归档"
			}
			fmt.Fprintf(w, "%d\t%s\t%.2f\t%s\t%s\t%s\n",
				p.ID, p.Name, p.HourlyRate, p.Currency, status,
				utils.FormatDate(p.CreatedAt))
		}

		w.Flush()
	},
}

var archiveProjectCmd = &cobra.Command{
	Use:   "archive [name]",
	Short: "归档项目",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		project, err := database.GetProjectByName(name)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if project == nil {
			fmt.Printf("项目 '%s' 不存在\n", name)
			return
		}

		if project.Archived {
			fmt.Printf("项目 '%s' 已归档\n", name)
			return
		}

		activeEntry, err := database.GetActiveTimeEntry()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if activeEntry != nil && activeEntry.ProjectID == project.ID {
			fmt.Printf("项目 '%s' 有正在进行的计时，请先停止后再归档\n", name)
			return
		}

		if err := database.ArchiveProject(project.ID); err != nil {
			fmt.Printf("归档项目失败: %v\n", err)
			return
		}

		fmt.Printf("项目 '%s' 已归档，归档后无法再记时但数据保留\n", name)
	},
}

var unarchiveProjectCmd = &cobra.Command{
	Use:   "unarchive [name]",
	Short: "取消归档项目",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		project, err := database.GetProjectByName(name)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if project == nil {
			fmt.Printf("项目 '%s' 不存在\n", name)
			return
		}

		if !project.Archived {
			fmt.Printf("项目 '%s' 未归档\n", name)
			return
		}

		if err := database.UnarchiveProject(project.ID); err != nil {
			fmt.Printf("取消归档失败: %v\n", err)
			return
		}

		fmt.Printf("项目 '%s' 已取消归档\n", name)
	},
}

var deleteProjectCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "删除项目（同时删除所有相关时间记录）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		project, err := database.GetProjectByName(name)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if project == nil {
			fmt.Printf("项目 '%s' 不存在\n", name)
			return
		}

		fmt.Printf("警告: 删除项目将同时删除所有相关时间记录！\n")
		fmt.Printf("确定要删除项目 '%s' 吗? (输入 yes 确认): ", name)

		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "yes" {
			fmt.Println("已取消删除")
			return
		}

		if err := database.DeleteProject(project.ID); err != nil {
			fmt.Printf("删除项目失败: %v\n", err)
			return
		}

		fmt.Printf("项目 '%s' 已删除\n", name)
	},
}

var setRateCmd = &cobra.Command{
	Use:   "set-rate [name]",
	Short: "设置项目时薪",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		hourlyRate, _ := cmd.Flags().GetFloat64("rate")
		currency, _ := cmd.Flags().GetString("currency")

		project, err := database.GetProjectByName(name)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		if project == nil {
			fmt.Printf("项目 '%s' 不存在\n", name)
			return
		}

		if currency == "" {
			currency = project.Currency
		}

		if err := database.SetProjectRate(project.ID, hourlyRate, currency); err != nil {
			fmt.Printf("设置时薪失败: %v\n", err)
			return
		}

		fmt.Printf("项目 '%s' 时薪已更新为: %.2f %s\n", name, hourlyRate, currency)
	},
}

func init() {
	createProjectCmd.Flags().Float64("rate", 0, "项目时薪")
	createProjectCmd.Flags().String("currency", "CNY", "币种 (CNY, USD, EUR 等)")

	listProjectsCmd.Flags().Bool("all", false, "包含已归档项目")

	setRateCmd.Flags().Float64("rate", 0, "新的时薪")
	setRateCmd.Flags().String("currency", "", "新的币种")
	setRateCmd.MarkFlagRequired("rate")

	ProjectCmd.AddCommand(createProjectCmd)
	ProjectCmd.AddCommand(listProjectsCmd)
	ProjectCmd.AddCommand(archiveProjectCmd)
	ProjectCmd.AddCommand(unarchiveProjectCmd)
	ProjectCmd.AddCommand(deleteProjectCmd)
	ProjectCmd.AddCommand(setRateCmd)
}

func formatProject(p *models.Project) string {
	return fmt.Sprintf("ID: %d, 名称: %s, 时薪: %.2f %s, 状态: %s",
		p.ID, p.Name, p.HourlyRate, p.Currency, map[bool]string{true: "已归档", false: "活跃"}[p.Archived])
}
