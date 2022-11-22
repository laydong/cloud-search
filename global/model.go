package global

import (
	"gorm.io/gorm"
)

type ModelID struct {
	ID uint `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 主键ID
}

type ModelTime struct {
	CreatedAt Time           `gorm:"column:created_at;NOT NULL;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt Time           `gorm:"column:updated_at;NOT NULL;comment:更新时间" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`                               // 删除时间
}
