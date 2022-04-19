package conversion

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestThatConverterRegistryReturnsOnlyConvertersThatMatchType(t *testing.T) {
	is, conReg := testSetup(t)

	mcs := conReg.DesignateConverters(context.Background(), []string{"urn:oma:lwm2m:ext:3303", "humidity"})
	is.Equal(len(mcs), 1)
}

func TestThatConverterRegistryReturnsEmptyIfNoTypesMatch(t *testing.T) {
	is, conReg := testSetup(t)

	mcs := conReg.DesignateConverters(context.Background(), []string{"co2"})
	is.Equal(len(mcs), 0)
}

func testSetup(t *testing.T) (*is.I, ConverterRegistry) {
	is := is.New(t)

	conReg := NewConverterRegistry()

	return is, conReg
}
