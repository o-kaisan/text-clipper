package di

import (
	"github.com/o-kaisan/text-clipper/infrastructure"
	"github.com/o-kaisan/text-clipper/service"
	"gorm.io/gorm"
)

type Container struct {
	// リポジトリ
	Cr *infrastructure.ClipRepositoryImpl
	// サービス
	Cs service.ClipService
}

func NewContainer(db *gorm.DB) Container {
	cri := infrastructure.NewClipRepositoryImpl(db)
	cs := service.ClipService{}

	return Container{
		Cr: cri,
		Cs: cs,
	}
}
