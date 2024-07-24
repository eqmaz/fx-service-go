package fiberHandlers

import (
	"fmt"
	"strings"

	"fx-service/internal/service/providers"
	"fx-service/internal/service/rates"
	"fx-service/internal/service/stats"
	"fx-service/pkg/config"
	util "fx-service/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// validateAndParseCurrencies checks for case sensitivity and whether the given currencies are supported
func validateAndParseCurrencies(cfg *config.Config, c *fiber.Ctx) (string, string, error) {
	// Get the currency codes from the URL
	ccyBase := c.Params("from")
	ccyQuote := c.Params("to")
	if ccyBase == "" || ccyQuote == "" {
		return "", "", fmt.Errorf("missing base or quote currency codes. Ensure URL and query is correct")
	}

	// Check if we need to make the currencies uppercase
	if !cfg.CurrenciesCaseSensitive {
		ccyBase = strings.ToUpper(ccyBase)
		ccyQuote = strings.ToUpper(ccyQuote)
	}

	if !cfg.IsCurrencySupported(ccyBase) {
		return "", "", fmt.Errorf("invalid base currency code, %s", ccyBase)
	}

	if !cfg.IsCurrencySupported(ccyQuote) {
		return "", "", fmt.Errorf("invalid quote currency code, %s", ccyQuote)
	}

	return ccyBase, ccyQuote, nil
}

// GetRate returns the exchange rate between two currencies
func GetRate(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ccyBase, ccyQuote, err := validateAndParseCurrencies(cfg, ctx)
		if err != nil {
			return replyError(ctx, http.StatusBadRequest, err.Error())
		}

		rateResult, err := rates.GetRate(ccyBase, ccyQuote, cfg.Mode)
		if err != nil {
			return replyError(ctx, http.StatusInternalServerError, err.Error())
		}

		result := fiber.Map{
			"base":   ccyBase,
			"quote":  ccyQuote,
			"rate":   rateResult.Rate,
			"cached": rateResult.WasCached,
		}

		if cfg.ShowProvider {
			result["provider"] = rateResult.Provider
		}

		return replyResult(ctx, result)
	}
}

// GetRates returns the exchange rates between a base currency and multiple quote currencies
func GetRates(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ccyBase := ctx.Query("base", "")
		quotes := ctx.Query("quote", "")
		if ccyBase == "" || quotes == "" {
			return replyError(ctx, http.StatusBadRequest, "Missing base or quote currency codes. Ensure URL and query is correct")
		}

		// Split the comma-delimited list of quote currencies
		ccyQuoteList := strings.Split(quotes, ",")

		// Check if we need to make the currencies uppercase
		if !cfg.CurrenciesCaseSensitive {
			ccyBase = strings.ToUpper(ccyBase)
			for i, ccy := range ccyQuoteList {
				ccyQuoteList[i] = strings.ToUpper(ccy)
			}
		}

		// Validate the currencies
		if !cfg.IsCurrencySupported(ccyBase) {
			return replyError(ctx, http.StatusBadRequest, "Invalid base currency code, "+ccyBase)
		}
		for _, ccy := range ccyQuoteList {
			if !cfg.IsCurrencySupported(ccy) {
				return replyError(ctx, http.StatusBadRequest, "Invalid quote currency code, "+ccy)
			}
		}

		// Get the rates from the provider (or from the cache) using the current strategy
		rateResult, err := rates.GetRates(ccyBase, ccyQuoteList, cfg.Mode)
		if err != nil {
			// TODO parse the different kinds of error and give friendly API responses
			//  instead of just returning the error message to the front end
			return replyError(ctx, http.StatusInternalServerError, err.Error())
		}

		result := fiber.Map{
			"base":   ccyBase,
			"quotes": rateResult.Rates,
			"cached": rateResult.WasCached,
		}

		if cfg.ShowProvider {
			result["provider"] = rateResult.Provider
		}

		return replyResult(ctx, result)
	}
}

func GetStatus(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		modeName := cfg.Mode.String()
		enabled := util.GetMapKeys(providers.EnabledProviders)
		available := util.GetMapKeys(providers.InstalledProviders)
		return replyResult(c, fiber.Map{
			"mode":  modeName,
			"stats": stats.GetInstance().GetStats(),
			"providers": fiber.Map{
				"enabled":   enabled,
				"available": available,
			},
		})
	}
}

func HealthCheck(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return replyResult(c, fiber.Map{
			"status": "healthy",
		})
	}
}
