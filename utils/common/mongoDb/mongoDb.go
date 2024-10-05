package mongoDb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

var client *mongo.Client

func Init(mongoAddr string) {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoAddr))

	//  mongodb://mongouser:O1*%5EJeWEP1yGqb!QqZrl@172.20.0.39:27017,172.20.0.31:27017,172.20.0.18:27017/test?authSource=admin
	if err != nil {
		log.Fatal("%v--%v", mongoAddr, err)
	}
	ctx := context.Background() //(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("%v--%v", mongoAddr, err)
	}
	//defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("%v--%v", mongoAddr, err)
	}
}
func GetClient() *mongo.Client {
	return client
}
