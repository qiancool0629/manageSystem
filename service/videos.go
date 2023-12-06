package service

import (
	"errors"
	"fmt"
	"manageSystem/config"
	"manageSystem/model"
	"manageSystem/query"
	"manageSystem/repository"

	uuid "github.com/satori/go.uuid"
)

type VideoRepoSrv interface {
	List(req *query.ListQuery) (Videos []*model.Video, err error)
	GetTotal() (total int64, err error)
	Get(Video model.Video) (*model.Video, error)
	Exist(Video *model.Video) *model.Video
	Add(Video model.Video) (*model.Video, error)
	Edit(Video model.Video) (bool, error)
	Delete(id string) (bool, error)
}

type VideoService struct {
	Repo repository.VideoRepoInterface
}

func (srv *VideoService) List(req *query.ListQuery) (Videos []*model.Video, err error) {
	if req.PageSize < 1 {
		req.PageSize = config.PAGE_SIZE
	}
	return srv.Repo.List(req)
}

func (srv *VideoService) GetTotal() (total int64, err error) {
	return srv.Repo.GetTotal()
}

func (srv *VideoService) Get(video model.Video) (*model.Video, error) {
	return srv.Repo.Get(video)
}

func (srv *VideoService) Exist(video *model.Video) *model.Video {
	return srv.Repo.Exist(video)
}

func (srv *VideoService) Add(video model.Video) (*model.Video, error) {
	if video.VideoPath == "" || video.VideoName == "" {
		return nil, errors.New("请输入视频名称或存放地址")
	}
	nameResult := srv.Repo.Exist(&video)
	if nameResult != nil {
		return nil, errors.New("视频名称或地址已经存在")
	}
	video.VideoId = uuid.NewV4().String()
	return srv.Repo.Add(video)
}

func (srv *VideoService) Edit(video model.Video) (bool, error) {
	if video.VideoId == "" {
		return false, fmt.Errorf("参数错误")
	}
	exist := srv.Repo.Exist(&video)
	if exist == nil {
		return false, errors.New("参数错误")
	}
	exist.VideoName = video.VideoName
	exist.VideoDetail = video.VideoDetail
	exist.VideoIntro = video.VideoIntro
	exist.VideoPath = video.VideoPath
	exist.VideoTag = video.VideoTag
	return srv.Repo.Edit(*exist)
}

func (srv *VideoService) Delete(id string) (bool, error) {
	if id == "" {
		return false, errors.New("参数错误")
	}
	v := model.Video{
		VideoId: id,
	}

	video := srv.Exist(&v)
	if video == nil {
		return false, errors.New("参数错误")
	}
	return srv.Repo.Delete(video.VideoId)
}
