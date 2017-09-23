// Copyright 2017 Kubernetes Community Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"net/http"

	"github.com/zabelinva/slavonapp/pkg/config"
	"github.com/zabelinva/slavonapp/pkg/handlers"
	"github.com/zabelinva/slavonapp/pkg/logger"
	stdlog "github.com/zabelinva/slavonapp/pkg/logger/standard"
	"github.com/zabelinva/slavonapp/pkg/router"
	"github.com/zabelinva/slavonapp/pkg/router/bitroute"
	"github.com/zabelinva/slavonapp/pkg/version"
)

// Setup configures the service
func Setup(cfg *config.Config) (r router.BitRoute, log logger.Logger, err error) {
	// Setup logger
	log = stdlog.New(&logger.Config{
		Level: cfg.LogLevel,
		Time:  true,
		UTC:   true,
	})

	log.Info("Version:", version.RELEASE)
	log.Warnf("%s log level is used", logger.LevelDebug.String())
	log.Infof("Service %s listened on %s:%d", config.SERVICENAME, cfg.LocalHost, cfg.LocalPort)

	// Define handlers
	h := handlers.New(log, cfg)

	// Register new router
	r = bitroute.New()

	// Response for undefined methods
	r.SetupNotFoundHandler(h.Base(notFound))

	// Configure router
	r.SetupMiddleware(h.Base)
	r.GET("/", h.Root)
	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Ready)
	r.GET("/info", h.Info)

	return
}

// Response for undefined methods
func notFound(c router.Control) {
	c.Code(http.StatusNotFound)
	c.Body("Method not found for " + c.Request().URL.Path)
}
