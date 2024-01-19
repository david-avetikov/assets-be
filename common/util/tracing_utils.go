package util

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
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
)

var b3Tracer = b3.HTTPFormat{}

const (
	traceIdKey = "TraceId"
	spanIdKey  = "SpanId"
)

func TraceAppender(context *gin.Context) {
	spanCtx, _ := b3Tracer.SpanContextFromRequest(context.Request)
	_, span := trace.StartSpanWithRemoteParent(context, "code", spanCtx)
	spanCtx = span.SpanContext()
	context.Set(traceIdKey, spanCtx.TraceID.String())
	context.Set(spanIdKey, spanCtx.SpanID.String())
}

func InitTracing() {
	otel.SetTracerProvider(jaegerTraceProvider())
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
func jaegerTraceProvider() *sdktrace.TracerProvider {
	endpoint := fmt.Sprintf(
		"http://%s:%s/api/traces",
		config.CoreConfig.Database.Jaeger.CollectorHost,
		config.CoreConfig.Database.Jaeger.CollectorPort,
	)
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(MustOne(jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))))),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ApplicationName),
			semconv.DeploymentEnvironmentKey.String(os.Getenv("PROFILE")),
		)),
	)
}
