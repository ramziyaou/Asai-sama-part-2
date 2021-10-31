package entity

// import (
// 	// "github.com/jinzhu/gorm"
// )

//Person object for REST(CRUD)
type User struct {
	// *gorm.Model
	ID       int    `gorm:"primaryKey;autoIncrement:false"`
	Username string `gorm:"foreignKey:UserRefer"`
	Password string
	// Wallets  WalletArray `gorm:"column:wallets;type:longtext"`
	Wallets string
}

// type WalletArray []CryptoWallet

// func (sla *WalletArray) Scan(src interface{}) error {
// 	return json.Unmarshal(src.([]byte), &sla)
// }

// func (sla WalletArray) Value() (driver.Value, error) {
// 	val, err := json.Marshal(sla)
// 	return string(val), err
// }
