package data

import (
	"Atreus/app/publish/service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// Video Database Model
type Video struct {
	Id            uint32 `gorm:"id"`
	AuthorID      uint32 `gorm:"column:author_id;not null"`
	Title         string `gorm:"column:title;not null;size:255"`
	PlayURL       string `gorm:"column:play_url;not null"`
	CoverURL      string `gorm:"column:cover_url;not null"`
	FavoriteCount uint32 `gorm:"column:favorite_count;not null"`
	CommentCount  uint32 `gorm:"column:comment_count;not null"`
	CreateAt      uint64 `gorm:"column:create_at"`
}

type publishRepo struct {
	data *Data
	log  *log.Helper
}

// NewPublishRepo .
func NewPublishRepo(data *Data, logger log.Logger) biz.PublishRepo {
	repo := &publishRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
	return repo
}

func (r *publishRepo) FindUserByID(ctx context.Context, id uint32) (*biz.User, error) {
	user := &biz.User{
		Signature: "this is June",
	}
	//err := r.data.db.Take(user, "id=?", id).Error

	return user, nil
}

func (r *publishRepo) CreateVideo(ctx context.Context, video biz.Video) error {
	v := &Video{
		Id:            video.ID,
		AuthorID:      video.Author.ID,
		Title:         video.Title,
		PlayURL:       video.PlayURL,
		CoverURL:      video.CoverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		CreateAt:      uint64(time.Now().Unix()),
	}
	r.data.db.WithContext(ctx).Create(v)

	return nil
}

func (r *publishRepo) FindVideoListByUserID(ctx context.Context, userId uint32) ([]*biz.Video, error) {
	user, _ := r.FindUserByID(ctx, userId)
	var videoList []*Video
	var vl []*biz.Video
	err := r.data.db.WithContext(ctx).Where("author_id = ?", userId).Find(&videoList).Error
	for _, video := range videoList {
		vl = append(vl, &biz.Video{
			ID:            video.Id,
			Author:        user,
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}

	return vl, err
}

func (r *publishRepo) FindVideoListByIDs(ctx context.Context, ids []uint32) ([]*biz.Video, error) {
	var videoList []*Video
	var vl []*biz.Video
	err := r.data.db.WithContext(ctx).Find(&videoList, ids).Error
	for _, video := range videoList {
		user, _ := r.FindUserByID(ctx, video.AuthorID)
		vl = append(vl, &biz.Video{
			ID:            video.Id,
			Author:        user,
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}

	return vl, err
}

func (r *publishRepo) UpdateVideo(ctx context.Context, videoID uint32, field string, count uint32, isPositive bool) (err error) {
	var video Video
	var newCount uint32

	if field == "favorite" {
		err = r.data.db.WithContext(ctx).Select("favorite_count").Where("id = ?", videoID).First(&video).Error
		if isPositive {
			newCount = video.FavoriteCount + count
		} else {
			newCount = video.FavoriteCount - count
		}
		err = r.data.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", newCount).Error
	} else { // field == "comment"
		err = r.data.db.WithContext(ctx).Model(&Video{}).Select("comment_count").Where("id = ?", videoID).First(&video).Error
		if isPositive {
			newCount = video.CommentCount + count
		} else {
			newCount = video.CommentCount - count
		}
		err = r.data.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("comment_count", newCount).Error
	}

	return err
}

func (r *publishRepo) FindVideoListByCount(context.Context, string, uint32) ([]*biz.Video, error) {
	var vl []*biz.Video
	var err error

	return vl, err
}
