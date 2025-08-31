package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	_ = godotenv.Load(".env")

	// Prefer full URI; otherwise build from parts safely
	mongoURI := os.Getenv("MONGODB_URL")
	if mongoURI == "" {
		user := os.Getenv("MONGODB_USER")
		pass := os.Getenv("MONGODB_PASSWORD")
		host := os.Getenv("MONGODB_HOST") // e.g. cluster0.xxx.mongodb.net
		if user != "" && pass != "" && host != "" {
			// URL-encode password to avoid auth failures on special characters
			encodedPass := url.QueryEscape(pass)
			authSource := os.Getenv("AUTH_SOURCE") // optional
			qs := "?retryWrites=true&w=majority"
			if authSource != "" {
				qs = qs + "&authSource=" + url.QueryEscape(authSource)
			}
			mongoURI = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", user, encodedPass, host, qs)
		}
	}

	if mongoURI == "" {
		log.Fatal("MONGODB_URL or MONGODB_USER/PASSWORD/HOST must be set")
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatal("failed to create mongo client: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal("mongo connect error: ", err)
	}
	// Fail fast if creds/network are wrong
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("mongo ping failed (check URI, user/password, IP whitelist, authSource): ", err)
	}

	fmt.Println("connected to MongoDB: ", mongoURI)
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "jwt_auth_db"
	}
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}
