package main

import (
	"context"
	"log"

	"github.com/bocdagla/tavern/adventurersguild/cmd/server/amqp"
	"github.com/bocdagla/tavern/adventurersguild/cmd/server/amqp/createQuest"
	"github.com/bocdagla/tavern/adventurersguild/cqrs/quest"
	mongoQuest "github.com/bocdagla/tavern/adventurersguild/pkg/domain/quest/mongo"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoDbConnectUri      = "<connection string>"
	questEN                = "quest"
	createQuestQN          = "createQuest"
	createQuestBindingKey  = "createQuest"
	createQuestConsumerTag = "createQuest"

	workerPoolSize = 2
)

func main() {

	//implement a better logger and add a interface for it
	logger := log.Default()

	//database conection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoDbConnectUri).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatalf("An error ocurred when trying to connect to the mongodb server: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("An error ocurred when trying to DISCONECT¿?¿? from the mongodb server: %v", err)
		}
	}()
	db := mongoClient.Database("adventurersguild")

	//repository creation
	questRepository := mongoQuest.New(db.Collection("quest"))

	//worker creation
	createQuestHandler := quest.NewHandler(logger, questRepository)
	createQuestWorker := createQuest.New(logger, createQuestHandler)

	//amqp connection
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("An error ocurred when trying to connect to the amqp server: %v", err)
	}
	defer conn.Close()
	c := amqp.New(conn, logger, createQuestWorker)
	ch, err := c.CreateChannel(questEN, createQuestQN, createQuestBindingKey, createQuestConsumerTag)
	if err != nil {
		log.Fatal(errors.Wrap(err, "CreateChannel"))
	}

	//consumer bootstrap
	if err = c.StartConsumer(workerPoolSize, createQuestQN, createQuestConsumerTag, ch); err != nil {
		log.Fatal(err)
	}
}
