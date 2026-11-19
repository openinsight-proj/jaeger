// Copyright (c) 2024 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
)

func TestComponents(t *testing.T) {
	factories, err := Components()

	require.NoError(t, err)

	assert.NotNil(t, factories.Extensions)
	assert.NotNil(t, factories.Receivers)
	assert.NotNil(t, factories.Exporters)
	assert.NotNil(t, factories.Processors)
	assert.NotNil(t, factories.Connectors)

	_, jaegerReceiverFactoryExists := factories.Receivers[component.MustNewType("jaeger")]
	assert.True(t, jaegerReceiverFactoryExists)
	assert.Contains(t, factories.Receivers, component.MustNewType("skywalking"))
	assert.Contains(t, factories.Receivers, component.MustNewType("fluentforward"))
	assert.Contains(t, factories.Exporters, component.MustNewType("elasticsearch"))
	assert.Contains(t, factories.Exporters, component.MustNewType("loadbalancing"))
	assert.Contains(t, factories.Exporters, component.MustNewType("prometheusremotewrite"))
	assert.Contains(t, factories.Exporters, component.MustNewType("sumologic"))
	assert.Contains(t, factories.Processors, component.MustNewType("transform"))
	assert.Contains(t, factories.Processors, component.MustNewType("probabilistic_sampler"))
	assert.Contains(t, factories.Connectors, component.MustNewType("servicegraph"))
}

func TestGetOtelcolFactories(t *testing.T) {
	errMock := errors.New("mock error")
	{
		b := defaultBuilders()
		b.extension = mockFactoryMap[extension.Factory](errMock)
		_, err := b.build()
		require.ErrorIs(t, err, errMock)
	}
	{
		b := defaultBuilders()
		b.receiver = mockFactoryMap[receiver.Factory](errMock)
		_, err := b.build()
		require.ErrorIs(t, err, errMock)
	}
	{
		b := defaultBuilders()
		b.exporter = mockFactoryMap[exporter.Factory](errMock)
		_, err := b.build()
		require.ErrorIs(t, err, errMock)
	}
	{
		b := defaultBuilders()
		b.processor = mockFactoryMap[processor.Factory](errMock)
		_, err := b.build()
		require.ErrorIs(t, err, errMock)
	}
	{
		b := defaultBuilders()
		b.connector = mockFactoryMap[connector.Factory](errMock)
		_, err := b.build()
		require.ErrorIs(t, err, errMock)
	}
}

func mockFactoryMap[FACTORY component.Factory](err error) func(factories ...FACTORY) (map[component.Type]FACTORY, error) {
	return func( /* factories */ ...FACTORY) (map[component.Type]FACTORY, error) {
		return nil, err
	}
}
