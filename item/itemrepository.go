package item

import (
	"log"
	"time"

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

type ItemRepository struct {
	DB *gorm.DB
}

func (g *ItemRepository) FindByID(id uint) *Item {
	var item Item
	result := g.DB.First(&item, id)
	if result.Error != nil {
		log.Panicf("failed to get item from DB: id=%d, err=%v", id, result.Error)
		return nil
	}
	return &item
}

func (g *ItemRepository) Copy(id uint) error {
	var item Item
	// 対象を取得する
	targetItem := g.DB.First(&item, id)
	if targetItem.Error != nil {
		log.Panicf("failed to get item from DB: id=%d, err=%v", id, targetItem.Error)
		return nil
	}

	now := time.Now()
	duplicatedItem := NewItem(item.Title, item.Content, item.IsActive, now, now, now)

	// 新しく作成する
	result := g.DB.Create(&duplicatedItem)
	if result.Error != nil {
		log.Panicf("failed to create duplicated item in DB: err=%v", result.Error)
	}
	return nil
}

func (g *ItemRepository) Create(item *Item) error {
	result := g.DB.Create(item)
	if result.Error != nil {
		log.Panicf("failed to create item in DB: err=%v", result.Error)
	}
	return nil
}

func (g *ItemRepository) Update(item *Item) error {
	result := g.DB.Model(&Item{}).Where("id = ?", item.ID).Updates(item)
	if result.Error != nil {
		log.Panicf("failed to update item in DB: err=%v", result.Error)
	}
	return nil
}

func (g *ItemRepository) ListOfActive(order string) ([]*Item, error) {
	// ORDER BYの条件を気にしない場合はorderに空文字を渡す
	condition := OrderMaps[order]
	var items []*Item
	result := g.DB.Order(condition).Where("is_active = ?", true).Find(&items)
	if result.Error != nil {
		log.Panicf("failed in saving item to DB: err=%v", result.Error)
	}

	return items, nil
}

func (g *ItemRepository) ListOfInactive(order string) ([]*Item, error) {
	// ORDER BYの条件を気にしない場合はorderに空文字を渡す
	condition := OrderMaps[order]
	var items []*Item
	result := g.DB.Order(condition).Where("is_active = ?", false).Find(&items)
	if result.Error != nil {
		log.Panicf("failed in saving item to DB: err=%v", result.Error)
	}

	return items, nil
}

func (g *ItemRepository) Delete(item *Item) error {

	result := g.DB.Delete(&item)
	if result.Error != nil {
		log.Panicf("failed in deleting item from DB: err=%v", result.Error)
	}
	return nil
}
