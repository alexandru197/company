package adapters

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexandru197/company/config"
	repository "github.com/alexandru197/company/repository/company"
	services "github.com/alexandru197/company/services/company"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
)

var jwtSecret string

type CompanyApp struct {
	CompanyService services.CompanyService
	Config         config.Config
}

func NewCompanyApp() *CompanyApp {
	var conf config.Config = config.InitConfig()

	jwtSecret = conf.JWT.Secret

	db, err := config.InitDB(conf)
	if err != nil {
		log.Panicf("Error initializing the Database")
	}

	var kafkaWriter *kafka.Writer = config.InitKafka(conf)

	companyRepository := repository.NewCompanyRepository(db)
	companyService := services.NewCompanyService(companyRepository, kafkaWriter)

	return &CompanyApp{
		CompanyService: companyService,
		Config:         conf,
	}
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: invalid token format", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token method conforms to HMAC.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *CompanyApp) Start() {
	router := app.SetupRoutes()

	// Start the server.
	port := strconv.Itoa(app.Config.Server.Port)
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
