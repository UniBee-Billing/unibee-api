// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"unibee/internal/dao/default/internal"
)

// internalCountryRateDao is internal type for wrapping internal DAO implements.
type internalCountryRateDao = *internal.CountryRateDao

// countryRateDao is the data access object for table country_rate.
// You can define custom methods on it to extend its functionality as you wish.
type countryRateDao struct {
	internalCountryRateDao
}

var (
	// CountryRate is globally public accessible object for table country_rate operations.
	CountryRate = countryRateDao{
		internal.NewCountryRateDao(),
	}
)

// Fill with you ideas below.