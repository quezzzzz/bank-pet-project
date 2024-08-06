package bank

type Customer struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Balance  int
	Password string `json:"password" binding:"required"`
}
