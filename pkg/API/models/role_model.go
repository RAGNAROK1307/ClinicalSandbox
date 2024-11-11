package models

type Role struct {
	ID          uint   `json:"id_rol" gorm:"primaryKey;column:id_rol" swaggerignore:"true"`
	NombreRol   string `json:"nombre_rol" gorm:"type:text"`
	Descripcion string `json:"descripcion" gorm:"type:text"`
}
