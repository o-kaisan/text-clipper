package sqlite

import (
	"fmt"
	"log"
	"time"

	"github.com/o-kaisan/text-clipper/domain/model"
	"gorm.io/gorm"
)

var (
	// Order Byするための条件式
	CreatedAtDesc  = "created_at desc"
	UpdatedAtDesc  = "updated_at desc"
	LastUsedAtDesc = "last_used_at desc"
	CreatedAtAsc   = "created_at asc"
	UpdatedAtAsc   = "updated_at asc"
	LastUsedAtAsc  = "last_used_at asc"
)

var OrderMaps = map[string]string{
	"createdAtDesc":  CreatedAtDesc,
	"updatedAtDesc":  UpdatedAtDesc,
	"lastUsedAtDesc": LastUsedAtDesc,
	"createdAtAsc":   CreatedAtAsc,
	"updatedAtAsc":   UpdatedAtAsc,
	"lastUsedAtAsc":  LastUsedAtAsc,
}

type ClipRepositoryImpl struct {
	DB *gorm.DB
}

func NewClipRepositoryImpl(db *gorm.DB) *ClipRepositoryImpl {
	return &ClipRepositoryImpl{db}
}

func (g *ClipRepositoryImpl) FindByID(id uint) (*model.Clip, error) {
	var clip model.Clip
	result := g.DB.First(&clip, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get clip from DB: id=%d, err=%v", id, result.Error)
	}
	return &clip, nil
}

func (g *ClipRepositoryImpl) Copy(id uint) error {
	var clip model.Clip
	// 対象を取得する
	targetClip := g.DB.First(&clip, id)
	if targetClip.Error != nil {
		log.Panicf("failed to get clip from DB: id=%d, err=%v", id, targetClip.Error)
		return nil
	}

	now := time.Now()
	duplicatedClip := model.NewClip(clip.Title, clip.Content, clip.IsActive, now, now, now)

	// 新しく作成する
	result := g.DB.Create(&duplicatedClip)
	if result.Error != nil {
		log.Panicf("failed to create duplicated clip in DB: err=%v", result.Error)
	}
	return nil
}

func (g *ClipRepositoryImpl) Create(clip *model.Clip) error {
	result := g.DB.Create(clip)
	if result.Error != nil {
		log.Panicf("failed to create clip in DB: err=%v", result.Error)
	}
	return nil
}

func (g *ClipRepositoryImpl) Update(clip *model.Clip) error {
	result := g.DB.Model(&model.Clip{}).Where("id = ?", clip.ID).Updates(clip)
	if result.Error != nil {
		log.Panicf("failed to update clip in DB: err=%v", result.Error)
	}
	return nil
}

func (g *ClipRepositoryImpl) ListOfActive(order string) ([]*model.Clip, error) {
	// ORDER BYの条件を気にしない場合はorderに空文字を渡す
	condition := OrderMaps[order]
	var clips []*model.Clip
	result := g.DB.Order(condition).Where("is_active = ?", true).Find(&clips)
	if result.Error != nil {
		log.Panicf("failed in saving clip to DB: err=%v", result.Error)
	}

	return clips, nil
}

func (g *ClipRepositoryImpl) ListOfInactive(order string) ([]*model.Clip, error) {
	// ORDER BYの条件を気にしない場合はorderに空文字を渡す
	condition := OrderMaps[order]
	var clips []*model.Clip
	result := g.DB.Order(condition).Where("is_active = ?", false).Find(&clips)
	if result.Error != nil {
		log.Panicf("failed in saving clip to DB: err=%v", result.Error)
	}

	return clips, nil
}

func (g *ClipRepositoryImpl) Delete(clip *model.Clip) error {

	result := g.DB.Delete(&clip)
	if result.Error != nil {
		log.Panicf("failed in deleting clip from DB: err=%v", result.Error)
	}
	return nil
}
