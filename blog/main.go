package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	tpb "grpc-blog/gunk/v1/categories"
	ppb "grpc-blog/gunk/v1/post"

	cc "grpc-blog/blog/core/categories"
	"grpc-blog/blog/services/categories"
	pc "grpc-blog/blog/core/post"
	"grpc-blog/blog/services/post"
	"grpc-blog/blog/storage/postgres"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {


	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("blog/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	grpcServer := grpc.NewServer()

	store, err := newDBFromConfig(config)
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err)
	}
	ccs := cc.NewCategorySvc(store)
	s := categories.NewCategoryServer(ccs)
	
	tpb.RegisterCategoryServiceServer(grpcServer, s)

	pcs := pc.NewPostSvc(store)
	p := post.NewPostServer(pcs)
	
	ppb.RegisterPostServiceServer(grpcServer, p)


	host, port := config.GetString("server.host"), config.GetString("server.port")

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		log.Fatalf("Failed to Listen: %s", err)
	}

	log.Printf("Server is starting on: http://%s:%s\n", host, port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %s", err)
	}
}

func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}