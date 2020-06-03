package api

import (
	"context"
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

	StreamSubject = "foo"
	StreamName = "foo-stream"
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
	e.POST("/kafka/data/:groupid", jellyfishPostData)
	e.POST("/data/:groupid", jellyfishPostData)
	e.POST("/v1/device/upload/cl", jellyfishPostCarelinkData)

	// platform data
	e.POST("/dataservices/v1/datasets/:dataSetId/data", platformDataSetsDataCreate)
	e.POST("/dataservices/v1/datasets/:dataSetId", platformDataSetsDelete)
	e.POST("/dataservices/v1/datasets/:dataSetId", platformDataSetsUpdate)

	e.POST("/dataservices/v1/users/:userId/data", platformUsersDataDelete)
	e.POST("/dataservices/v1/users/:userId/datasets", platformUsersDataSetsCreate)

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
func jellyfishPostData(c echo.Context) error {
	// to produce messages
	topic := "data"
	partition := 0
	host := "kafka-kafka-bootstrap"
	port := 9092
	hostStr := fmt.Sprintf("%s:%d", host,port)

	conn, err := kafka.DialLeader(context.Background(), "tcp", hostStr, topic, partition)
	if err != nil {

		fmt.Printf("Error making connection: %s", err.Error())
		return c.String(http.StatusOK, err.Error())
	}

	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	conn.WriteMessages(
		kafka.Message{Value: []byte("Hello jellyfish!")},
	)

	conn.Close()
	return c.String(http.StatusOK, "Hello, World!")
}

func jellyfishPostCarelinkData(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// platform data
func platformDataSetsDataCreate(c echo.Context) error {
	//dataSetId := c.QueryParam("dataSetId")
	return c.String(http.StatusOK, "Hello, World!")
}

func platformDataSetsDelete(c echo.Context) error {
	//dataSetId := c.QueryParam("dataSetId")
	return c.String(http.StatusOK, "Hello, World!")
}

func platformDataSetsUpdate(c echo.Context) error {
	//userId := c.QueryParam("userId")
	return c.String(http.StatusOK, "Hello, World!")
}

func platformUsersDataDelete(c echo.Context) error {
	//userId := c.QueryParam("userId")
	return c.String(http.StatusOK, "Hello, World!")
}

func platformUsersDataSetsCreate(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

