package server

/*
 * Copyright © 2024, "DEADLINE TEAM" LLC
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are not permitted.
 *
 * THIS SOFTWARE IS PROVIDED BY "DEADLINE TEAM" LLC "AS IS" AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL "DEADLINE TEAM" LLC BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * No reproductions or distributions of this code is permitted without
 * written permission from "DEADLINE TEAM" LLC.
 * Do not reverse engineer or modify this code.
 *
 * © "DEADLINE TEAM" LLC, All rights reserved.
 */

import (
	"assets/common/config"
	"assets/common/controller"
	"assets/common/iface"
	"assets/common/middleware"
	"assets/common/util"
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"strconv"
	"time"
)

var startTime = time.Now()

type httpServer struct {
	httpServer  *http.Server
	startupFunc func(*gin.Engine)
}

func NewHttpServer(startupFunc func(*gin.Engine)) iface.Runnable {
	server := &httpServer{startupFunc: startupFunc}
	closer.Bind(func() {
		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()
		log.WithError(server.httpServer.Shutdown(ctx)).Info("Http server stopped")
	})
	return server
}

func (server *httpServer) Run() {
	if log.GetLevel() < log.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery)
	router.Use(util.AddCorsHeaders)
	router.Use(util.TraceAppender)
	router.Use(middleware.HttpRequestLogger)
	router.Use(otelgin.Middleware(util.ApplicationName))

	controller.Register(router, controller.NewHealthCheckController())
	controller.Register(router, controller.NewSwaggerController())
	controller.Register(router, controller.NewPprofController())
	controller.Register(router, controller.NewOptionsController())
	controller.Register(router, controller.NewMetricsController())

	server.startupFunc(router)

	server.httpServer = &http.Server{
		Addr:           ":" + strconv.Itoa(config.CoreConfig.Server.Port),
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if config.CoreConfig.Server.SslKeyPath != "" && config.CoreConfig.Server.SslCrtPath != "" {
		log.Infof("Https server started at %dms", time.Since(startTime).Milliseconds())
		util.Must(server.httpServer.ListenAndServeTLS(config.CoreConfig.Server.SslCrtPath, config.CoreConfig.Server.SslKeyPath))
	} else {
		log.Infof("Http server started at %dms", time.Since(startTime).Milliseconds())
		util.Must(server.httpServer.ListenAndServe())
	}
}
