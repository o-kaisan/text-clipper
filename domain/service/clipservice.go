package service

import (
	"fmt"
	"time"

	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/domain/model"
	"github.com/o-kaisan/text-clipper/domain/repository"
)

type ClipService struct {
	cr repository.ClipRepository
}

func NewClipService(cs repository.ClipRepository) ClipService {
	return ClipService{cs}
}

func (cs ClipService) DeleteClip(cid uint) error {
	clip := cs.cr.FindByID(cid)
	err := cs.cr.Delete(clip)
	if err != nil {
		return fmt.Errorf("cannot delete item: title=%s id=%d err=%w", clip.Title, clip.ID, err)
	}
	return nil
}

func (cs ClipService) GetActiveClips() ([]*model.Clip, error) {
	order := common.Env("TEXT_CLIPPER_SORT", "createdAtDesc")
	clips, err := cs.cr.ListOfActive(order)
	if err != nil {
		return nil, fmt.Errorf("cannot get all items: %w", err)
	}
	return clips, nil
}

func (cs ClipService) GetArchivedClips() ([]*model.Clip, error) {
	order := common.Env("TEXT_CLIPPER_SORT", "createdAtDesc")
	clips, err := cs.cr.ListOfInactive(order)
	if err != nil {
		return nil, fmt.Errorf("cannot get all items: %w", err)
	}
	return clips, nil
}

func (cs ClipService) RegisterClip(cid uint, title string, content string) error {
	now := time.Now()
	var err error
	if cid == 0 { // 新規登録
		newClip := model.NewClip(title, content, common.BoolPtr(true), now, now, now)
		err = cs.cr.Create(newClip)
	} else { //上書き
		targetClip := cs.cr.FindByID(cid)
		targetClip.Title = title
		targetClip.Content = content
		err = cs.cr.Update(targetClip)
	}
	if err != nil {
		return fmt.Errorf("can not save new text: %w", err)
	}
	return nil
}

func (cs ClipService) CopyClip(cid uint) error {
	return cs.cr.Copy(cid)
}

func (cs ClipService) ActivateClip(uid uint) error {
	clip := cs.cr.FindByID(uid)
	clip.IsActive = common.BoolPtr(true) // 無効化する
	err := cs.cr.Update(clip)
	if err != nil {
		return fmt.Errorf("cannot delete item: title=%s id=%d err=%w", clip.Title, clip.ID, err)
	}
	return nil
}

func (cs ClipService) DeactivateClip(cid uint) error {
	clip := cs.cr.FindByID(cid)
	clip.IsActive = common.BoolPtr(false) // 無効化する
	err := cs.cr.Update(clip)
	if err != nil {
		return fmt.Errorf("cannot deactivate clip: title=%s id=%d err=%w", clip.Title, clip.ID, err)
	}
	return nil
}

func (cs ClipService) UpdateLastUsedAt(cid uint) error {
	target := cs.cr.FindByID(cid)
	target.LastUsedAt = time.Now()
	return cs.cr.Update(target)
}
