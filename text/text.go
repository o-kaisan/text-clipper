package text

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Text struct {
	gorm.Model
	Title      string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastUsedAt time.Time
}

var (
	// Order Byするための条件式
	createdAtDesc  = "created_at desc"
	updatedAtDesc  = "updated_at desc"
	lastUsedAtDesc = "last_used_at desc"
	createdAtAsc   = "created_at asc"
	updatedAtAsc   = "updated_at asc"
	lastUsedAtAsc  = "last_used_at asc"
)

var orderMap = map[string]string{
	"createdAtDesc":  createdAtDesc,
	"updatedAtDesc":  updatedAtDesc,
	"lastUsedAtDesc": lastUsedAtDesc,
	"createdAtAsc":   createdAtAsc,
	"updatedAtAsc":   updatedAtAsc,
	"lastUsedAtAsc":  lastUsedAtAsc,
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

func (g *GormRepository) FindByID(id uint) *Text {
	var text Text
	result := g.DB.First(&text, id)
	if result.Error != nil {
		log.Panicf("failed to get text from DB: id=%d, err=%v", id, result.Error)
		return nil
	}
	return &text
}

func (g *GormRepository) Crete(text *Text) error {
	result := g.DB.Create(text)
	if result.Error != nil {
		log.Panicf("failed to create text in DB: err=%v", result.Error)
	}
	return nil
}

func (g *GormRepository) Update(text *Text) error {
	result := g.DB.Model(&Text{}).Where("id = ?", text.ID).Updates(map[string]interface{}{
		"Title":   text.Title,
		"Content": text.Content,
	})
	if result.Error != nil {
		log.Panicf("failed to update text in DB: err=%v", result.Error)
	}
	return nil
}

// ORDER BYの条件を気にしない場合は空文字を渡す
func (g *GormRepository) List(order string) ([]*Text, error) {
	condition := orderMap[order]
	var texts []*Text
	result := g.DB.Order(condition).Find(&texts)
	if result.Error != nil {
		log.Panicf("failed in saving text to DB: err=%v", result.Error)
	}

	return texts, nil
}

func (g *GormRepository) Delete(text *Text) error {

	result := g.DB.Delete(&text)
	if result.Error != nil {
		log.Panicf("failed in deleting text from DB: err=%v", result.Error)
	}
	return nil
}
