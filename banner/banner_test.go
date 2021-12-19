package banner

import (
	"testing"
	"time"
)

func TestBannerShow(t *testing.T) {
	Show("1.0.0", "1.17.3", time.Now().Format("2006-01-02 15:04:05"), "Jeff")
}
