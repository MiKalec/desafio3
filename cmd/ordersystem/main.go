package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/MiKalec/desafio3/configs"
	"github.com/MiKalec/desafio3/internal/event/handler"
	"github.com/MiKalec/desafio3/internal/infra/database"
	"github.com/MiKalec/desafio3/internal/infra/graph"
	"github.com/MiKalec/desafio3/internal/infra/grpc/pb"
	"github.com/MiKalec/desafio3/internal/infra/grpc/service"
	"github.com/MiKalec/desafio3/internal/infra/web/webserver"
	"github.com/MiKalec/desafio3/internal/usecase"
	"github.com/MiKalec/desafio3/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	if err := configs.EnsureDatabaseAndTable(cfg); err != nil {
		panic(err)
	}

	db, err := sql.Open(cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(cfg)
	if rabbitMQChannel == nil {
		log.Println("RabbitMQ indisponível: eventos OrderCreated não serão publicados. Suba com docker compose up para habilitar.")
	}

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	orderRepository := database.NewOrderRepository(db)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			webOrderHandler.Create(w, r)
		case http.MethodGet:
			webOrderHandler.List(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", cfg.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", cfg.GraphQLServerPort)
	http.ListenAndServe(":"+cfg.GraphQLServerPort, nil)
}

func getRabbitMQChannel(cfg *configs.Conf) *amqp.Channel {
	url := cfg.RabbitMQURL
	if url == "" {
		url = os.Getenv("RABBITMQ_URL")
	}
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("RabbitMQ: conexão falhou (%v)", err)
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("RabbitMQ: canal falhou (%v)", err)
		conn.Close()
		return nil
	}
	return ch
}
