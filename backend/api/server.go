package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/validator"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/validator/govalidator"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/services/password"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/errorhandler"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/logger"
	"github.com/getsentry/sentry-go"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/cache"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/cache/redis"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database/mongodb"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/generated"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/resolver"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/middlewares"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/services/jwt"
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

	initializeLogger()

	router := mux.NewRouter()

	router.Use(middlewares.LoggingMiddleware(initializeLogger()))
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
	return resolver.NewResolver(
		initializeDB(),
		initializeValidator(),
		initializePasswordService(),
		initializeErrorHandler(),
		initializeLogger(),
		initializeJWTService(),
	)
}

func initializeValidator() validator.Validator {
	return govalidator.New(initializeDB())
}

func initializePasswordService() password.Service {
	return password.NewBcryptService()
}

func initializeJWTService() jwt.Service {
	secret := os.Getenv("JWT_SECRET")
	return jwt.NewService(secret, initializeCache())
}

func initializeDB() database.DB {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot connect to mongoDB"))
	}

	db := client.Database(os.Getenv("MONGODB_DB_NAME"))

	return mongodb.NewMongoDB(db)
}

func initializeCache() cache.Cache {
	return redis.NewClient(redis.Options{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func initializeLogger() logger.Logger {
	return logger.NewGoKitLogger(os.Stdout)
}

func initializeErrorHandler() errorhandler.ErrorHandler {
	errHandler, err := errorhandler.NewSentryErrorHandler(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: os.Getenv("SENTRY_DSN"),
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return errHandler
}
