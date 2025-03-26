package repository

import "github.com/o-kaisan/text-clipper/model"

type ClipRepository interface {
	FindByID() *model.Clip
	Copy() error
	Create() error
	Update() error
	ListOfActive() []*model.Clip
	ListOfInactive() []*model.Clip
	Delete() error
}
