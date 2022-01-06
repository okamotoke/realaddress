package realaddress

import (
	"testing"
)

func Test_getRandomAddress(t *testing.T) {
	var (
		testAddressFile = "data/test1.CSV"
		lienCount       = 3
	)

	t.Run("should return 9071801", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			got, err := getRandomAddress(testAddressFile, lienCount)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.GetPostalCode() != "9071801" {
				t.Errorf("unexpected postal code %s", got.GetPostalCode())
			}

			if got.GetPrefecture() != "沖縄県" {
				t.Errorf("unexpected prefecture %s", got.GetPrefecture())
			}

			if got.GetCity() != "八重山郡　与那国町" {
				t.Errorf("unexpected city %s", got.GetCity())
			}

			if got.GetTown() != "与那国" {
				t.Errorf("unexpected town %s", got.GetTown())
			}
		}
	})
}
