package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"museum-server/internal/dto"
	"museum-server/pkg/response"
	validatorpkg "museum-server/pkg/validator"
)

func ValidateJSON(target interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetType := reflect.TypeOf(target)
		if targetType.Kind() == reflect.Ptr {
			targetType = targetType.Elem()
		}

		req := reflect.New(targetType).Interface()

		if err := c.ShouldBindJSON(req); err != nil {
			response.Error(c, http.StatusBadRequest, 400, "请求参数格式错误: "+err.Error())
			c.Abort()
			return
		}

		result := validateByType(req)
		if !result.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": result.Message,
				"errors":  result.Errors,
			})
			c.Abort()
			return
		}

		c.Set("validated_data", req)
		c.Next()
	}
}

func ValidateQuery(target interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetType := reflect.TypeOf(target)
		if targetType.Kind() == reflect.Ptr {
			targetType = targetType.Elem()
		}

		req := reflect.New(targetType).Interface()

		if err := c.ShouldBindQuery(req); err != nil {
			response.Error(c, http.StatusBadRequest, 400, "请求参数格式错误: "+err.Error())
			c.Abort()
			return
		}

		result := validateByType(req)
		if !result.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": result.Message,
				"errors":  result.Errors,
			})
			c.Abort()
			return
		}

		c.Set("validated_data", req)
		c.Next()
	}
}

func ValidateUri(target interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetType := reflect.TypeOf(target)
		if targetType.Kind() == reflect.Ptr {
			targetType = targetType.Elem()
		}

		req := reflect.New(targetType).Interface()

		if err := c.ShouldBindUri(req); err != nil {
			response.Error(c, http.StatusBadRequest, 400, "请求参数格式错误: "+err.Error())
			c.Abort()
			return
		}

		result := validateByType(req)
		if !result.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": result.Message,
				"errors":  result.Errors,
			})
			c.Abort()
			return
		}

		c.Set("validated_data", req)
		c.Next()
	}
}

func validateByType(req interface{}) validatorpkg.ValidationResult {
	switch r := req.(type) {
	case *dto.LoginRequest:
		return validatorpkg.ValidateLoginRequest(r)
	case *dto.RegisterRequest:
		return validatorpkg.ValidateRegisterRequest(r)
	case *dto.UpdateUserRequest:
		return validatorpkg.ValidateUpdateUserRequest(r)
	case *dto.CollectionCategoryRequest:
		return validatorpkg.ValidateCollectionCategoryRequest(r)
	case *dto.CollectionRequest:
		return validatorpkg.ValidateCollectionRequest(r)
	case *dto.CollectionListQuery:
		return validatorpkg.ValidateCollectionListQuery(r)
	case *dto.ExhibitionRequest:
		return validatorpkg.ValidateExhibitionRequest(r)
	case *dto.ExhibitionListQuery:
		return validatorpkg.ValidateExhibitionListQuery(r)
	case *dto.ReservationRequest:
		return validatorpkg.ValidateReservationRequest(r)
	case *dto.ReservationCancelRequest:
		return validatorpkg.ValidateStruct(r)
	case *dto.TimeSlotRequest:
		return validatorpkg.ValidateTimeSlotRequest(r)
	case *dto.BatchTimeSlotRequest:
		return validatorpkg.ValidateBatchTimeSlotRequest(r)
	case *dto.GuideScheduleRequest:
		return validatorpkg.ValidateGuideScheduleRequest(r)
	case *dto.GuideContentRequest:
		return validatorpkg.ValidateGuideContentRequest(r)
	case *dto.ResearchApplicationRequest:
		return validatorpkg.ValidateResearchApplicationRequest(r)
	case *dto.ApplicationReviewRequest:
		return validatorpkg.ValidateApplicationReviewRequest(r)
	case *dto.StatisticsQuery:
		return validatorpkg.ValidateStatisticsQuery(r)
	case *dto.MuseumRequest:
		return validatorpkg.ValidateMuseumRequest(r)
	default:
		return validatorpkg.ValidateStruct(req)
	}
}

func GetValidatedData(c *gin.Context) interface{} {
	data, exists := c.Get("validated_data")
	if !exists {
		return nil
	}
	return data
}
