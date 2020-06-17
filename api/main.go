package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	kafka "github.com/segmentio/kafka-go"
)

var (
	Host = "localhost"
	Port = 8080
	ServerString = fmt.Sprintf("%s:%d", Host, Port)
	ServerTimeoutAmount = 20
)

func MainLoop() {
	// Echo instance
	e := echo.New()
	e.Logger.Print("Starting Main Loop")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// jellyfish
	e.POST("*", postData)
	e.PUT("*", postData)
	e.DELETE("*", postData)

	// Start server
	e.Logger.Printf("Starting Server at: %s\n", ServerString)
	go func() {
		if err := e.Start(ServerString); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ServerTimeoutAmount) * time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Handler
// jellyfish
func postData(c echo.Context) error {
	// to produce messages
	topic := "data"
	partition := 0
	host := "kafka-kafka-bootstrap.kafka.svc.cluster.local"
	port := 9092
	hostStr := fmt.Sprintf("%s:%d", host,port)

	req := echo.Map{}
	data := echo.Map{}
	body := echo.Map{}
	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	data["body"] = body
	data["headers"] = c.Request().Header
	data["method"] = c.Request().Method
	data["url"] = c.Request().URL
	data["host"] = c.Request().Host

	req["data"] = data
	req["source"] = "api"

	reqString, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Problem converting to json: %s", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Connect to kafka
	conn, err := kafka.DialLeader(context.Background(), "tcp", hostStr, topic, partition)
	if err != nil {

		fmt.Printf("Error making connection: %s", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Write out html
	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	conn.WriteMessages(
		kafka.Message{Value: []byte(reqString)},
	)

	conn.Close()
	return c.NoContent(http.StatusOK)
}

