package metricscmd

import (
	"log/slog"

	"github.com/neonmei/szgen/internal/config"
	"github.com/spf13/cobra"
)

// buildMetricTask builds a task of the given kind from the task defaults, overridden by explicitly set flags.
func buildMetricTask(cmd *cobra.Command, kind string) (*config.MetricTask, error) {
	options := []config.MetricTaskOption{
		config.WithKind(kind),
	}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		options = append(options, config.WithName(name))
	}
	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		options = append(options, config.WithDescription(description))
	}
	if cmd.Flags().Changed("unit") {
		unit, _ := cmd.Flags().GetString("unit")
		options = append(options, config.WithUnit(unit))
	}
	if cmd.Flags().Changed("generator") {
		generator, _ := cmd.Flags().GetString("generator")
		options = append(options, config.WithGenerator(generator))
	}
	if cmd.Flags().Changed("value") {
		value, _ := cmd.Flags().GetString("value")
		options = append(options, config.WithValue(value))
	}
	if cmd.Flags().Changed("type") {
		valueType, _ := cmd.Flags().GetString("type")
		options = append(options, config.WithType(valueType))
	}
	if cmd.Flags().Changed("count") {
		count, _ := cmd.Flags().GetInt("count")
		options = append(options, config.WithCount(count))
	}
	if cmd.Flags().Changed("rate") {
		rate, _ := cmd.Flags().GetDuration("rate")
		options = append(options, config.WithRate(rate))
	}
	if cmd.Flags().Changed("attributes") {
		strAttrs, _ := cmd.Flags().GetStringToString("attributes")
		attrs := make(map[string]any, len(strAttrs))
		for k, v := range strAttrs {
			attrs[k] = v
		}
		options = append(options, config.WithMetricAttributes(attrs))
	}

	mc := config.NewMetricTask(options...)

	if err := mc.Validate(); err != nil {
		return nil, err
	}

	slog.Debug(
		"loaded config",
		"metric", mc.Name,
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
