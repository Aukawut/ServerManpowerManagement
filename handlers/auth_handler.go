package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Aukawut/ServerManpowerManagement/config"
	"github.com/Aukawut/ServerManpowerManagement/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nmcclain/ldap"
)

func VerifyToken(tokenString string) (jwt.MapClaims, error) {

	// Retrieve the secret key from environment variables
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return nil, fmt.Errorf("secret_key not set in .env file")
	}

	// Convert the JWT secret key to a byte slice
	secretKey := []byte(jwtSecret)

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	// Check if the token is valid
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func AuthenticateUserDomain(username string, password string) (bool, string) {

	// Get LDAP server from environment variable
	ldapServer := os.Getenv("LDAP_IP")

	// Make sure LDAP_SERVER is set
	if ldapServer == "" {
		return false, fmt.Sprintf("LDAP_SERVER not defined in .env file,%v", "")
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, 389))
	if err != nil {

		return false, fmt.Sprintf("failed to connect to LDAP server: %v", err)
	}
	defer l.Close() // Close the connection once done

	err = l.Bind(username, password)
	if err != nil {

		return false, fmt.Sprintf("failed to authenticate user: %v", username)
	}

	// user is authenticated
	return true, ""
}

func GenerateToken(user model.UsersInfoLogin) (string, error) {

	// Retrieve the secret key from the environment variables
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return "", fmt.Errorf("secret_key not set in .env file") // Lowercase error message
	}

	// Convert the JWT secret key to a byte slice
	secretKey := []byte(jwtSecret)

	claims := jwt.MapClaims{

		"employee_code": user.UHR_EmpCode,
		"role":          user.RoleName,
		"department":    user.UHR_Department,
		"position":      user.UHR_Position,
		"exp":           time.Now().Add(time.Hour * 3).Unix(), // Token expiration (3 Hours)
	}

	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err) // Lowercase error message
	}

	return signedToken, nil
}

func LoginDomain(c *fiber.Ctx) error {

	req := model.BodyLoginDomain{}
	var usersDetail []model.UsersInfoLogin

	var err error

	if err = c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": err.Error(),
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Username and Password is required!",
		})
	}

	// Check User on Database

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	resultUser, errQuery := db.Query(`SELECT  [UHR_EmpCode],UHR_Position,UHR_Department,
CONCAT(UHR_FirstName_en,' ',UHR_LastName_en) as [UHR_FullName],UHR_Sex,r.ROLE_NAME as RoleName
FROM [DB_MANPOWER_MGT].[dbo].[V_AllUserPSTH] h
LEFT JOIN TBL_USERS u ON h.UHR_EmpCode COLLATE Thai_CI_AS = u.EMPLOYEE_CODE COLLATE Thai_CI_AS
LEFT JOIN TBL_ROLES r ON u.ROLE_ID = r.ROLE_ID
WHERE AD_UserLogon = @username AND UHR_StatusToUse = 'ENABLE' AND u.EMPLOYEE_CODE IS NOT NULL`, sql.Named("username", req.Username))

	if errQuery != nil {
		fmt.Println("Query failed: " + errQuery.Error())
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Query failed: " + errQuery.Error(),
		})

	}

	// Loop for append
	for resultUser.Next() {

		var userInfo model.UsersInfoLogin
		errScan := resultUser.Scan(

			&userInfo.UHR_EmpCode,
			&userInfo.UHR_Position,
			&userInfo.UHR_Department,
			&userInfo.UHR_FullName,
			&userInfo.UHR_Sex,
			&userInfo.RoleName,
		)

		if errScan != nil {
			fmt.Println("Scan error ", errScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errScan.Error(),
			})

		} else {
			usersDetail = append(usersDetail, userInfo)
		}

	}

	// Check Childen in array struct
	if len(usersDetail) > 0 {

		// Ldap Authen
		usernameAd := req.Username + "@" + os.Getenv("LDAP_DNS") // awk@psth.com

		verifiesDomain, errorMsg := AuthenticateUserDomain(usernameAd, req.Password)

		if verifiesDomain && errorMsg == "" {
			token, _ := GenerateToken(usersDetail[0])

			// Success
			return c.JSON(fiber.Map{
				"err":     false,
				"msg":     "Success",
				"status":  "Ok",
				"results": usersDetail,
				"token":   token,
			})

		} else {
			// Login Failed
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Sorry, Login is failed!!",
			})
		}

	} else {
		return c.JSON(fiber.Map{
			"err":  true,
			"msg":  "You don't have permission!",
			"user": usersDetail,
		})
	}

}

func CheckToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Authorization header must start with 'Bearer '",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, errToken := VerifyToken(token) // Ensure VerifyToken is implemented correctly

	if errToken != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errToken.Error(),
		})
	} else {
		// Auth Success
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "Auth success",
			"decoded": decoded,
			"status":  "Ok",
		})
	}

}
