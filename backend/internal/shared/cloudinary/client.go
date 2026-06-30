package cloudinary

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Client struct {
	cld *cloudinary.Cloudinary

	folder string
}

type UploadResult struct {
	URL string

	PublicID string

	ResourceType string

	Bytes int

	Format string
}

func New(
	cfg *config.Config,
) (*Client, error) {

	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)

	if err != nil {
		return nil, err
	}

	return &Client{

		cld: cld,

		folder: cfg.CloudinaryFolder,
	}, nil
}

func (c *Client) Upload(
	ctx context.Context,
	file interface{},
) (*UploadResult, error) {

	resp, err := c.cld.Upload.Upload(

		ctx,

		file,

		uploader.UploadParams{

			Folder: c.folder,
		},
	)

	if err != nil {
		return nil, err
	}

	return &UploadResult{

		URL: resp.SecureURL,

		PublicID: resp.PublicID,

		ResourceType: resp.ResourceType,

		Bytes: resp.Bytes,

		Format: resp.Format,
	}, nil
}
