package ginHandlers

import (
	"fmt"
	"fx-service/internal/service/providers"
	"fx-service/internal/service/rates"
	"fx-service/internal/service/stats"
	"fx-service/pkg/config"
	util "fx-service/pkg/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// validateAndParseCurrencies checks for case sensitivity and whether the given currencies are supported
func validateAndParseCurrencies(cfg *config.Config, c *gin.Context) (string, string, error) {
	ccyBase := c.Param("from")
	ccyQuote := c.Param("to")
	if ccyBase == "" || ccyQuote == "" {
		return "", "", fmt.Errorf("missing base or quote currency codes. Ensure URL and query is correct")
	}

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

func GetRate(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ccyBase, ccyQuote, err := validateAndParseCurrencies(cfg, c)
		if err != nil {
			replyError(c, http.StatusBadRequest, err.Error())
			return
		}

		rateResult, err := rates.GetRate(ccyBase, ccyQuote, cfg.Mode)
		if err != nil {
			replyError(c, http.StatusInternalServerError, err.Error())
			return
		}

		result := gin.H{
			"base":   ccyBase,
			"quote":  ccyQuote,
			"rate":   rateResult.Rate,
			"cached": rateResult.WasCached,
		}

		if cfg.ShowProvider {
			result["provider"] = rateResult.Provider
		}

		replyResult(c, result)
	}
}

func GetRates(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ccyBase := c.Query("base")
		quotes := c.Query("quote")
		if ccyBase == "" || quotes == "" {
			replyError(c, http.StatusBadRequest, "Missing base or quote currency codes. Ensure URL and query is correct")
			return
		}

		ccyQuoteList := strings.Split(quotes, ",")

		if !cfg.CurrenciesCaseSensitive {
			ccyBase = strings.ToUpper(ccyBase)
			for i, ccy := range ccyQuoteList {
				ccyQuoteList[i] = strings.ToUpper(ccy)
			}
		}

		if !cfg.IsCurrencySupported(ccyBase) {
			replyError(c, http.StatusBadRequest, "Invalid base currency code, "+ccyBase)
			return
		}
		for _, ccy := range ccyQuoteList {
			if !cfg.IsCurrencySupported(ccy) {
				replyError(c, http.StatusBadRequest, "Invalid quote currency code, "+ccy)
				return
			}
		}

		rateResult, err := rates.GetRates(ccyBase, ccyQuoteList, cfg.Mode)
		if err != nil {
			replyError(c, http.StatusInternalServerError, err.Error())
			return
		}

		result := gin.H{
			"base":   ccyBase,
			"quotes": rateResult.Rates,
			"cached": rateResult.WasCached,
		}

		if cfg.ShowProvider {
			result["provider"] = rateResult.Provider
		}

		replyResult(c, result)
	}
}

func GetStatus(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		modeName := cfg.Mode.String()
		enabled := util.GetMapKeys(providers.EnabledProviders)
		available := util.GetMapKeys(providers.InstalledProviders)
		replyResult(c, gin.H{
			"mode":  modeName,
			"stats": stats.GetInstance().GetStats(),
			"providers": gin.H{
				"enabled":   enabled,
				"available": available,
			},
		})
	}
}

func HealthCheck(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	}
}
