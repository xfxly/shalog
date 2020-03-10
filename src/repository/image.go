package repository

import (
	"github.com/icowan/shalom/src/repository/types"
	"github.com/jinzhu/gorm"
)

type ImageRepository interface {
	FindByPostIdLast(postId int64) (res *types.Image, err error)
	FindByPostIds(ids []int64) (res []*types.Image, err error)
	AddImage(img *types.Image) error
	ExistsImageByMd5(val string) bool
	FindImageByMd5(val string) (img *types.Image, err error)
	FindById(id int64) (img types.Image, err error)
	FindAll(pageSize, offset int) (res []types.Image, count int64, err error)
}

type image struct {
	db *gorm.DB
}

func (c *image) FindAll(pageSize, offset int) (res []types.Image, count int64, err error) {
	err = c.db.Model(&types.Image{}).
		Select("id,image_name,image_path,post_id,created_at,image_size,image_status,client_original_mame").
		Where("post_id IS NULL OR post_id = 0").
		Count(&count).Offset(offset).Limit(pageSize).
		Order("id desc").Find(&res).Error
	return
}

func (c *image) FindById(id int64) (img types.Image, err error) {
	err = c.db.Model(&types.Image{}).Where("id = ?", id).First(&img).Error
	return
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &image{db: db}
}

func (c *image) FindByPostIdLast(postId int64) (res *types.Image, err error) {
	var i types.Image
	if err = c.db.Last(&i, "post_id=?", postId).Error; err == nil {
		return &i, nil
	}
	return
}

func (c *image) FindByPostIds(ids []int64) (res []*types.Image, err error) {
	if err = c.db.Raw("SELECT image_name,image_path,MAX(id) id,post_id,real_path FROM `images`  WHERE `images`.`deleted_at` IS NULL AND ((post_id in (?))) GROUP BY post_id ORDER BY created_at DESC", ids).
		Scan(&res).Error; err != nil {
		return
	}
	return
}

func (c *image) AddImage(img *types.Image) error {
	return c.db.Save(img).Error
}

func (c *image) ExistsImageByMd5(val string) bool {
	var img types.Image
	if err := c.db.Where("md5 = ?", val).First(&img).Error; err != nil {
		return false
	}
	if img.Md5 != "" {
		return true
	}
	return false
}

func (c *image) FindImageByMd5(val string) (img *types.Image, err error) {
	var rs types.Image
	err = c.db.Where("md5 = ?", val).First(&rs).Error
	return &rs, err
}
