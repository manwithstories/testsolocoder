package handlers

import (
	"smart-energy-platform/models"
	"smart-energy-platform/services"
	"smart-energy-platform/utils"

	"github.com/gin-gonic/gin"
)

type CreateSceneRequest struct {
	FamilyID    uint                  `json:"familyId" binding:"required"`
	Name        string                `json:"name" binding:"required,min=2,max=100"`
	Description string                `json:"description"`
	Icon        string                `json:"icon"`
	IsActive    bool                  `json:"isActive"`
	Conditions  []SceneConditionReq   `json:"conditions"`
	Actions     []SceneActionReq      `json:"actions" binding:"required,min=1"`
}

type SceneConditionReq struct {
	Type     string `json:"type" binding:"required,oneof=time device sensor"`
	DeviceID *uint  `json:"deviceId"`
	Operator string `json:"operator" binding:"omitempty,oneof=eq neq gt lt"`
	Value    string `json:"value"`
	TimeExpr string `json:"timeExpr"`
}

type SceneActionReq struct {
	DeviceID uint   `json:"deviceId" binding:"required"`
	Action   string `json:"action" binding:"required,oneof=on off toggle dim set_temp"`
	Value    string `json:"value"`
}

type UpdateSceneRequest struct {
	Name        string                `json:"name" binding:"min=2,max=100"`
	Description string                `json:"description"`
	Icon        string                `json:"icon"`
	IsActive    *bool                 `json:"isActive"`
	Conditions  []SceneConditionReq   `json:"conditions"`
	Actions     []SceneActionReq      `json:"actions"`
}

func ListScenes(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var familyIDs []uint
	if familyID > 0 {
		familyIDs = []uint{familyID}
	} else {
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
	}

	var scenes []models.Scene
	query := models.DB.Preload("Conditions").Preload("Actions")

	if len(familyIDs) > 0 {
		query = query.Where("family_id IN ?", familyIDs)
	}

	query.Find(&scenes)
	utils.Success(c, scenes)
}

func CreateScene(c *gin.Context) {
	userID := c.GetUint("userId")

	var req CreateSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if !hasFamilyAccess(userID, req.FamilyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	scene := models.Scene{
		FamilyID:    req.FamilyID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
	}

	for _, cond := range req.Conditions {
		scene.Conditions = append(scene.Conditions, models.SceneCondition{
			Type:     cond.Type,
			DeviceID: cond.DeviceID,
			Operator: cond.Operator,
			Value:    cond.Value,
			TimeExpr: cond.TimeExpr,
		})
	}

	for _, act := range req.Actions {
		scene.Actions = append(scene.Actions, models.SceneAction{
			DeviceID: act.DeviceID,
			Action:   act.Action,
			Value:    act.Value,
		})
	}

	if err := models.DB.Create(&scene).Error; err != nil {
		utils.InternalError(c, "Failed to create scene")
		return
	}

	utils.Success(c, scene)
}

func GetScene(c *gin.Context) {
	userID := c.GetUint("userId")
	sceneID := parseUintParam(c, "id")

	var scene models.Scene
	if err := models.DB.Preload("Conditions").Preload("Actions").First(&scene, sceneID).Error; err != nil {
		utils.NotFound(c, "Scene not found")
		return
	}

	if !hasFamilyAccess(userID, scene.FamilyID) {
		utils.Forbidden(c, "No access to this scene")
		return
	}

	utils.Success(c, scene)
}

func UpdateScene(c *gin.Context) {
	userID := c.GetUint("userId")
	sceneID := parseUintParam(c, "id")

	var scene models.Scene
	if err := models.DB.Preload("Conditions").Preload("Actions").First(&scene, sceneID).Error; err != nil {
		utils.NotFound(c, "Scene not found")
		return
	}

	if !hasFamilyAdminAccess(userID, scene.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	var req UpdateSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	models.DB.Model(&scene).Updates(updates)

	if req.Conditions != nil {
		models.DB.Where("scene_id = ?", sceneID).Delete(&models.SceneCondition{})
		for _, cond := range req.Conditions {
			models.DB.Create(&models.SceneCondition{
				SceneID:  sceneID,
				Type:     cond.Type,
				DeviceID: cond.DeviceID,
				Operator: cond.Operator,
				Value:    cond.Value,
				TimeExpr: cond.TimeExpr,
			})
		}
	}

	if req.Actions != nil {
		models.DB.Where("scene_id = ?", sceneID).Delete(&models.SceneAction{})
		for _, act := range req.Actions {
			models.DB.Create(&models.SceneAction{
				SceneID:  sceneID,
				DeviceID: act.DeviceID,
				Action:   act.Action,
				Value:    act.Value,
			})
		}
	}

	models.DB.Preload("Conditions").Preload("Actions").First(&scene, sceneID)
	utils.Success(c, scene)
}

func DeleteScene(c *gin.Context) {
	userID := c.GetUint("userId")
	sceneID := parseUintParam(c, "id")

	var scene models.Scene
	if err := models.DB.First(&scene, sceneID).Error; err != nil {
		utils.NotFound(c, "Scene not found")
		return
	}

	if !hasFamilyAdminAccess(userID, scene.FamilyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	models.DB.Where("scene_id = ?", sceneID).Delete(&models.SceneCondition{})
	models.DB.Where("scene_id = ?", sceneID).Delete(&models.SceneAction{})
	models.DB.Delete(&scene)

	utils.Success(c, nil)
}

func ExecuteScene(c *gin.Context) {
	userID := c.GetUint("userId")
	sceneID := parseUintParam(c, "id")

	var scene models.Scene
	if err := models.DB.Preload("Actions").First(&scene, sceneID).Error; err != nil {
		utils.NotFound(c, "Scene not found")
		return
	}

	if !hasFamilyAccess(userID, scene.FamilyID) {
		utils.Forbidden(c, "No access to this scene")
		return
	}

	if err := services.ExecuteSceneNow(sceneID); err != nil {
		utils.InternalError(c, "Failed to execute scene")
		return
	}

	utils.Success(c, gin.H{
		"scene":    scene.Name,
		"executed": true,
	})
}
