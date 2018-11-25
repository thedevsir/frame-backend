package main

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/spf13/cobra"
	"github.com/thedevsir/frame-backend/app/model"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/config/database"
	"github.com/thedevsir/frame-backend/services/encrypt"
	"github.com/thedevsir/frame-backend/services/utils"
)

func main() {

	config.Composer(".env")

	var adminPassword string
	var cmdAdminUserInstall = &cobra.Command{
		Use:   "create-admin",
		Short: "Create a admin with username and password",
		Run: func(cmd *cobra.Command, args []string) {

			db := &database.Composer{
				Locals:   "resource/locals/locals.json",
				Addrs:    strings.Split(config.DBAddress, ","),
				Database: config.DBName,
				Username: config.DBUsername,
				Password: config.DBPassword,
				Source:   config.DBSource,
			}
			db.Shoot()

			adminModel := database.Connection.Model(model.AdminCollection)
			admin := &model.Admin{}
			adminModel.New(admin)
			hash, _ := encrypt.Hash(adminPassword)

			admin.Id = bson.ObjectIdHex("000000000000000000000000")
			admin.Username = "root"
			admin.Password = hash
			admin.IsActive = true
			admin.CreatedAt = time.Now()
			admin.UpdatedAt = time.Now()

			err := admin.Save()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Root admin successfully created!")
		},
	}
	cmdAdminUserInstall.Flags().StringVarP(&adminPassword, "password", "p", "admin", "Password of root admin")

	var rootCmd = &cobra.Command{Use: "cmd"}
	rootCmd.AddCommand(cmdAdminUserInstall)
	rootCmd.Execute()
}
