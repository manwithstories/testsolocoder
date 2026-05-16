package backupcmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"timetrack/internal/database"
	"time"

	"github.com/spf13/cobra"
)

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份数据",
	Long:  "创建SQLite数据库备份文件到指定位置",
	Run: func(cmd *cobra.Command, args []string) {
		outputPath, _ := cmd.Flags().GetString("output")

		db, err := database.GetDB()
		if err != nil {
			fmt.Printf("获取数据库连接失败: %v\n", err)
			return
		}

		srcPath, err := database.GetDBPath()
		if err != nil {
			fmt.Printf("获取数据库路径失败: %v\n", err)
			return
		}

		if outputPath == "" {
			backupDir := filepath.Join(filepath.Dir(srcPath), "backups")
			if err := os.MkdirAll(backupDir, 0755); err != nil {
				fmt.Printf("创建备份目录失败: %v\n", err)
				return
			}
			outputPath = filepath.Join(backupDir, fmt.Sprintf("timetrack_backup_%s.db", time.Now().Format("20060102_150405")))
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			fmt.Printf("打开源数据库失败: %v\n", err)
			return
		}
		defer srcFile.Close()

		destFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("创建备份文件失败: %v\n", err)
			return
		}
		defer destFile.Close()

		if err := db.Close(); err != nil {
			fmt.Printf("关闭数据库连接失败: %v\n", err)
			return
		}

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			fmt.Printf("复制数据库失败: %v\n", err)
			return
		}

		fmt.Printf("备份成功！\n")
		fmt.Printf("备份文件: %s\n", outputPath)

		fileInfo, err := os.Stat(outputPath)
		if err == nil {
			fmt.Printf("文件大小: %.2f MB\n", float64(fileInfo.Size())/1024/1024)
		}
	},
}

var RestoreCmd = &cobra.Command{
	Use:   "restore [backup_file]",
	Short: "恢复数据",
	Long:  "从备份文件恢复数据，冲突时按模式选择跳过或覆盖",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupPath := args[0]
		mode, _ := cmd.Flags().GetString("mode")

		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			fmt.Printf("备份文件不存在: %s\n", backupPath)
			return
		}

		currentDB, err := database.GetDB()
		if err != nil {
			fmt.Printf("获取当前数据库连接失败: %v\n", err)
			return
		}
		defer currentDB.Close()

		backupDB, err := sql.Open("sqlite3", backupPath)
		if err != nil {
			fmt.Printf("打开备份数据库失败: %v\n", err)
			return
		}
		defer backupDB.Close()

		stats := struct {
			projectsInserted int
			projectsUpdated  int
			projectsSkipped  int
			entriesInserted  int
			entriesUpdated   int
			entriesSkipped   int
		}{}

		projectRows, err := backupDB.Query(`SELECT id, name, hourly_rate, currency, archived, created_at, updated_at FROM projects`)
		if err != nil {
			fmt.Printf("读取备份项目数据失败: %v\n", err)
			return
		}
		defer projectRows.Close()

		for projectRows.Next() {
			var id int64
			var name, currency string
			var hourlyRate float64
			var archived bool
			var createdAt, updatedAt time.Time

			err := projectRows.Scan(&id, &name, &hourlyRate, &currency, &archived, &createdAt, &updatedAt)
			if err != nil {
				fmt.Printf("读取项目记录失败: %v\n", err)
				return
			}

			var existingID int64
			err = currentDB.QueryRow(`SELECT id FROM projects WHERE id = ?`, id).Scan(&existingID)

			if err == sql.ErrNoRows {
				_, err = currentDB.Exec(`
					INSERT INTO projects (id, name, hourly_rate, currency, archived, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, ?, ?)
				`, id, name, hourlyRate, currency, archived, formatTime(createdAt), formatTime(updatedAt))
				if err != nil {
					fmt.Printf("插入项目失败: %v\n", err)
					return
				}
				stats.projectsInserted++
			} else if err == nil {
				action := mode
				if action == "" {
					fmt.Printf("\n项目冲突: ID=%d, 名称='%s'\n", id, name)
					fmt.Printf("请选择处理方式:\n")
					fmt.Printf("  1) 覆盖 - 用备份数据覆盖当前数据\n")
					fmt.Printf("  2) 跳过 - 保留当前数据\n")
					fmt.Printf("  3) 全部覆盖 - 后续所有冲突都覆盖\n")
					fmt.Printf("  4) 全部跳过 - 后续所有冲突都跳过\n")
					fmt.Printf("请输入选项 (1/2/3/4): ")

					reader := bufio.NewReader(os.Stdin)
					choice, _ := reader.ReadString('\n')
					choice = strings.TrimSpace(choice)

					switch choice {
					case "1":
						action = "overwrite"
					case "2":
						action = "skip"
					case "3":
						action = "overwrite"
						mode = "overwrite"
					case "4":
						action = "skip"
						mode = "skip"
					default:
						action = "skip"
					}
				}

				if action == "overwrite" {
					_, err = currentDB.Exec(`
						UPDATE projects SET name = ?, hourly_rate = ?, currency = ?, archived = ?, updated_at = ?
						WHERE id = ?
					`, name, hourlyRate, currency, archived, formatTime(updatedAt), id)
					if err != nil {
						fmt.Printf("更新项目失败: %v\n", err)
						return
					}
					stats.projectsUpdated++
				} else {
					stats.projectsSkipped++
				}
			} else {
				fmt.Printf("查询项目失败: %v\n", err)
				return
			}
		}

		entryRows, err := backupDB.Query(`
			SELECT id, project_id, start_time, end_time, paused, paused_at, total_paused_seconds, created_at, updated_at
			FROM time_entries
		`)
		if err != nil {
			fmt.Printf("读取备份时间记录失败: %v\n", err)
			return
		}
		defer entryRows.Close()

		for entryRows.Next() {
			var id, projectID int64
			var startTime, createdAt, updatedAt time.Time
			var endTime, pausedAt sql.NullTime
			var paused bool
			var totalPausedSeconds int64

			err := entryRows.Scan(&id, &projectID, &startTime, &endTime, &paused, &pausedAt, &totalPausedSeconds, &createdAt, &updatedAt)
			if err != nil {
				fmt.Printf("读取时间记录失败: %v\n", err)
				return
			}

			var existingID int64
			err = currentDB.QueryRow(`SELECT id FROM time_entries WHERE id = ?`, id).Scan(&existingID)

			if err == sql.ErrNoRows {
				_, err = currentDB.Exec(`
					INSERT INTO time_entries (id, project_id, start_time, end_time, paused, paused_at, total_paused_seconds, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
				`, id, projectID, formatTime(startTime), nullTimeToPtr(endTime), paused, nullTimeToPtr(pausedAt), totalPausedSeconds, formatTime(createdAt), formatTime(updatedAt))
				if err != nil {
					fmt.Printf("插入时间记录失败: %v\n", err)
					return
				}

				tagRows, err := backupDB.Query(`
					SELECT t.name FROM tags t
					JOIN time_entry_tags tet ON t.id = tet.tag_id
					WHERE tet.time_entry_id = ?
				`, id)
				if err == nil {
					for tagRows.Next() {
						var tagName string
						if tagRows.Scan(&tagName) == nil {
							var tagID int64
							currentDB.QueryRow(`SELECT id FROM tags WHERE name = ?`, tagName).Scan(&tagID)
							if tagID == 0 {
								res, _ := currentDB.Exec(`INSERT INTO tags (name) VALUES (?)`, tagName)
								tagID, _ = res.LastInsertId()
							}
							currentDB.Exec(`INSERT OR IGNORE INTO time_entry_tags (time_entry_id, tag_id) VALUES (?, ?)`, id, tagID)
						}
					}
					tagRows.Close()
				}

				stats.entriesInserted++
			} else if err == nil {
				action := mode
				if action == "" {
					fmt.Printf("\n时间记录冲突: ID=%d, 项目ID=%d, 开始时间=%s\n", id, projectID, formatTime(startTime))
					fmt.Printf("请选择处理方式:\n")
					fmt.Printf("  1) 覆盖 - 用备份数据覆盖当前数据\n")
					fmt.Printf("  2) 跳过 - 保留当前数据\n")
					fmt.Printf("  3) 全部覆盖 - 后续所有冲突都覆盖\n")
					fmt.Printf("  4) 全部跳过 - 后续所有冲突都跳过\n")
					fmt.Printf("请输入选项 (1/2/3/4): ")

					reader := bufio.NewReader(os.Stdin)
					choice, _ := reader.ReadString('\n')
					choice = strings.TrimSpace(choice)

					switch choice {
					case "1":
						action = "overwrite"
					case "2":
						action = "skip"
					case "3":
						action = "overwrite"
						mode = "overwrite"
					case "4":
						action = "skip"
						mode = "skip"
					default:
						action = "skip"
					}
				}

				if action == "overwrite" {
					_, err = currentDB.Exec(`
						UPDATE time_entries SET project_id = ?, start_time = ?, end_time = ?, paused = ?, paused_at = ?,
						total_paused_seconds = ?, updated_at = ?
						WHERE id = ?
					`, projectID, formatTime(startTime), nullTimeToPtr(endTime), paused, nullTimeToPtr(pausedAt), totalPausedSeconds, formatTime(updatedAt), id)
					if err != nil {
						fmt.Printf("更新时间记录失败: %v\n", err)
						return
					}
					stats.entriesUpdated++
				} else {
					stats.entriesSkipped++
				}
			} else {
				fmt.Printf("查询时间记录失败: %v\n", err)
				return
			}
		}

		fmt.Println("\n恢复完成！")
		fmt.Printf("项目: 新增 %d, 更新 %d, 跳过 %d\n", stats.projectsInserted, stats.projectsUpdated, stats.projectsSkipped)
		fmt.Printf("时间记录: 新增 %d, 更新 %d, 跳过 %d\n", stats.entriesInserted, stats.entriesUpdated, stats.entriesSkipped)
	},
}

const timeLayout = "2006-01-02 15:04:05"

func formatTime(t time.Time) string {
	return t.Format(timeLayout)
}

func nullTimeToPtr(nt sql.NullTime) *string {
	if nt.Valid {
		s := nt.Time.Format(timeLayout)
		return &s
	}
	return nil
}

func init() {
	BackupCmd.Flags().StringP("output", "o", "", "备份文件输出路径，默认备份到 ~/.timetrack/backups/")
	RestoreCmd.Flags().String("mode", "", "恢复模式: overwrite (覆盖所有冲突) 或 skip (跳过所有冲突)，不指定则逐条询问")
}
