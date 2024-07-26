package handlers

import (
	"digikala/DataBase"
	"digikala/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWTManager(manager models.Manager) (string, error) {
	claims := &jwt.MapClaims{
		"username": manager.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	newToken := models.Token{
		Token:     tokenStr,
		UserID:    manager.ID,
		ExpiresAt: expirationTime,
	}

	result := DataBase.DB.Create(&newToken)
	if result.Error != nil {
		return "", result.Error
	}

	return tokenStr, nil
}

func AuthenticateManager(username, password string) (models.Manager, string, error) {
	var manager models.Manager
	log.Println("Authenticating user:", username)
	result := DataBase.DB.Where("username = ?", username).First(&manager)
	if result.Error != nil {
		log.Println("User not found:", result.Error)
		return manager, "", echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(manager.Password), []byte(password))
	if err != nil {
		log.Println("Invalid password:", err)
		return manager, "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid password")
	}

	token, err := GenerateJWTManager(manager)
	if err != nil {
		return manager, "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to Generate Token")
	}

	return manager, token, nil
}

func CreateManager(c echo.Context) error {
	log.Println("Create Manager called")
	manager := new(models.Manager)
	if err := c.Bind(manager); err != nil {
		return err
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(manager.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"messgae": "Failed to hash password"})
	}
	manager.Password = string(hashPass)
	log.Println("the hash password is:", manager.Password)

	result := DataBase.DB.Create(manager)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	log.Println("user created")

	token, err := GenerateJWTManager(*manager)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to generate token"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfuly",
		"user":    manager,
		"token":   token,
	})
}

func LoginManager(c echo.Context) error {
	loginInfo := new(models.LoginInfo)
	if err := c.Bind(loginInfo); err != nil {
		log.Println("Failed to bind credentials:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request data"})
	}

	log.Println("Credentials received:", loginInfo)
	manager, token, err := AuthenticateManager(loginInfo.Username, loginInfo.Password)
	if err != nil {
		log.Println("Authentication failed:", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid username or password"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   token,
		"user":    manager,
	})
}

func ReadManager(c echo.Context) error {
	managerID := c.Param("id")
	var manager models.Manager
	if result := DataBase.DB.First(&manager, managerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, result.Error)
	}

	return c.JSON(http.StatusOK, manager)
}

func ReadManagers(c echo.Context) error {
	var manager []models.Manager
	result := DataBase.DB.Find(&manager)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, manager)
}

func UpdateManagers(c echo.Context) error {
	managerID := c.Param("id")
	var manager models.Manager
	if result := DataBase.DB.First(&manager, managerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, result.Error)
	}

	updatedManager := new(models.Manager)
	if err := c.Bind(updatedManager); err != nil {
		return err
	}

	if updatedManager.Name != "" {
		manager.Name = updatedManager.Name
	}


	if updatedManager.Password != "" {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(updatedManager.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		manager.Password = string(hashPass)
	}

	if result := DataBase.DB.Save(&manager); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, manager)
}

func DeleteManagers(c echo.Context) error {
	managerID := c.Param("id")

	var manager models.Manager
	if result := DataBase.DB.First(&manager, managerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "user not found"})
	}

	if result := DataBase.DB.Delete(&manager); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfuly"})

}
