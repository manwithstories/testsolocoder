package migrations

import (
	"medical-platform/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Doctor{},
		&models.Schedule{},
		&models.Patient{},
		&models.HealthRecord{},
		&models.Appointment{},
		&models.Consultation{},
		&models.Prescription{},
		&models.PrescriptionItem{},
		&models.ExaminationReport{},
		&models.Notification{},
		&models.Payment{},
		&models.Review{},
	)
	if err != nil {
		return err
	}

	return seedData(db)
}

func seedData(db *gorm.DB) error {
	var deptCount int64
	db.Model(&models.Department{}).Count(&deptCount)
	if deptCount > 0 {
		return nil
	}

	departments := []models.Department{
		{Name: "内科", Description: "内科诊疗", Location: "门诊楼2层"},
		{Name: "外科", Description: "外科诊疗", Location: "门诊楼2层"},
		{Name: "儿科", Description: "儿科诊疗", Location: "门诊楼1层"},
		{Name: "妇产科", Description: "妇产科诊疗", Location: "门诊楼3层"},
		{Name: "眼科", Description: "眼科诊疗", Location: "门诊楼4层"},
		{Name: "耳鼻喉科", Description: "耳鼻喉科诊疗", Location: "门诊楼4层"},
		{Name: "口腔科", Description: "口腔科诊疗", Location: "门诊楼1层"},
		{Name: "皮肤科", Description: "皮肤科诊疗", Location: "门诊楼3层"},
	}

	for _, dept := range departments {
		if err := db.Create(&dept).Error; err != nil {
			return err
		}
	}

	return nil
}
