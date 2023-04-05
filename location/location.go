package location

import (
	"fmt"
	"net/http"


	//"strconv"

	"github.com/gin-gonic/gin"
	"main.go/databasehandler"
	"main.go/fcm"
	"main.go/myStructs"
)

func UpdateLocation(c *gin.Context) {
	var location myStructs.LocationUpdate

	if err := c.ShouldBindJSON(&location); err != nil {
		fmt.Printf("error: %s \n ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fmt.Printf("current distance from origin is: %f \n ", location.User_distance)


	if location.CurrentLatitude != "" && location.CurrentLongitude != "" {
		go handleNotifications(location.UserId, location.User_distance, float32(location.MaxDistance))
	}
	
	

	status, response := databasehandler.UpdateLocation(location.UserId, location.CurrentLatitude, location.CurrentLongitude, fmt.Sprintf("%f", location.MaxDistance) , location.OriginLatitude, location.OriginLongitude, location.User_distance)

	getUserLocation, mylocation :=  databasehandler.GetUserLocationDetails(location.UserId);

	devices, devicesStatus := databasehandler.GetUsersLocation()

	if devicesStatus == 200 {
		if getUserLocation == 200 {
			if status == 200 {
				c.JSON(http.StatusOK, gin.H{"message": response, "users": devices, "you": mylocation })
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": response})
			}
		}
		
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "server error"})
	}
}

func GetLocation(c *gin.Context) {

	userLocation, status := databasehandler.GetUsersLocation()

	if status == 200 {
		c.JSON(http.StatusOK, gin.H{
			"message": "get devices location success",
			"devices": userLocation,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
	}
}

func handleNotifications(user_id string, user_distance float32, max_distance float32) {

	if user_distance > max_distance {
		IsNotificationSent, firstName, middleName := databasehandler.GetFcmDetails(user_id)

		if IsNotificationSent != true {
			fcm.SendNotification( user_id, firstName, middleName)
			databasehandler.UpdateNotificationSent(user_id, true)

		}

	} else {
		databasehandler.UpdateNotificationSent(user_id, false)
	}
}
