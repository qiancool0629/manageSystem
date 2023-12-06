package repository

import (
	"fmt"
	"manageSystem/model"
	"manageSystem/query"
	"manageSystem/utils"

	"gorm.io/gorm"
)

type VideoRepository struct {
	DB *gorm.DB
}

type VideoRepoInterface interface {
	List(req *query.ListQuery) (Videos []*model.Video, err error)
	GetTotal() (total int64, err error)
	Get(Video model.Video) (*model.Video, error)
	Exist(Video *model.Video) *model.Video
	Add(Banner model.Video) (*model.Video, error)
	Edit(Banner model.Video) (bool, error)
	Delete(id string) (bool, error)
}

func (repo *VideoRepository) List(req *query.ListQuery) (videos []*model.Video, err error) {
	db := repo.DB
	limit, offset := utils.Page(req.PageSize, req.Page) // 分页

	if err := db.Order("id desc").Limit(limit).Offset(offset).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (repo *VideoRepository) GetTotal() (total int64, err error) {
	var videos []model.Video
	db := repo.DB
	if err := db.Find(&videos).Count(&total).Error; err != nil {
		return total, err
	}
	return total, nil
}

func (repo *VideoRepository) Get(video model.Video) (*model.Video, error) {
	if err := repo.DB.Where(&video).Find(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func (repo *VideoRepository) Exist(video *model.Video) *model.Video {
	var total int64
	repo.DB.Where(&video).Find(&video).Count(&total)
	if total > 0 {
		return video
	}
	return nil
}

func (repo *VideoRepository) Add(video model.Video) (*model.Video, error) {
	if exist := repo.Exist(&video); exist != nil {
		return nil, fmt.Errorf("视频已存在")
	}
	err := repo.DB.Create(&video).Error
	if err != nil {
		return nil, fmt.Errorf("视频注册失败")
	}
	return &video, nil
}

func (repo *VideoRepository) Edit(video model.Video) (bool, error) {
	err := repo.DB.Model(&video).Where("video_id=?", video.VideoId).Updates(map[string]interface{}{
		"video_name":   video.VideoName,
		"video_intro":  video.VideoIntro,
		"video_path":   video.VideoPath,
		"video_detail": video.VideoDetail,
		"video_tag":    video.VideoTag,
		"category_id":  video.CategoryId,
	}).Error
	//err := repo.DB.Save(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *VideoRepository) Delete(video model.Video) (bool, error) {
	err := repo.DB.Model(&video).Where("video_id=?", video.VideoId).Delete(&video).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
