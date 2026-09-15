package authenticate

import (
	"fmt"
	models "maintenance-system-go/Models"
)


func (a *Authenticate) HandleLoginForm(loginRequest LoginRequest) error {
	// Implement the login logic here
	fmt.Println(loginRequest)
	var user models.User
	if err := a.Database.Where("name = ?", loginRequest.Name).First(&user).Error; err != nil {
		return err
	}
	fmt.Println(user)

	// err := bcrypt.CompareHashAndPassword([]byte(loginRequest.Password), []byte(rawPassword))
	// if err != nil {
	//     // sai mật khẩu
	// }

	return nil
}