package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	// "strconv"
	"time"

	"github.com/Laevateinn17/travelohi-backend/db"
	"github.com/Laevateinn17/travelohi-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "vincent ganteng"

const recaptchaSecretKey = "6LehlFopAAAAANZItAlnBGdpWPkBm634fN5wuOLr"

func VerifyRecaptcha(responseToken string) (bool, error) {
	// Prepare the request data
	data := url.Values{}
	data.Set("secret", recaptchaSecretKey)
	data.Set("response", responseToken)

	// Make a POST request to the reCAPTCHA verification endpoint
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", data)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Parse the JSON response
	var recaptchaResponse map[string]interface{}
	if err := json.Unmarshal(body, &recaptchaResponse); err != nil {
		return false, err
	}

	// Check if the verification was successful
	success, ok := recaptchaResponse["success"].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected response format")
	}

	return success, nil
}

func HandleRegister(c *fiber.Ctx) error {

	var data struct {
		User     models.User     `json:"user"`
		UserAuth models.UserAuth `json:"userAuth"`
		Captcha  string          `json:"captcha"`
	}

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Error binding JSON: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Invalid JSON"})
	}

	success, err := VerifyRecaptcha(data.Captcha)

	if err != nil || !success {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	database, _ := db.Connect()
	err = RegisterUser(database, &data.User, &data.UserAuth)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = SendEmail(&SMTP_SERVER, Email, []string{data.UserAuth.Email}, EmailPassword, "Account Registered Successfully", "Your account is registered successfully.\n")

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}

func Ping(c *fiber.Ctx) error {
	c.SendStatus(http.StatusOK)
	return nil
}

func HandleLogin(c *fiber.Ctx) error {
	var data struct {
		UserAuth *models.UserAuth `json:"userAuth"`
		Captcha  string           `json:"captcha"`
	}

	if err := c.BodyParser(&data); err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Invalid JSON"})
	}

	success, err := VerifyRecaptcha(data.Captcha)

	if err != nil || !success {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	database, _ := db.Connect()
	userAuth, err := GetUserAuth(database, data.UserAuth)
	if err != nil {
		fmt.Println(1)
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	expiryTime := time.Now().Add(time.Hour * 24)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(userAuth.ID)),
		ExpiresAt: jwt.NewNumericDate(expiryTime),
	})

	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		fmt.Println(3)
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "could not login",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expiryTime,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	cookie2 := fiber.Cookie{
		Name:    "session",
		Value:   "active",
		Expires: expiryTime,
	}

	c.Cookie(&cookie2)
	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"messsage": "success",
	})
}

func HandleLoginByEmail(c *fiber.Ctx) error {
	var userAuth *models.UserAuth

	if err := c.BodyParser(&userAuth); err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Invalid JSON"})
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userAuth, err = GetUserAuthByEmail(database, userAuth.Email)

	if err != nil {
		c.Status(http.StatusBadRequest)

		return c.JSON(fiber.Map{
			"error": err.Error()})
	}

	expiryTime := time.Now().Add(time.Hour * 24)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(userAuth.ID)),
		ExpiresAt: jwt.NewNumericDate(expiryTime),
	})

	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "could not login",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expiryTime,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	cookie2 := fiber.Cookie{
		Name:    "session",
		Value:   "active",
		Expires: expiryTime,
	}

	c.Cookie(&cookie2)
	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"messsage": "success",
	})

}

func GetUserData(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	claims := token.Claims

	id, err := claims.GetIssuer()

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	database, _ := db.Connect()
	var userAuth models.UserAuth
	err = database.Model(&models.UserAuth{}).Where("id = ?", id).First(&userAuth).Error

	if err == nil {
		var user models.User
		database.Model(&user).Where("id = ?", userAuth.UserID).First(&user)
		return c.JSON(user)
	}
	// fmt.Println(err.Error())
	return nil
}

func HandleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "deleted",
		Expires:  time.Now().Add(-(time.Hour * 2)), // Add negative time means it happens in the past :P
		HTTPOnly: true,
	})

	c.ClearCookie("session")
	fmt.Println("clearing cookies")
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func CreateOTPRequest(c *fiber.Ctx) error {
	var userAuth models.UserAuth

	if err := c.BodyParser(&userAuth); err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "could not login",
		})
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	if len(userAuth.Email) <= 0 || !models.DoesEmailExist(database, userAuth.Email) {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "invalid email"})
	}

	otp, err := CreateOTP(database, &userAuth)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	err = SendEmail(&SMTP_SERVER, Email, []string{userAuth.Email}, EmailPassword, "Your TraveloHI Verification Code", "Your verification code is <b>"+otp.Code+"</b>")

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf(err.Error())
	}

	c.Status(http.StatusOK)

	return c.JSON(otp)
}

func ValidateOTPRequest(c *fiber.Ctx) error {
	var otp models.OTP

	if err := c.BodyParser(&otp); err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("error binding json")
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed connecting to database")
	}

	err = ValidateOTP(database, &otp)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf(err.Error())
	}

	c.Status(http.StatusOK)
	return nil
}

func HandleGetSecurityQuestion(c *fiber.Ctx) error {
	var userAuth models.UserAuth

	if err := c.BodyParser(&userAuth); err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("error binding json")
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed connecting to database")
	}

	result, err := GetSecurityQuestion(database, &userAuth)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf(err.Error())
	}

	c.JSON(fiber.Map{"securityQuestion": result})
	return nil
}

func ValidateSecurityAnswer(c *fiber.Ctx) error {
	var userAuth models.UserAuth

	if err := c.BodyParser(&userAuth); err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("error binding json")
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed connecting to database")
	}

	result, err := GetSecurityAnswer(database, &userAuth)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf(err.Error())
	}

	if result != userAuth.SecurityAnswer {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf("wrong answer")
	}

	return nil
}

func HandleChangePassword(c *fiber.Ctx) error {
	var userAuth models.UserAuth

	if err := c.BodyParser(&userAuth); err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Invalid JSON"})
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	err = ChangePassword(database, &userAuth)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	return nil
}

func HandleGetAirportsData(c *fiber.Ctx) error {

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	airports, err := GetAirports(database)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	return c.JSON(fiber.Map{
		"airports": airports,
	})
}

func HandleGetAirlinesData(c *fiber.Ctx) error {

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	airports, err := GetAirlines(database)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	return c.JSON(fiber.Map{
		"airlines": airports,
	})
}

func HandleGetAirplanesData(c *fiber.Ctx) error {

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	airports, err := GetAirplanes(database)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	return c.JSON(fiber.Map{
		"airplanes": airports,
	})
}

func HandleAddFlightSchedule(c *fiber.Ctx) error {
	var flightSchedule models.FlightSchedule

	if err := c.BodyParser(&flightSchedule); err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return fmt.Errorf("error binding json")
	}
	// flightSchedule := string(c.BodyRaw())

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	isAvailable, err := IsAirplaneAvailable(database, flightSchedule.Airplane.ID, flightSchedule.DepartureTime, flightSchedule.ArrivalTime)

	if !isAvailable {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf("airplane is not available within the time period")
	}

	if err = AddFlightSchedule(database, &flightSchedule); err != nil {
		c.Status(http.StatusBadRequest)
		return err
	}

	return nil
}

func HandleGetAvailableAirplanes(c *fiber.Ctx) error {
	var data struct {
		Airline       models.Airline `json:"airline"`
		DepartureTime time.Time      `json:"departureTime"`
		ArrivalTime   time.Time      `json:"arrivalTime"`
	}

	if err := c.BodyParser(&data); err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("error binding json")
	}

	fmt.Println(data)

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	airplanes, err := GetAirplanes(database)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed retrieving airplanes data")
	}

	var ret []models.Airplane

	for _, airplane := range airplanes {
		if avail, _ := IsAirplaneAvailable(database, airplane.ID, data.DepartureTime, data.ArrivalTime); avail && airplane.Airline.ID == data.Airline.ID {
			ret = append(ret, *airplane)
		}
	}

	return c.JSON(fiber.Map{
		"airplanes": ret,
	})
}

func HandleAPI(c *fiber.Ctx) error {
	var object string

	if err := c.BodyParser(&object); err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("error binding json")
	}

	_, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	return nil
}
