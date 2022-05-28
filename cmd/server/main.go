package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vanamelnik/gophkeeper/models"
	"github.com/vanamelnik/gophkeeper/server/gophkeeper"
	"github.com/vanamelnik/gophkeeper/server/storage/postgres"
	"github.com/vanamelnik/gophkeeper/server/users"
)

const (
	cfgFileName = "server_config.yaml"
)

func main() {
	ctx := context.Background()
	log.Println("Starting GophKeeper server")
	must(loadConfig())

	s, err := postgres.NewStorage(viper.GetString("databaseDSN"), postgres.WithDesctructiveReset())
	must(err)
	defer s.Close()

	u := users.NewService(
		s,
		viper.GetString("tokens.secretKey"),
		viper.GetDuration("tokens.accessTokenDuration"),
		viper.GetDuration("tokens.refreshTokenDuration"),
	)
	userID, err := u.CreateUser(ctx, "jmvmiller@gmail.com", "myPassword")
	must(err)
	access, refresh, err := u.CreateSession(ctx, userID)
	must(err)
	time.Sleep(time.Second * 2)
	_, err = u.Authenticate(ctx, access)
	fmt.Println(err)
	access, refresh, err = u.RefreshTheTokens(ctx, refresh)
	must(err)
	id, err := u.Authenticate(ctx, access)
	fmt.Println(id, err)

	g := gophkeeper.NewService(s)
	defer g.Close()

	now := time.Now()
	pwdID := uuid.New()
	cardID := uuid.New()
	g.PublishUserData(ctx, id, []models.Event{
		models.Event{
			Operation: models.OpCreate,
			Item: models.Item{
				ID:        pwdID,
				Version:   0,
				CreatedAt: &now,
				DeletedAt: nil,
				Payload: models.PasswordData{
					Password: "bebebe",
				},
				Meta: `{"login": "qweqweqwe"}`,
			},
		},
	})
	time.Sleep(time.Millisecond)
	g.PublishUserData(ctx, id, []models.Event{
		models.Event{
			Operation: models.OpUpdate,
			Item: models.Item{
				ID:        pwdID,
				Version:   1,
				CreatedAt: &now,
				DeletedAt: nil,
				Payload: models.PasswordData{
					Password: "bebebe1",
				},
				Meta: `{"login": "qweqweqw1"}`,
			},
		},
		models.Event{
			Operation: models.OpCreate,
			Item: models.Item{
				ID:        cardID,
				Version:   0,
				CreatedAt: &now,
				DeletedAt: nil,
				Payload: models.CardData{
					Number:         "123456789",
					CardholderName: "IVAN TOPORYSHKIN",
					Date:           "01.01.2001",
					CVC:            765,
				},
				Meta: `{"bank": "sber"}`,
			},
		},
	})
	time.Sleep(100 * time.Millisecond)
	version, err := u.GetDataVersion(ctx, id)
	must(err)
	fmt.Printf("DataVersion is %d\n", version)

	data, err := g.GetUserData(ctx, id, map[uuid.UUID]uint64{
		pwdID: 1,
	})
	must(err)
	for _, item := range data.Items {
		fmt.Printf("%+v\n", item)
	}
	sID, err := u.GetSessionID(refresh)
	must(err)
	must(u.Logout(ctx, sID))
	_, _, err = u.RefreshTheTokens(ctx, refresh)
	fmt.Println(err)
}

func loadConfig() error {
	viper.SetConfigFile(cfgFileName)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	log.Println("Config loaded")
	return nil
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
