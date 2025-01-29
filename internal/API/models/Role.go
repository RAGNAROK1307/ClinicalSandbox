package models

type Role struct {
	IDRole      uint   `json:"id_rol" gorm:"primaryKey;column:id_rol" swaggerignore:"true"`
	RoleName    string `json:"nombre_rol" gorm:"type:varchar(50);column:nombre_rol"`
	Description string `json:"descripcion" gorm:"type:text;column:description"`
}
