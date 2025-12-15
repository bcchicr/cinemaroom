package services

import (
	"testing"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure/services/password"
)

func TestPasswordHashingAndVerification(t *testing.T) {
	somePassword, _ := vo.NewPassword("some_password")
	otherPassword, _ := vo.NewPassword("other_password")

	s := password.NewService()

	t.Run("same password hash", func(t *testing.T) {
		someHash, err := s.Hash(somePassword)
		if err != nil {
			t.Fatalf("expected correct hashing, got error")
		}

		isVerified := s.Verify(somePassword, someHash)
		if isVerified == false {
			t.Fatalf("expected hash verification")
		}
	})

	t.Run("other password hash", func(t *testing.T) {
		otherHash, err := s.Hash(otherPassword)
		if err != nil {
			t.Fatalf("expected correct hashing, got error")
		}

		isVerified := s.Verify(somePassword, otherHash)
		if isVerified == true {
			t.Fatalf("not expected hash verification")
		}
	})

}
