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

// DeleteClip は、指定したIDのクリップを削除(物理削除)します。
func (cs ClipService) DeleteClip(cid uint) error {
	var err error
	clip, err := cs.cr.FindByID(cid)
	if err != nil {
		return fmt.Errorf("cannot delete clip: title=%s, id=%d err=%w", clip.Title, clip.ID, err)
	}

	err = cs.cr.Delete(clip)
	if err != nil {
		return fmt.Errorf("cannot delete clip: title=%s id=%d err=%w", clip.Title, clip.ID, err)
	}
	return nil
}

// GetActiveClips は、アクティブなクリップのリストを取得します。
func (cs ClipService) GetActiveClips() ([]*model.Clip, error) {
	order := common.Env("TEXT_CLIPPER_SORT", "createdAtDesc")
	clips, err := cs.cr.ListOfActive(order)
	if err != nil {
		return nil, fmt.Errorf("cannot get active clips: %w", err)
	}
	return clips, nil
}

// GetArchivedClips は、アーカイブされたクリップのリストを取得します。
func (cs ClipService) GetArchivedClips() ([]*model.Clip, error) {
	order := common.Env("TEXT_CLIPPER_SORT", "createdAtDesc")
	clips, err := cs.cr.ListOfInactive(order)
	if err != nil {
		return nil, fmt.Errorf("cannot get archived clips: %w", err)
	}
	return clips, nil
}

// RegisterClip は、クリップを登録または更新します。
func (cs ClipService) RegisterClip(cid uint, title string, content string) error {
	var err error
	now := time.Now()
	if cid == 0 { // 新規登録
		newClip := model.NewClip(title, content, common.BoolPtr(true), now, now, now)
		err = cs.cr.Create(newClip)
	} else { //上書き
		targetClip, err := cs.cr.FindByID(cid)
		if err != nil {
			return fmt.Errorf("cannot update clip: title=%s, id=%d err=%w", title, cid, err)
		}
		targetClip.Title = title
		targetClip.Content = content
		err = cs.cr.Update(targetClip)
	}
	if err != nil {
		return fmt.Errorf("can not save new text: %w", err)
	}
	return nil
}

// CopyClip は、指定したIDのクリップをコピーします。
func (cs ClipService) CopyClip(cid uint) error {
	return cs.cr.Copy(cid)
}

// ActivateClip は、指定したIDのクリップを有効化します。
func (cs ClipService) ActivateClip(uid uint) error {
	var err error
	clip, err := cs.cr.FindByID(uid)
	if err != nil {
		return fmt.Errorf("cannot activate clip: title=%s, id=%d err=%w", clip.Title, clip.ID, err)
	}

	clip.IsActive = common.BoolPtr(true) // 無効化する
	err = cs.cr.Update(clip)
	if err != nil {
		return fmt.Errorf("cannot delete item: title=%s id=%d err=%w", clip.Title, clip.ID, err)
	}
	return nil
}

// DeactivateClip は、指定したIDのクリップを無効化(論理削除)します。
func (cs ClipService) DeactivateClip(cid uint) error {
	var err error
	clip, err := cs.cr.FindByID(cid)
	if err != nil {
		return fmt.Errorf("cannot deactivate clip: title=%s, id=%d err=%w", clip.Title, cid, err)
	}
	clip.IsActive = common.BoolPtr(false) // 無効化する
	err = cs.cr.Update(clip)
	if err != nil {
		return fmt.Errorf("cannot deactivate clip: title=%s id=%d err=%w", clip.Title, clip.ID, err)
	}
	return nil
}

// UpdateLastUsedAt は、指定したIDのクリップの最終利用日時を更新します。
func (cs ClipService) UpdateLastUsedAt(cid uint) error {
	target, err := cs.cr.FindByID(cid)
	if err != nil {
		return fmt.Errorf("could not update clip: id=%d err=%w", cid, err)
	}
	target.LastUsedAt = time.Now()
	return cs.cr.Update(target)
}
