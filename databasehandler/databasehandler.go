package databasehandler

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"main.go/myStructs"
)

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	CheckError(err)

	return os.Getenv(key)
}

func DbConnect() *sql.DB {
	dsn := GoDotEnvVariable("DATABASEURL")

	db, err := sql.Open("postgres", dsn)

	CheckError(err)
	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func BeforeSave(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}


func SaveUser(firstName string, middleName string, email string, firebase_id string, phoneNumber string, password string) (myStructs.User, int, error) {
	userUrl := "INSERT INTO users(first_name,middle_name, email, phone_number, firebase_id , password) VALUES($1, $2, $3, $4, $5, $6)"
	status := 500

	var userDetails myStructs.User

	mypassword, _ := BeforeSave(password)
	insertUser, err := DbConnect().Exec(userUrl, firstName, middleName, email, phoneNumber, firebase_id, mypassword)
	defer DbConnect().Close()
	if err != nil {
		return userDetails, status, err
	}

	affected, affectederr := insertUser.RowsAffected()

	if affectederr != nil {
		return userDetails, status, err
	}

	if affected > 0 {
		var loginError error = nil
		userDetails, loginError = Login(email)

		if loginError != nil {
			return userDetails, status, loginError
		}
	}

	status = 200

	return userDetails, status, nil

}

func Login(email string) (myStructs.User, error) {
	var data myStructs.User

	loginQuery := fmt.Sprintf("SELECT  id, is_admin, first_name, middle_name, email, phone_number, password, profile_photo  FROM users WHERE email = '%v'  ", email)
	fmt.Printf("login querry is: %s \n", loginQuery)

	rows, err := DbConnect().Query(loginQuery)

	if err != nil {
		fmt.Printf("login querry is: %s \n", err.Error())
		return data, err
	}

	for rows.Next() {
		err = rows.Scan(&data.Id, &data.Is_admin, &data.First_name, &data.Middle_name, &data.Email, &data.Phone_number, &data.Password, &data.Profile_photo)
		CheckError(err)
	}

	return data, nil
}

func UpdateLocation(user_id string, cur_lat string, curr_lng string, max_dis string, orig_lat string, orig_lng string, user_distance float32) (int, string) {
	var loginQuery string

	if cur_lat == "" && curr_lng == "" && orig_lat != "0.0" && orig_lng != "0.0" {
		fmt.Printf("if 1: \n")
		loginQuery = fmt.Sprintf("Update distance set max_distance = '%v', origin_latitude = '%v', origin_longitude = '%v', notification_sent = '%v', latest_update = CURRENT_TIMESTAMP  WHERE user_id = '%v'", max_dis, orig_lat, orig_lng, false, user_id)
	}

	if cur_lat != "" && curr_lng != "" && orig_lat != "0.0" && orig_lng != "0.0" {
		fmt.Printf("if 2: \n")
		loginQuery = fmt.Sprintf("Update distance set current_latitude = '%v', current_longitude = '%v', latest_update = CURRENT_TIMESTAMP, user_distance = '%v'  WHERE user_id = '%v'", cur_lat, curr_lng, user_distance, user_id)
	}
	if cur_lat == "" && curr_lng == "" && orig_lat == "0.0" && orig_lng == "0.0" {
		fmt.Printf("here ---------------------------------------------> 3: \n")
		loginQuery = fmt.Sprintf("Update distance set max_distance = '%v', notification_sent = '%v', latest_update = CURRENT_TIMESTAMP  WHERE user_id = '%v'", max_dis, false, user_id)
	}

	status := 500
	dbRResponse := "failed to update location"
	fmt.Printf("update  querry is: ---------------------------------------------> %s \n", loginQuery)

	rows, err := DbConnect().Exec(loginQuery)

	if err != nil {
		fmt.Printf("login querry is: %s \n", err.Error())
		return status, err.Error()
	}

	rows_affected, _ := rows.RowsAffected()

	if rows_affected > 0 {
		status = 200
		dbRResponse = "update location successfully"
	} else {
		status, dbRResponse = addLocationuser_id(user_id, cur_lat, curr_lng, max_dis, orig_lat, orig_lng)
	}

	return status, dbRResponse
}

func GetUserLocationDetails(user_id string) (int, myStructs.MyLocation) {
	status := 500
	var myLocation myStructs.MyLocation
	query := "SELECT max_distance, origin_latitude, origin_longitude from distance where user_id = $1"

	fmt.Printf("login querry is: %s \n", query)

	rows, err := DbConnect().Query(query, user_id)
	if err != nil {
		fmt.Printf("get location error is: %s \n", err.Error())
		return status, myLocation
	}

	for rows.Next() {
		err = rows.Scan(&myLocation.MaxDistance, &myLocation.OriginLatitude, &myLocation.OriginLongitude)
		fmt.Printf("latitude is: %s \n", myLocation.OriginLatitude)
		fmt.Printf("longitude is: %s \n", myLocation.OriginLongitude)
		CheckError(err)
		status = 200

	}

	return status, myLocation
}

func addLocationuser_id(user_id string, cur_lat string, curr_lng string, max_dis string, orig_lat string, orig_lng string) (int, string) {
	status := 500
	dbRResponse := "failed to update location"

	locationQuery := "INSERT INTO distance( user_id, current_latitude, current_longitude, max_distance, origin_latitude, origin_longitude, user_distance )VALUES ($1, $2, $3, $4, $5, $6 ,$7)"
	fmt.Printf("insert querry is: %s \n", locationQuery)

	rows, err := DbConnect().Exec(locationQuery, user_id, cur_lat, curr_lng, 10.0, cur_lat, curr_lng, 0.00)

	if err != nil {
		fmt.Printf("login querry is: %s \n", err.Error())
		return status, err.Error()
	}
	rowsAffected, err := rows.RowsAffected()
	CheckError(err)

	if rowsAffected > 0 {
		status = 200
		dbRResponse = " update location successfully"
	}

	return status, dbRResponse
}

func UpdateProfile(image string, id int) (int, string) {
	status := 500
	dbRResponse := "failed to update profile"

	profileQuery := "UPDATE users SET profile_photo = $1 WHERE id = $2"
	fmt.Printf("update  querry is: %s \n", profileQuery)

	rows, err := DbConnect().Exec(profileQuery, image, id)

	if err != nil {
		fmt.Printf("update profile error is: %s \n", err.Error())
		return status, err.Error()
	}
	rowsAffected, err := rows.RowsAffected()
	CheckError(err)

	if rowsAffected > 0 {
		status = 200
		dbRResponse = "update profile successfully"
	}

	return status, dbRResponse
}
func PromoteUser(user_id string, status bool) (resultStatus int) {
	requestStatus := 500
	query := "UPDATE users SET is_admin = $1 WHERE id = $2"

	fmt.Printf("update  querry is: %s \n", query)

	rows, err := DbConnect().Exec(query, status, user_id)

	if err != nil {
		fmt.Printf("update profile error is: %s \n", err.Error())
		return requestStatus
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		fmt.Printf("update profile error is: %s \n", err.Error())
		return requestStatus
	}

	if rowsAffected > 0 {
		requestStatus = 200

	}

	return requestStatus
}

func GetUsersLocation() ([]myStructs.LocationUpdate, int) {

	query := "SELECT  u.first_name,     u.middle_name, u.phone_number,u.email, u.id, d.user_distance, d.current_latitude, d.current_longitude,d.max_distance,d.origin_latitude,d.origin_longitude,d.latest_update FROM users u INNER JOIN distance d ON  u.id = d.user_id;"
	rows, err := DbConnect().Query(query)
	defer DbConnect().Close()
	CheckError(err)

	response := 500

	var currentUser myStructs.LocationUpdate
	var userSlice []myStructs.LocationUpdate

	for rows.Next() {
		err = rows.Scan(&currentUser.FirstName, &currentUser.MiddleName, &currentUser.PhoneNumber, &currentUser.Email, &currentUser.UserId, &currentUser.User_distance, &currentUser.CurrentLatitude, &currentUser.CurrentLongitude, &currentUser.MaxDistance, &currentUser.OriginLatitude, &currentUser.OriginLongitude, &currentUser.LastUpdate)
		CheckError(err)

		fmt.Printf("update  querry is: %s \n", currentUser.FirstName)
		userSlice = append(userSlice, currentUser)
		response = 200

	}

	return userSlice, response
}

func GetFcmDetails(id string) (bool, string, string) {

	query := "SELECT d.notification_sent, u.first_name, u.middle_name FROM distance d INNER JOIN users u ON d.user_id = u.id WHERE u.id = $1"

	fmt.Printf("update  querry is: %s \n", query)

	rows, err := DbConnect().Query(query, id)
	CheckError(err)

	var fcmModel myStructs.FcmModel

	for rows.Next() {
		err = rows.Scan(&fcmModel.IsNotificationSent, &fcmModel.FirstName, &fcmModel.LastName)
		CheckError(err)
	}

	notificationsent := fcmModel.IsNotificationSent
	firstName := fcmModel.FirstName
	lastName := fcmModel.LastName

	return notificationsent, firstName, lastName

}
func UpdateNotificationSent(id string, notification_sent bool) {
	notificationQuery := "UPDATE distance SET notification_sent = $1 WHERE user_id = $2"
	fmt.Printf("update  querry is: %s \n", notificationQuery)

	rows, err := DbConnect().Exec(notificationQuery, notification_sent, id)

	rowsAffected, err := rows.RowsAffected()
	CheckError(err)

	if rowsAffected > 0 {
		fmt.Printf("notification updated successfully \n \n")
	}
}
