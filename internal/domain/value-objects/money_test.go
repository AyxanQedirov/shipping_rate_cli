package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"shipping_rate_cli/internal/domain/enums"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
)

type MoneySuite struct {
	suite.Suite
}

func (this *MoneySuite) Test_Add_RetursNil_InCaseOfMatchingCurrencies() {
	//Arrange
	money := valueobjects.Money{Amount: 10.25, Currency: enums.EUR}
	//Act
	err := money.Add(valueobjects.Money{Amount: 8.50, Currency: enums.EUR})

	//Assert
	this.NoError(err)
	this.Equal(18.75, money.Amount)
	this.Equal(enums.EUR, money.Currency)
}

func (this *MoneySuite) Test_Add_RetursError_InCaseOfDifferentCurrencies() {
	//Arrange
	money := valueobjects.Money{Amount: 10, Currency: enums.EUR}
	original := money

	//Act
	err := money.Add(valueobjects.Money{Amount: 8, Currency: enums.Currency("USD")})

	//Assert
	this.EqualError(err, "adding money failed: currencies are not the same and multicurrency is not supported")
	this.Equal(original, money, "failed addition must not mutate the amount or currency")
}

func TestMoneySuite(t *testing.T) {
	suite.Run(t, new(MoneySuite))
}
