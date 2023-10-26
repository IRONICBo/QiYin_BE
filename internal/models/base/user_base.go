package base

// UserBase user base model.
type UserBase struct {
	Model

	UUID        string `json:"uuid"     gorm:"index;type:varchar(36);not null;comment:UUID"`
	Email       string `json:"email"    gorm:"type:varchar(255);not null;unique;comment:Email"`
	Nickname    string `json:"nickname" gorm:"type:varchar(20);not null;comment:Nickname"`
	Avatar      string `json:"avatar"   gorm:"type:varchar(255);not null;comment:Avatar"`
	Description string `json:"description" gorm:"type:varchar(255);not null;comment:Description"`
	IsEnable    bool   `json:"is_enable"   gorm:"type:tinyint(1);not null;default:1;comment:IsEnable,1:enable,0:disable"`
}
