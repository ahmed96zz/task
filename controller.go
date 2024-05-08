package main

import (
	"errors"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"task/task"
	"time"

	"github.com/gin-gonic/gin"
)

type ProfileReq struct {
	FirstName string `json:"name"`
	LastName  string `json:"phone_number"`
}

type OTPReq struct {
	OTP   string `json:"otp"`
	Phone string `json:"phone_number"`
}

// creatUser ...
func creatUser(c *gin.Context) {
	var user ProfileReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	queries := task.New(db)
	if !IsValidUserInputs(user) {
		c.JSON(400, gin.H{"error": "invalid inputs"})
		return
	}
	//check if phone number exist
	exsitingUser, err := queries.PhoneExisted(ctx, user.LastName)
	if err != nil {
		log.Println("Error in checking existing phone", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if exsitingUser {
		log.Println("phone already exist!")
		err = errors.New("phone already exist")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = UserTransaction(user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "User created successfully!",
	})
}

func UserTransaction(profile ProfileReq) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("Error in UserTransAction 1", err.Error())
		return err
	}
	defer tx.Rollback(ctx)
	queries := task.New(db)
	qtx := queries.WithTx(tx)

	userid, err := qtx.CreateProfile(ctx, task.CreateProfileParams{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	})
	if err != nil {
		log.Println("Error in UserTransAction 2", err.Error())
		return err
	}
	err = qtx.CreateUser(ctx, task.CreateUserParams{
		PhoneNumber: profile.LastName,
		ID:          userid,
	})
	if err != nil {
		log.Println("Error in UserTransAction 3", err.Error())
		return err
	}
	return tx.Commit(ctx)
}

func GetUsers(c *gin.Context) {

	queries := task.New(db)
	usersCh := make(chan []task.Profile)
	errCh := make(chan error)

	go func() {
		users, err := queries.SelectProfiles(ctx)
		if err != nil {
			log.Println("error getting Users", err.Error())
			errCh <- err
			return
		}
		usersCh <- users
	}()
	usersList := <-usersCh
	go func() {
		sort.SliceStable(usersList, func(i, j int) bool {
			return usersList[i].FirstName < usersList[j].FirstName
		})
		errCh <- nil
	}()
	err := <-errCh
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{
			"users": usersList,
		})
	}

}

func createOTP(c *gin.Context) {
	var user ProfileReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	queries := task.New(db)
	exsitingUser, err := queries.PhoneExisted(ctx, user.LastName)
	if err != nil {
		log.Println("Error in checking otp phone", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !exsitingUser {
		log.Println("phone doesn't exist")
		err = errors.New("phone doesn't exist")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	otp, err := generateOTP(user.LastName)
	if err != nil {
		log.Println("Error in createOTP 1", err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = func() error {
		tx, err := db.Begin(ctx)
		if err != nil {
			log.Println("Error in otp transaction 1", err.Error())
			return err
		}
		queries := task.New(db)
		qtx := queries.WithTx(tx)
		defer tx.Rollback(ctx)
		err = qtx.SetOTP(ctx, task.SetOTPParams{Otp: otp, PhoneNumber: user.LastName})
		if err != nil {
			log.Println("Error in createOTP 2", err.Error())
			return err
		}
		return tx.Commit(ctx)

	}()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{
			"message": "otp created successfully!", "otp": otp,
		})
	}

}

func generateOTP(number string) (otp string, err error) {
	buffer := make([]byte, 4)
	_, err = rand.Read(buffer)
	if err != nil {
		log.Println("error in generateOTP 1", err)
		return
	}
	numberInt, _ := strconv.ParseFloat(number, 64)
	base := strconv.FormatInt(time.Now().UnixNano()+int64(numberInt), 10)
	baseLength := len(base)
	for i := 0; i < 4; i++ {
		buffer[i] = base[int(buffer[i]%byte(baseLength))]
	}
	otp = string(buffer)
	return
}

func verifyOTP(c *gin.Context) {
	var otpReq OTPReq
	if err := c.ShouldBindJSON(&otpReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := func() error {
		tx, err := db.Begin(ctx)
		if err != nil {
			log.Println("Error in otp transaction 1", err.Error())
			return err
		}
		defer tx.Rollback(ctx)
		queries := task.New(db)
		qtx := queries.WithTx(tx)
		isValid, err := qtx.IsValidUserOTP(ctx, task.IsValidUserOTPParams{
			Otp:         otpReq.OTP,
			PhoneNumber: otpReq.Phone,
		})
		if err != nil {
			log.Println("error checking valid otp 1", err.Error())
			return err
		}
		if !isValid {
			return errors.New("otp doesn't exist for this phone")
		}
		isNotExpired, err := qtx.IsOTPExpired(ctx, task.IsOTPExpiredParams{
			Otp:         otpReq.OTP,
			PhoneNumber: otpReq.Phone,
		})
		if err != nil {
			log.Println("error checking valid otp 2", err.Error())
			return err
		}
		if !isNotExpired {
			return errors.New("otp has been expired")
		}
		return tx.Commit(ctx)
	}()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{
			"message": "opt is valid!",
		})
	}

}

func IsValidUserInputs(user ProfileReq) (valid bool) {
	if strings.TrimSpace(user.FirstName) == "" || strings.TrimSpace(user.LastName) == "" {
		return
	}
	_, err := strconv.Atoi(user.LastName)
	if err != nil {
		return
	}
	valid = true
	return
}
