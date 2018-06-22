package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RobinBaeckman/ragnar"
	"github.com/RobinBaeckman/ragnar/controller"
	"github.com/RobinBaeckman/ragnar/db/mysql"
	"github.com/RobinBaeckman/ragnar/errors"
	"github.com/RobinBaeckman/ragnar/middleware"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	parseConfig()
	re := newRedis()
	l := log.New(os.Stdout, viper.GetString("app.log_prefix"), 3)
	db := mysql.NewDB()
	defer db.Close()

	r := mux.NewRouter()

	s := &mysql.UserService{db}
	c := ragnar.NewUserCache(s, re)

	// User user
	r.Handle("/v1/users", Adapt(
		errors.Check(controller.UserStore(c)),
		middleware.Log(l),
	)).Methods("POST")
	r.Handle("/v1/users/{id}", Adapt(
		errors.Check(controller.UserShow(c)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("GET")
	r.Handle("/v1/users", Adapt(
		errors.Check(controller.UserIndex(s)),
		middleware.Log(l),
	)).Methods("GET")
	r.Handle("/v1/login", Adapt(
		errors.Check(controller.Login(c)),
		middleware.Log(l),
	)).Methods("POST")

	http.Handle("/", r)

	l.Println("Started")
	l.Fatal(http.ListenAndServe(viper.GetString("app.host")+":"+viper.GetString("app.port"), nil))
}

func Adapt(h http.Handler, adapters ...middleware.Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func parseConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		fmt.Errorf("Fatal error config file: %s \n", err)
	}
}

func newRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
