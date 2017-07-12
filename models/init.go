package models

import (
	"errors"
	"log"
	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
	"strconv"
	"time"
)

var (
	Db            *nodb.DB
	ErrNotExist   = errors.New("not exist")
	ErrParamError = errors.New("param error")
)

func Init() error {
	cfg := config.NewConfigDefault()
	cfg.DataDir = "./metas"
	

	var err error
	// init nosql
	db, err := nodb.Open(cfg)
	if err != nil {
		return err
	}

	// select db
	Db, err = db.Select(0)
	if err != nil {
		return err
	}

	var confp string = "./usercfg.conf"
	conf_fh, err := get_config(confp)
        if err != nil {
		log.Printf("read usercfg config file error")
		return errors.New("read usercfg config file error")
	}

	sections := get_sections(conf_fh)
	var passwordStr string

	// add admin
	for _, v := range sections {
		if v == "default" {
			continue
		}
		userItem, err := get_users(conf_fh, v)
		if err != nil {
			log.Printf("parse user for section %s error", v, err)
		}
		otpT, err := getOtpPass(userItem.secret)
		if err != nil {
			log.Printf("get totp for section %s error", v, err)
			otpT = 100000
		}
		Tremain, err := getOtpPassremain(userItem.secret)
		log.Printf("totp: %d remain seconds\n", Tremain)

		passwordStr = userItem.pass + strconv.FormatUint(uint64(otpT), 10)
		log.Printf("pass for user %s: %s\n",userItem.name, passwordStr)
		userFind, err := GetUserByName(userItem.name)
		if err == ErrNotExist {
			err = AddUser(&User{
				Name:     userItem.name,
				Password: passwordStr,
				Database: userItem.dbs,
			})
		} else {
			log.Printf("update pass for user %s: %s\n",userItem.name, passwordStr)
			err = UpdateUser(&User{
				Id:       userFind.Id,
				Name:     userItem.name,
				Password: passwordStr,
				Database: userItem.dbs,
			})
		}
		go loopUpdate(userItem)

	}
	return err
}

func loopUpdate(userItem *usermsg) {
	for {
		userFind, err := GetUserByName(userItem.name)
		if err != nil {
			log.Printf("cann't find Id for user %s ...\n", userItem.name)
			continue
		}
		Tremain, err := getOtpPassremain(userItem.secret)
		time.Sleep(time.Duration(Tremain) * time.Second)
		totp, err := getOtpPass(userItem.secret)
		passwordStr := userItem.pass + strconv.FormatUint(uint64(totp), 10)
		err = UpdateUser(&User{
			Id:       userFind.Id,
			Name:     userItem.name,
			Password: passwordStr,
			Database: userItem.dbs,
		})
		if err == nil {
			log.Printf("update %s for password %s ok\n", userItem.name, passwordStr)
		}
	}
}
