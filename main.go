package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	DBUsername    string `mapstructure:"DB_USERNAME"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBAddress     string `mapstructure:"DB_ADDRESS"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBName        string `mapstructure:"DB_NAME"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.AutomaticEnv()

	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_NAME", "test")

	if !viper.IsSet("DB_ADDRESS") {
		return nil, fmt.Errorf("DB_ADDRESS must be set")
	}

	if !viper.IsSet("DB_USERNAME") {
		return nil, fmt.Errorf("DB_USERNAME must be set")
	}

	if !viper.IsSet("DB_PASSWORD") {
		return nil, fmt.Errorf("DB_PASSWORD must be set")
	}

	viper.Set("DB_ADDRESS", viper.GetString("DB_ADDRESS"))
	viper.Set("DB_USERNAME", viper.GetString("DB_USERNAME"))
	viper.Set("DB_PASSWORD", viper.GetString("DB_PASSWORD"))

	err = viper.Unmarshal(&config)
	return
}

type DatabaseCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// getCreds responds with credentials for the database
func (env *Env) getCreds(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, env.creds)
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Env struct {
	db    *sql.DB
	creds DatabaseCreds
}

func NewDB(dataSource string, creds DatabaseCreds) (*Env, error) {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Env{db: db, creds: creds}, nil
}

// getProducts responds with the list of all products as JSON.
func (env *Env) getProducts(c *gin.Context) {

	rows, err := env.db.Query("select id,name,description from products")

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	defer rows.Close()
	products := make([]Product, 0)
	for rows.Next() {
		var id string
		var name string
		var description string
		err = rows.Scan(&id, &name, &description)
		if err != nil {
			log.Printf("error occurred while reading the database rows: %v", err)
			break
		}
		products = append(products, Product{
			ID:          id,
			Name:        name,
			Description: description,
		})
	}

	c.IndentedJSON(http.StatusOK, products)
}

func health(c *gin.Context) {
	c.Status(http.StatusOK)
}

func main() {

	configFilePath := viper.GetString("CONFIG_FILE_PATH")
	if configFilePath == "" {
		configFilePath = "."
	}

	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	datasource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUsername, config.DBPassword, config.DBAddress, config.DBPort, config.DBName)

	hclog.Default().Info(datasource)

	env, err := NewDB(datasource, DatabaseCreds{Username: config.DBUsername, Password: config.DBPassword})
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	router := gin.Default()
	router.GET("/", health)
	router.GET("/products", env.getProducts)
	router.GET("/creds", env.getCreds)

	router.Run(config.ServerAddress)
}
