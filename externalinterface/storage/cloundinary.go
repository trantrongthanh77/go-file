package storage

import (
	"context"
	"fmt"
	"go-file/common/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary interface {
	GetSubFolder(parent string) *admin.FoldersResult
	GetAssets(path string) *admin.SearchResult
	UpdateAsset(asset interface{}, params uploader.UploadParams) (*uploader.UploadResult, error)
}

type CloudinaryImpl struct {
	cloundinary *cloudinary.Cloudinary
	ctx         context.Context
}

var _ Cloudinary = &CloudinaryImpl{}

func NewCloudinary(conf *config.Config) *CloudinaryImpl {
	cld, _ := cloudinary.NewFromURL(conf.CloudinaryURL)
	cld.Config.URL.Secure = true
	ctx := context.Background()
	return &CloudinaryImpl{
		cloundinary: cld,
		ctx:         ctx,
	}
}

func (c *CloudinaryImpl) GetSubFolder(parent string) *admin.FoldersResult {
	if parent == "" {
		parent = "samples"
	}
	folders, err := c.cloundinary.Admin.SubFolders(c.ctx, admin.SubFoldersParams{Folder: parent})
	if err != nil {
		fmt.Println("error")
	}
	return folders
}

func (c *CloudinaryImpl) GetAssets(path string) *admin.SearchResult {
	asset, err := c.cloundinary.Admin.Search(c.ctx, search.Query{
		Expression: fmt.Sprintf("folder:%s", path),
	})
	if err != nil {
		fmt.Println("error", err)
	}
	return asset
}

func (c *CloudinaryImpl) UpdateAsset(asset interface{}, params uploader.UploadParams) (*uploader.UploadResult, error) {
	result, err := c.cloundinary.Upload.Upload(c.ctx, asset, params)
	if err != nil {
		fmt.Println("error", err)
	}
	return result, err
}
