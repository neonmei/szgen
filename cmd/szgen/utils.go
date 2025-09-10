package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/neonmei/szgen/internal/config"
	"github.com/spf13/cobra"
)

func buildMetricConfig(cmd *cobra.Command, metricType string) (*config.MetricTask, error) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	unit, _ := cmd.Flags().GetString("unit")
	generator, _ := cmd.Flags().GetString("generator")
	value, _ := cmd.Flags().GetString("value")
	valueType, _ := cmd.Flags().GetString("type")
	count, _ := cmd.Flags().GetInt("count")
	rate, _ := cmd.Flags().GetDuration("rate")
	attributes, _ := cmd.Flags().GetStringToString("attributes")

	options := []config.MetricTaskOption{
		config.WithKind(metricType),
		config.WithName(name),
		config.WithDescription(description),
		config.WithUnit(unit),
		config.WithGenerator(generator),
		config.WithValue(value),
		config.WithType(valueType),
		config.WithCount(count),
		config.WithRate(rate),
		config.WithMetricAttributes(attributes),
	}

	mc := config.NewMetricTask(options...)

	if err := mc.Validate(); err != nil {
		return nil, err
	}

	slog.Debug("loaded config",
		"name", mc.Name,
		"kind", mc.Kind,
		"type", mc.Type,
		"rate", mc.Rate,
		"count", mc.Count,
		"value", mc.Value,
		"attributes", mc.Attributes,
		"generator", mc.Generator,
		"description", mc.Description,
		"unit", mc.Unit,
	)
	return mc, nil
}

func setupContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigCh:
			slog.Info("Received signal, shutting down gracefully", "signal", sig)
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(sigCh)
		close(sigCh)
	}()

	return ctx, cancel
}
