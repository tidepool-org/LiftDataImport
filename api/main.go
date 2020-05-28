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
	lift "github.com/liftbridge-io/go-liftbridge"
)

var (
	Host = "localhost"
	Port = 8080
	ServerString = fmt.Sprintf("%s:%d", Host, Port)
	ServerTimeoutAmount = 20

	StreamSubject = "foo"
	StreamName = "foo-stream"

)

type StreamContext struct {
	echo.Context
	Client lift.Client
}

func MainLoop() {
	// Echo instance
	e := echo.New()
	e.Logger.Print("Starting Main Loop")

	// Create Liftbridge client.
	addrs := []string{"localhost:9292", "localhost:9293", "localhost:9294"}
	client, err := lift.Connect(addrs)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Create a stream attached to the NATS subject "foo".
	if err := client.CreateStream(context.Background(), StreamSubject, StreamName); err != nil {
		if err != lift.ErrStreamExists {
			panic(err)
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &StreamContext{c, client}
			return next(cc)
		}
	})

	// Routes
	// jellyfish
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
	//groupId := c.QueryParam("groupid")
	cc := c.(*StreamContext)
	// Publish a message to "foo".
	if _, err := cc.Client.Publish(context.Background(), StreamName, []byte("hello")); err != nil {
		panic(err)
	}

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

