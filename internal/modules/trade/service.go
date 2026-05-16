package trade

import (
	"bcv/internal/domain"
	"bcv/internal/platform/providers/binance"
	"context"
	"errors"
	"log/slog"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

func FetchBinanceRates() (*domain.BinanceRate, error) {
	var sellPrice, buyPrice decimal.Decimal
	c := context.Background()
	g, _ := errgroup.WithContext(c)

	g.Go(func() error {
		var err error
		buyPrice, err = binance.GetBinanceRates("BUY", 50000)
		return err
	})

	g.Go(func() error {
		var err error
		sellPrice, err = binance.GetBinanceRates("SELL", 50000, "Banesco", "PagoMovil", "BANK", "BancoDeVenezuela", "Mercantil", "Bancamiga", "Provincial", "BNCBancoNacional", "BBVABank", "Bancaribe", "Banplus", "SpecificBank", "BancoPlaza", "BancoVeneCredit", "BancoDelTesoro", "BancoActivo", "BFC", "BDDT", "N58")
		return err
	})

	if err := g.Wait(); err != nil {
		slog.Error("Error fetching Binance rates", "error", err)
		return nil, errors.New("Error fetching Binance rates")
	}

	return &domain.BinanceRate{
		SellPrice: sellPrice.InexactFloat64(),
		BuyPrice:  buyPrice.InexactFloat64(),
	}, nil
}
