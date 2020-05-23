package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database/mongodb"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/generated"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/resolver"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/middlewares"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/services/jwt"
	goKitLog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := mux.NewRouter()

	router.Use(middlewares.LoggingMiddleware(goKitLog.NewLogfmtLogger(os.Stdout)))
	router.Use(middlewares.EnrichUserID(initializeJWTService()))

	router.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", initializeGraphQLServer())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initializeGraphQLServer() *handler.Server {
	return handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: initializeResolver(),
			},
		),
	)
}
func initializeResolver() *resolver.Resolver {
	return resolver.NewResolver(initializeDB())
}
func initializeJWTService() jwt.Service {
	secret := os.Getenv("JWT_SECRET")
	return jwt.NewService(secret)
}
func initializeDB() database.DB {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot connect to mongoDB"))
	}

	db := client.Database(os.Getenv("MONGODB_DB_NAME"))

	return mongodb.NewMongoDB(db)
}
