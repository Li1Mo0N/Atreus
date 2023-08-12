package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"math"
)

// Video is a video model
type Video struct {
	ID            uint32
	Author        *User
	PlayURL       string
	CoverURL      string
	FavoriteCount uint32
	CommentCount  uint32
	IsFavorite    bool
	Title         string
}

// User is a user model.
type User struct {
	ID              uint32
	Name            string
	FollowCount     uint32
	FollowerCount   uint32
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  uint32
	WorkCount       uint32
	FavoriteCount   uint32
}

// PublishRepo is a publishing repo.
type PublishRepo interface {
	FindUserByID(context.Context, uint32) (*User, error)
	CreateVideo(context.Context, Video) error
	FindVideoListByUserID(context.Context, uint32) ([]*Video, error)
	FindVideoListByIDs(context.Context, []uint32) ([]*Video, error)
	FindVideoListByCount(context.Context, string, uint32) ([]*Video, error)
	UpdateVideo(context.Context, uint32, string, uint32, bool) error
}

// PublishUsecase is a publishing usecase.
type PublishUsecase struct {
	repo PublishRepo
	log  *log.Helper
}

// NewPublishUsecase new a publishing usecase.
func NewPublishUsecase(repo PublishRepo, logger log.Logger) *PublishUsecase {
	return &PublishUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Auth use jwt verified token to get user id
func Auth(ctx context.Context, t0ken string) (*User, bool) {
	// TODO

	return &User{}, true
}

// UploadVideo use ftp, ffmpeg to get play url and cover url
func UploadVideo(ctx context.Context, video []byte) (string, string) {
	// TODO

	return "", ""
}

// PublishVideo auth user, upload video and save video
func (u PublishUsecase) PublishVideo(ctx context.Context, t string, title string, videoData []byte) error {
	user, ok := Auth(ctx, t)
	if !ok {
		return nil
	}
	playUrl, coverUrl := UploadVideo(ctx, videoData)
	video := Video{
		Author:   user,
		Title:    title,
		PlayURL:  playUrl,
		CoverURL: coverUrl,
	}
	err := u.repo.CreateVideo(ctx, video)

	return err
}

// GetVideoListByUserID get video list by id
func (u PublishUsecase) GetVideoListByUserID(ctx context.Context, userID uint32) ([]*Video, error) {
	videoList, err := u.repo.FindVideoListByUserID(ctx, userID)

	return videoList, err
}

func (u PublishUsecase) GetVideoListByIDs(ctx context.Context, ids []uint32) ([]*Video, error) {
	videoList, err := u.repo.FindVideoListByIDs(ctx, ids)

	return videoList, err
}

func (u PublishUsecase) GetVideoListByCount(ctx context.Context, time string, num uint32) ([]*Video, error) {
	// TODO

	return nil, nil
}

// UpdateVideoCount update video comment/favorite count
func (u PublishUsecase) UpdateVideoCount(ctx context.Context, videoID uint32, field string, count int32) (err error) {
	var uintCount uint32
	if count < 0 {
		uintCount = uint32(count)                  // 转换为 uint32，将得到一个正数
		uintCount = math.MaxUint32 - uintCount + 1 // 补码求反
		err = u.repo.UpdateVideo(ctx, videoID, field, uintCount, false)
	} else {
		err = u.repo.UpdateVideo(ctx, videoID, field, uint32(count), true)
	}

	return err
}
