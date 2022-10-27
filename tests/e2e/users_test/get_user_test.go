package tests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/HackIllinois/api/common/configloader"
	user_models "github.com/HackIllinois/api/services/user/models"
	"github.com/HackIllinois/api/tests/common"
	"github.com/dghubble/sling"
	"go.mongodb.org/mongo-driver/mongo"
)

var admin_client *sling.Sling
var client *mongo.Client

func TestMain(m *testing.M) {
	cfg, err := configloader.Load(os.Getenv("HI_CONFIG"))

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	admin_client = common.GetSlingClient("Admin")
	client = common.GetLocalMongoSession()

	user_db_name, err := cfg.Get("USER_DB_NAME")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	client.Database(user_db_name).Drop(context.Background())
	return_code := m.Run()
	os.Exit(return_code)
}

func TestGetCurrentUserInfo(t *testing.T) {
	firstName := [5]string{"Bryant", "Tim", "Jareth", "Ananya", "Ashay"}
	lastName := [5]string{"Zhou", "Gonzalez", "Gomes", "Sehgal", "blah"}
	received_user := user_models.UserInfo{}
	for i, fName := range firstName {
		user := user_models.UserInfo{
			ID:        "github000001" + strconv.Itoa(i),
			Username:  fName,
			FirstName: fName,
			LastName:  lastName[i],
			Email:     fName + "@gmail.com",
		}
		_, err := admin_client.New().Post("/user/").BodyJSON(user).ReceiveSuccess((&received_user))
		if err != nil {
			t.Errorf("Unable to make req")
		}
	}
}
