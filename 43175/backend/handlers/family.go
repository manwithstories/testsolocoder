package handlers

import (
	"net/http"
	"smart-energy-platform/models"
	"smart-energy-platform/services"
	"smart-energy-platform/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateFamilyRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateFamilyRequest struct {
	Name        string `json:"name" binding:"min=2,max=100"`
	Description string `json:"description"`
}

type InviteMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin member guest"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member guest"`
}

func ListFamilies(c *gin.Context) {
	userID := c.GetUint("userId")

	var members []models.FamilyMember
	models.DB.Where("user_id = ? AND status = ?", userID, 1).Find(&members)

	var familyIDs []uint
	for _, m := range members {
		familyIDs = append(familyIDs, m.FamilyID)
	}

	var families []models.Family
	if len(familyIDs) > 0 {
		models.DB.Where("id IN ?", familyIDs).
			Preload("Owner").
			Preload("Members.User").
			Find(&families)
	}

	utils.Success(c, families)
}

func CreateFamily(c *gin.Context) {
	userID := c.GetUint("userId")

	var req CreateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	family := models.Family{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := models.DB.Create(&family).Error; err != nil {
		utils.InternalError(c, "Failed to create family")
		return
	}

	member := models.FamilyMember{
		FamilyID: family.ID,
		UserID:   userID,
		Role:     "admin",
		Status:   1,
	}
	models.DB.Create(&member)

	utils.Success(c, family)
}

func GetFamily(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseUintParam(c, "id")

	if !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var family models.Family
	if err := models.DB.Preload("Owner").
		Preload("Members.User").
		First(&family, familyID).Error; err != nil {
		utils.NotFound(c, "Family not found")
		return
	}

	utils.Success(c, family)
}

func UpdateFamily(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseUintParam(c, "id")

	if !hasFamilyAdminAccess(userID, familyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	var req UpdateFamilyRequest
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

	if err := models.DB.Model(&models.Family{}).Where("id = ?", familyID).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update family")
		return
	}

	var family models.Family
	models.DB.Preload("Members.User").First(&family, familyID)

	utils.Success(c, family)
}

func DeleteFamily(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseUintParam(c, "id")

	var family models.Family
	if err := models.DB.First(&family, familyID).Error; err != nil {
		utils.NotFound(c, "Family not found")
		return
	}

	if family.OwnerID != userID {
		utils.Forbidden(c, "Only owner can delete the family")
		return
	}

	models.DB.Where("family_id = ?", familyID).Delete(&models.FamilyMember{})
	models.DB.Where("family_id = ?", familyID).Delete(&models.Device{})
	models.DB.Delete(&family)

	utils.Success(c, nil)
}

func InviteMember(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseUintParam(c, "id")

	if !hasFamilyAdminAccess(userID, familyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	var req InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var existing models.Invitation
	result := models.DB.Where("family_id = ? AND email = ? AND status = ?", familyID, req.Email, "pending").
		First(&existing)

	if result.Error == nil {
		utils.Error(c, http.StatusConflict, 409, "Invitation already sent")
		return
	}

	var user models.User
	models.DB.Where("email = ?", req.Email).First(&user)

	invitation := models.Invitation{
		FamilyID:  familyID,
		InviterID: userID,
		Email:     req.Email,
		Role:      req.Role,
		Status:    "pending",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := models.DB.Create(&invitation).Error; err != nil {
		utils.InternalError(c, "Failed to create invitation")
		return
	}

	if user.ID > 0 {
		services.SendNotification(user.ID, "invitation",
			"Family Invitation",
			"You have been invited to join a family group.")
	}

	utils.Success(c, invitation)
}

func RemoveMember(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseUintParam(c, "id")
	memberID := parseUintParam(c, "memberId")

	if !hasFamilyAdminAccess(userID, familyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	if memberID == userID {
		utils.Error(c, http.StatusBadRequest, 400, "Cannot remove yourself")
		return
	}

	if err := models.DB.Where("id = ? AND family_id = ?", memberID, familyID).
		Delete(&models.FamilyMember{}).Error; err != nil {
		utils.InternalError(c, "Failed to remove member")
		return
	}

	utils.Success(c, nil)
}

func UpdateMemberRole(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseUintParam(c, "id")
	memberID := parseUintParam(c, "memberId")

	if !hasFamilyAdminAccess(userID, familyID) {
		utils.Forbidden(c, "Admin access required")
		return
	}

	var req UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := models.DB.Model(&models.FamilyMember{}).
		Where("id = ? AND family_id = ?", memberID, familyID).
		Update("role", req.Role).Error; err != nil {
		utils.InternalError(c, "Failed to update member role")
		return
	}

	utils.Success(c, nil)
}

func ListInvitations(c *gin.Context) {
	userID := c.GetUint("userId")

	var user models.User
	models.DB.First(&user, userID)

	var invitations []models.Invitation
	models.DB.Where("email = ? AND status = ?", user.Email, "pending").
		Preload("Family").
		Preload("Inviter").
		Find(&invitations)

	utils.Success(c, invitations)
}

func AcceptInvitation(c *gin.Context) {
	userID := c.GetUint("userId")
	invitationID := parseUintParam(c, "id")

	var invitation models.Invitation
	if err := models.DB.First(&invitation, invitationID).Error; err != nil {
		utils.NotFound(c, "Invitation not found")
		return
	}

	if invitation.Status != "pending" {
		utils.Error(c, http.StatusBadRequest, 400, "Invitation already processed")
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		utils.Error(c, http.StatusBadRequest, 400, "Invitation has expired")
		return
	}

	invitation.Status = "accepted"
	models.DB.Save(&invitation)

	member := models.FamilyMember{
		FamilyID: invitation.FamilyID,
		UserID:   userID,
		Role:     invitation.Role,
		Status:   1,
	}
	models.DB.Create(&member)

	utils.Success(c, nil)
}

func RejectInvitation(c *gin.Context) {
	invitationID := parseUintParam(c, "id")

	var invitation models.Invitation
	if err := models.DB.First(&invitation, invitationID).Error; err != nil {
		utils.NotFound(c, "Invitation not found")
		return
	}

	invitation.Status = "rejected"
	models.DB.Save(&invitation)

	utils.Success(c, nil)
}

func parseUintParam(c *gin.Context, name string) uint {
	val, _ := strconv.ParseUint(c.Param(name), 10, 64)
	return uint(val)
}

func hasFamilyAccess(userID, familyID uint) bool {
	var member models.FamilyMember
	err := models.DB.Where("family_id = ? AND user_id = ? AND status = ?", familyID, userID, 1).
		First(&member).Error
	return err == nil
}

func hasFamilyAdminAccess(userID, familyID uint) bool {
	var member models.FamilyMember
	err := models.DB.Where("family_id = ? AND user_id = ? AND role = ? AND status = ?",
		familyID, userID, "admin", 1).First(&member).Error
	return err == nil
}
