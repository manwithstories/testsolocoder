package service

import (
	"errors"
	"fmt"
	"survey-platform/internal/dto"
	"survey-platform/internal/model"
	"survey-platform/internal/repository"
	"time"

	"github.com/google/uuid"
)

type DistributionService struct {
	distRepo   *repository.DistributionRepository
	surveyRepo *repository.SurveyRepository
	emailSvc   *EmailService
}

func NewDistributionService(
	distRepo *repository.DistributionRepository,
	surveyRepo *repository.SurveyRepository,
	emailSvc *EmailService,
) *DistributionService {
	return &DistributionService{
		distRepo:   distRepo,
		surveyRepo: surveyRepo,
		emailSvc:   emailSvc,
	}
}

func (s *DistributionService) CreateLink(surveyID, userID uint, req *dto.CreateDistributionLinkRequest) (*dto.DistributionLinkResponse, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, errors.New("survey not found")
	}
	if survey.UserID != userID {
		return nil, errors.New("permission denied")
	}

	token := uuid.New().String()
	link := &model.DistributionLink{
		SurveyID:  surveyID,
		LinkToken: token,
		Channel:   req.Channel,
		MaxUses:   req.MaxUses,
		ExpiresAt: req.ExpiresAt,
		IsActive:  true,
	}

	if err := s.distRepo.CreateLink(link); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("http://localhost:3000/survey/fill/%s", token)
	qrCodeURL := ""
	if req.Channel == "qrcode" {
		qrCodeURL, _ = s.generateQRCode(fullURL)
	}

	return &dto.DistributionLinkResponse{
		ID:        link.ID,
		LinkToken: link.LinkToken,
		Channel:   link.Channel,
		MaxUses:   link.MaxUses,
		UseCount:  link.UseCount,
		ExpiresAt: link.ExpiresAt,
		IsActive:  link.IsActive,
		FullURL:   fullURL,
		QRCodeURL: qrCodeURL,
		CreatedAt: link.CreatedAt,
	}, nil
}

func (s *DistributionService) GetByToken(token string) (*model.DistributionLink, error) {
	link, err := s.distRepo.FindByToken(token)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, errors.New("link not found")
	}

	if !link.IsActive {
		return nil, errors.New("link is deactivated")
	}

	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return nil, errors.New("link has expired")
	}

	if link.MaxUses > 0 && link.UseCount >= link.MaxUses {
		return nil, errors.New("link has reached maximum uses")
	}

	return link, nil
}

func (s *DistributionService) ListBySurveyID(surveyID, userID uint) ([]*model.DistributionLink, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, errors.New("survey not found")
	}
	if survey.UserID != userID {
		return nil, errors.New("permission denied")
	}

	return s.distRepo.ListBySurveyID(surveyID)
}

func (s *DistributionService) SendInvitations(surveyID, userID uint, req *dto.SendInvitationsRequest) error {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}

	var invitations []*model.Invitation
	for _, email := range req.Emails {
		token := uuid.New().String()
		invitation := &model.Invitation{
			SurveyID:  surveyID,
			Email:     email,
			LinkToken: token,
			Status:    1,
		}
		invitations = append(invitations, invitation)

		surveyLink := &model.DistributionLink{
			SurveyID:  surveyID,
			LinkToken: token,
			Channel:   "email",
			ExpiresAt: req.ExpiresAt,
			IsActive:  true,
		}
		s.distRepo.CreateLink(surveyLink)
	}

	if err := s.distRepo.BatchCreateInvitations(invitations); err != nil {
		return err
	}

	go s.sendInvitationEmails(invitations, survey, req.CustomMessage)

	return nil
}

func (s *DistributionService) sendInvitationEmails(invitations []*model.Invitation, survey *model.Survey, customMessage string) {
	for _, invitation := range invitations {
		subject := fmt.Sprintf("邀请您参与问卷调查: %s", survey.Title)
		surveyLink := fmt.Sprintf("http://localhost:3000/survey/fill/%s", invitation.LinkToken)
		body := fmt.Sprintf(`
			<h2>问卷调查邀请</h2>
			<p>尊敬的用户，您好！</p>
			<p>我们诚邀您参与以下问卷调查：</p>
			<h3>%s</h3>
			<p>%s</p>
			%s
			<p>请点击以下链接参与调查：</p>
			<a href="%s">%s</a>
			<p>感谢您的支持与配合！</p>
		`, survey.Title, survey.Description, customMessage, surveyLink, surveyLink)

		err := s.emailSvc.SendEmail(invitation.Email, subject, body)
		if err != nil {
			s.distRepo.UpdateInvitationStatus(invitation.ID, 5, err.Error())
			time.Sleep(30 * time.Second)
			if invitation.RetryCount < 3 {
				err = s.emailSvc.SendEmail(invitation.Email, subject, body)
				if err == nil {
					s.distRepo.UpdateInvitationStatus(invitation.ID, 2, "")
				}
			}
		} else {
			s.distRepo.UpdateInvitationStatus(invitation.ID, 2, "")
		}
	}
}

func (s *DistributionService) ListInvitations(surveyID, userID uint, page, pageSize int) ([]*model.Invitation, int64, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, 0, err
	}
	if survey == nil {
		return nil, 0, errors.New("survey not found")
	}
	if survey.UserID != userID {
		return nil, 0, errors.New("permission denied")
	}

	return s.distRepo.ListInvitations(surveyID, page, pageSize)
}

func (s *DistributionService) generateQRCode(content string) (string, error) {
	return "", nil
}

func (s *DistributionService) DeleteLink(id, userID uint) error {
	link, err := s.distRepo.FindByToken(fmt.Sprintf("%d", id))
	if err != nil {
		return err
	}
	if link == nil {
		return errors.New("link not found")
	}

	survey, err := s.surveyRepo.FindByID(link.SurveyID)
	if err != nil {
		return err
	}
	if survey == nil || survey.UserID != userID {
		return errors.New("permission denied")
	}

	return s.distRepo.DeleteLink(id)
}

func (s *DistributionService) UpdateInvitationOpened(token string) error {
	invitation, err := s.distRepo.FindInvitationByToken(token)
	if err != nil {
		return err
	}
	if invitation == nil {
		return errors.New("invitation not found")
	}

	return s.distRepo.UpdateInvitationStatus(invitation.ID, 3, "")
}
