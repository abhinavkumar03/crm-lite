package service

import (
	"context"
	"mime/multipart"

	"github.com/abhinavkumar03/crm-lite/backend/internal/media/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cloudinary"
)

type Service struct {
	client *cloudinary.Client
}

func New(
	client *cloudinary.Client,
) *Service {

	return &Service{

		client: client,
	}
}

func (s *Service) Upload(

	ctx context.Context,

	file multipart.File,

) (*dto.UploadResponse, error) {

	result, err := s.client.Upload(
		ctx,
		file,
	)

	if err != nil {
		return nil, err
	}

	return &dto.UploadResponse{

		URL: result.URL,

		PublicID: result.PublicID,

		ResourceType: result.ResourceType,

		Bytes: result.Bytes,

		Format: result.Format,
	}, nil
}
