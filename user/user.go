package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"main.go/databasehandler"
	"main.go/myStructs"
	"main.go/sendmail"
)

func UpdateProfile(c *gin.Context) {

	var user myStructs.User

	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Printf("error: %s \n ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	status, response := databasehandler.UpdateProfile(user.Profile_photo, user.Id)

	userData, getdetailsErr := databasehandler.Login(user.Email)
	if getdetailsErr != nil {
		fmt.Printf("error: %s \n ", getdetailsErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Server error"})
		return
	}

	if status == 200 {
		c.JSON(http.StatusOK, gin.H{"message": response, "user": userData})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": response})
	}
}


func RequestPromotion(c *gin.Context) {
	var requestDetails myStructs.LoginData
	if err := c.ShouldBindJSON(&requestDetails); err != nil {
		fmt.Printf("error: %s \n ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := databasehandler.Login(requestDetails.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
	} else {
		go sendmail.SendMail("charles5muchogo@gmail.com", user)
		c.JSON(http.StatusOK, gin.H{"message": "request sent successfully"})
	}

}

func PromoteUser(c *gin.Context) {
	id := c.Query("id")
	fmt.Printf("userid is updated successfully %s  \n", id)

	requestStatus := databasehandler.PromoteUser(id, true)
	fmt.Printf("userid is updated successfully %d  \n", requestStatus)

	if requestStatus == 200 {

		c.JSON(http.StatusOK, gin.H{"message": "User has been promoted to admin successfully"})

	} else {
		c.JSON(http.StatusOK, gin.H{"message": "could not promote user"})

	}
}

type DB struct {
    *gorm.DB
}

func InitDB() (*DB, error) {
    db, err := gorm.Open("postgres", databasehandler.GoDotEnvVariable("DATABASEURL"))
    if err != nil {
        return nil, err
    }
    db.AutoMigrate(&myStructs.User{})
    return &DB{db}, nil
}

func encryptPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}
func Signup(c *gin.Context) {
    db, err := InitDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer db.Close()

    var user myStructs.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Assign a default profile photo if none was provided
    if user.Profile_photo == "" {
        user.Profile_photo = "https://www.pngitem.com/pimgs/m/30-307416_profile-icon-png-image-free-download-searchpng-employee.png"
    }

    hashedPassword, err := encryptPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    user.Password = hashedPassword

    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"user": user, "message":"signup success"})
}


func Login(c *gin.Context) {
    db, err := InitDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer db.Close()

    var loginUser myStructs.LoginUser
    if err := c.ShouldBindJSON(&loginUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user myStructs.User
    if err := db.Where("email = ?", loginUser.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    c.JSON(http.StatusOK,gin.H{"user": user, "message":"login success"})
}
