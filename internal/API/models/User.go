package models

type User struct {
	IDUser   uint   `json:"id_usuarios" gorm:"primaryKey;column:id_usuarios" swaggerignore:"true"`
	IDRole   uint   `json:"id_rol" gorm:"column:id_rol"` // Llave foránea que referencia la tabla roles
	UserName string `json:"nombre_usuario" gorm:"type:varchar(50);column:nombre_usuario"`
	Password string `json:"contraseña" gorm:"type:varchar(50);column:contraseña"`
	Role     Role   `gorm:"foreignKey:IDRole;references:IDRole"` // Carga la información del rol relacionado
}

func (User) TableName() string {
	return "usuarios"
}
