package toolkit

import (
	"fmt"
	"testing"
	"time"

	"github.com/PKUJohnson/solar/std/toolkit/dateutil"
)

func TestGetStandarTime(t *testing.T) {
	fmt.Println(dateutil.GetStandarTime(time.Now()))
}
