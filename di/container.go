package di

import (
	"github.com/o-kaisan/text-clipper/domain/service"
	"github.com/o-kaisan/text-clipper/infrastructure/sqlite"
	"gorm.io/gorm"
)

type Container struct {
	// リポジトリ
	Cr *sqlite.ClipRepositoryImpl
	// サービス
	Cs service.ClipService
}

func NewContainer(db *gorm.DB) Container {
	cri := sqlite.NewClipRepositoryImpl(db)
	cs := service.NewClipService(cri)

	return Container{
		Cr: cri,
		Cs: cs,
	}
}
