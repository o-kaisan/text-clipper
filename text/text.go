package text

import (
	"log"

	"gorm.io/gorm"
)

type Text struct {
	gorm.Model
	Title    string
	Contents string
}

type Repository interface {
	Get(id uint) error
	Save(text *Text) error
	Update(text *Text) error
	List() ([]*Text, error)
	Delete(text Text) error
}

type GormRepository struct {
	DB *gorm.DB
}

// TODO Error返せるように
func (g *GormRepository) FindByID(id uint) *Text {
	var text Text
	result := g.DB.First(&text, id)
	if result.Error != nil {
		log.Panicf("failed to get text from DB: id=%d, err=%v", id, result.Error)
		return nil
	}
	return &text
}

// TODO Error返せるように
func (g *GormRepository) Crete(text *Text) error {
	result := g.DB.Create(text)
	if result.Error != nil {
		log.Panicf("failed to create text in DB: err=%v", result.Error)
	}
	return nil
}

// TODO Error返せるように
func (g *GormRepository) Update(text *Text) error {
	result := g.DB.Model(&Text{}).Where("id = ?", text.ID).Updates(map[string]interface{}{
		"Title":    text.Title,
		"Contents": text.Contents,
	})
	if result.Error != nil {
		log.Panicf("failed to update text in DB: err=%v", result.Error)
	}
	return nil
}

// TODO Error返せるように
func (g *GormRepository) List() ([]*Text, error) {
	var texts []*Text
	result := g.DB.Find(&texts)
	if result.Error != nil {
		log.Panicf("failed in saving text to DB: err=%v", result.Error)
	}

	return texts, nil
}

// TODO Error返せるように
func (g *GormRepository) Delete(text *Text) error {

	result := g.DB.Delete(&text)
	if result.Error != nil {
		log.Panicf("failed in deleting text from DB: err=%v", result.Error)
	}
	return nil
}
