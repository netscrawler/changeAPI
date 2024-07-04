package converter

import (
	"context"
	"converter/internal/lib/logger/sl"
	"converter/internal/lib/rateExtract"
	"fmt"
	"log/slog"
)

type Converter struct {
	log *slog.Logger
}

// New returns a new instance of the Converter service.
func New(
	log *slog.Logger,

) *Converter {
	return &Converter{
		log: log,
	}
}

func (c *Converter) Convert(
	ctx context.Context,
	amount uint32,
	currency string,
) (uint32, float32, error) {
	const op = "converter.convert"
	log := c.log.With(
		slog.String("op", op),
		slog.String("currency", currency),
		slog.Any("amount", amount),
	)
	log.Info("convertation")
	var convertedAmount uint32
	//TODO extract rate
	rate, err := rateExtract.GetExchangeRate(currency)
	if err != nil {
		log.Error("Error extracting rate", sl.Err(err))
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	convertedAmount = uint32(float32(amount) / rate.Rate)
	return convertedAmount, rate.Rate, nil
}
