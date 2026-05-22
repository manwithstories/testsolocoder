package repository

import (
	"fmt"
	"qa-platform/config"
	"qa-platform/models"
	"qa-platform/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	utils.LogInfo("数据库连接成功")

	err = autoMigrate()
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	utils.LogInfo("数据库迁移完成")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.Answer{},
		&models.Comment{},
		&models.Category{},
		&models.Tag{},
		&models.QuestionTag{},
		&models.AuditRecord{},
		&models.Report{},
		&models.SensitiveWord{},
		&models.PointLog{},
		&models.Reward{},
		&models.RewardExchange{},
		&models.Favorite{},
		&models.Follow{},
		&models.Notification{},
		&models.ExpertApplication{},
	)
}

func SeedData() error {
	categories := []models.Category{
		{Name: "编程开发", Description: "编程语言、框架、工具等相关问题", Icon: "💻", SortOrder: 1},
		{Name: "人工智能", Description: "机器学习、深度学习、NLP等AI相关问题", Icon: "🤖", SortOrder: 2},
		{Name: "数据科学", Description: "数据分析、数据可视化、大数据等", Icon: "📊", SortOrder: 3},
		{Name: "云计算", Description: "云服务、容器、微服务等", Icon: "☁️", SortOrder: 4},
		{Name: "前端开发", Description: "HTML/CSS/JS、框架、UI等", Icon: "🎨", SortOrder: 5},
		{Name: "后端开发", Description: "服务端开发、API、数据库等", Icon: "⚙️", SortOrder: 6},
		{Name: "移动开发", Description: "iOS、Android、跨平台等", Icon: "📱", SortOrder: 7},
		{Name: "网络安全", Description: "安全、加密、渗透测试等", Icon: "🔒", SortOrder: 8},
	}
	for _, cat := range categories {
		DB.FirstOrCreate(&cat, models.Category{Name: cat.Name})
	}

	tags := []models.Tag{
		{Name: "Go", Description: "Go语言相关"},
		{Name: "Python", Description: "Python语言相关"},
		{Name: "JavaScript", Description: "JavaScript语言相关"},
		{Name: "Vue", Description: "Vue框架相关"},
		{Name: "React", Description: "React框架相关"},
		{Name: "Docker", Description: "Docker容器相关"},
		{Name: "Kubernetes", Description: "K8s相关"},
		{Name: "MySQL", Description: "MySQL数据库相关"},
		{Name: "Redis", Description: "Redis缓存相关"},
		{Name: "算法", Description: "算法与数据结构"},
	}
	for _, tag := range tags {
		DB.FirstOrCreate(&tag, models.Tag{Name: tag.Name})
	}

	sensitiveWords := []models.SensitiveWord{
		{Word: "敏感词1", Category: "政治", Level: 1, ReplaceTo: "***"},
		{Word: "敏感词2", Category: "政治", Level: 1, ReplaceTo: "***"},
		{Word: "敏感词3", Category: "色情", Level: 2, ReplaceTo: "***"},
		{Word: "敏感词4", Category: "暴力", Level: 2, ReplaceTo: "***"},
	}
	for _, sw := range sensitiveWords {
		DB.FirstOrCreate(&sw, models.SensitiveWord{Word: sw.Word})
	}

	rewards := []models.Reward{
		{Name: "徽章-初学者", Description: "完成第一个回答获得", Image: "🏅", PointsCost: 100, Stock: -1},
		{Name: "徽章-专家", Description: "获得专家认证", Image: "🎖️", PointsCost: 1000, Stock: -1},
		{Name: "会员月卡", Description: "30天会员权限", Image: "👑", PointsCost: 500, Stock: 100},
		{Name: "主题皮肤", Description: "专属主题皮肤", Image: "🎨", PointsCost: 200, Stock: -1},
	}
	for _, reward := range rewards {
		DB.FirstOrCreate(&reward, models.Reward{Name: reward.Name})
	}

	adminUser := models.User{
		Username: "admin",
		Email:    "admin@qa.com",
		Nickname: "管理员",
		Role:     "admin",
		Status:   "active",
		Level:    10,
		Points:   10000,
		IsExpert: true,
	}
	adminUser.Password, _ = utils.HashPassword("admin123")
	DB.FirstOrCreate(&adminUser, models.User{Username: adminUser.Username})

	return nil
}
