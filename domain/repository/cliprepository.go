package repository

import "github.com/o-kaisan/text-clipper/domain/model"

type ClipRepository interface {
	FindByID(id uint) *model.Clip
	Copy(id uint) error
	Create(clip *model.Clip) error
	Update(clip *model.Clip) error
	ListOfActive(order string) ([]*model.Clip, error)
	ListOfInactive(order string) ([]*model.Clip, error)
	Delete(clip *model.Clip) error
}
