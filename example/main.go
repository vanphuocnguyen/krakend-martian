package main

import (
	"context"
	"flag"
	"log"

	"github.com/devopsfaith/krakend-gologging"
	"github.com/devopsfaith/krakend-viper"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/proxy"
	krakendgin "github.com/luraproject/lura/router/gin"
	"github.com/luraproject/lura/transport/http/client"

	"github.com/muoingon/krakend-martian"
)

func main() {
	port := flag.Int("p", 0, "Port of the service")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "/etc/krakend/configuration.json", "Path to the configuration filename")
	flag.Parse()

	parser := viper.New()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}

	logger, err := gologging.NewLogger(serviceConfig.ExtraConfig)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	logger.Debug("config:", serviceConfig)

	ctx, cancel := context.WithCancel(context.Background())

	backendFactory := martian.NewBackendFactory(logger, client.DefaultHTTPRequestExecutor(client.NewHTTPClient))

	routerFactory := krakendgin.NewFactory(krakendgin.Config{
		Engine:         gin.Default(),
		Logger:         logger,
		HandlerFactory: krakendgin.EndpointHandler,
		ProxyFactory:   proxy.NewDefaultFactory(backendFactory, logger),
	})

	routerFactory.NewWithContext(ctx).Run(serviceConfig)

	cancel()
}
