package creditcard

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// various credit card errors
var (
	ErrExpired   = errors.New("credit card expired")
	ErrInvalid   = errors.New("credit card invalid")
	ErrUnaccpted = errors.New("credit card unaccepted")
)

// Luhn uses a checksum formula to validate numbers,
// mainly single-digit error,
// from https://en.wikipedia.org/wiki/Luhn_algorithm.
func Luhn(purported string) bool {
	sum := 0
	l := len(purported)
	if l < 2 {
		return false
	}
	parity := l % 2 // the check digit is at odd or even
	for i, n := range purported {
		d, err := strconv.Atoi(string(n))
		if err != nil {
			return false
		}
		if i%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return (sum%10 == 0)
}

// Expiration is the expire date on the credit card.
type Expiration struct {
	Year  int
	Month int
}

// Expired checkes if the date is expired.
func (e Expiration) Expired() bool {
	if e.Year == 0 || e.Month == 0 {
		return true
	}
	t := time.Now()
	return t.Year() > e.Year || (t.Year() == e.Year && int(t.Month()) > e.Month)
}

// Card represents a credit card.
type Card struct {
	Number string
	Cvv    int32
	Exp    Expiration
	// internal
	err error
	msg string
}

func (c *Card) first(n int) int64 {
	if len(c.Number) < n {
		return 0
	}
	i, err := strconv.ParseInt(c.Number[0:n], 0, 64)
	if err != nil {
		return 0
	}
	return i
}

// known brand codes
const (
	MasterCard = "master"
	VisaCard   = "visa"
)

// Brand tries to get the brand code from the card number.
// More from https://creditcardjs.com/credit-card-type-detection.
func (c *Card) Brand() string {
	if c.Number[0] == '4' {
		return VisaCard
	}
	f2 := c.first(2)
	if f2 >= 50 && f2 <= 55 {
		return MasterCard
	}
	return ""
}

// Validate validates the credit card and
// return true if it passes all the test.
// Otherwise use Err() and Msg() to get
// the error and detail about the first
// failed validation.
func (c *Card) Validate() bool {
	// checksum
	if !Luhn(c.Number) {
		c.err = ErrInvalid
		c.msg = fmt.Sprintf("Invalid credit card number: %s; might be a typo", c.Number)
		return false
	}
	// brand filter
	b := c.Brand()
	switch b {
	case MasterCard, VisaCard:
		// pass
	default:
		c.err = ErrUnaccpted
		c.msg = fmt.Sprintf("Sorry only VISA or MasterCard is accepted")
		return false
	}
	// expire
	if c.Exp.Expired() {
		c.err = ErrExpired
		c.msg = fmt.Sprintf("Your credit card expired on %d/%d", c.Exp.Month, c.Exp.Year)
		return false
	}
	return true
}

// Err returns the error during validation.
func (c *Card) Err() error {
	return c.err
}

// Msg returns the detail error message.
func (c *Card) Msg() string {
	return c.msg
}
