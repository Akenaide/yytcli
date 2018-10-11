package lib

import (
	"testing"
)

const NbCard = 117

func TestImport(t *testing.T) {
	cardMap := GetCards([]string{"tmhmovie"})
	if len(cardMap) != NbCard {
		t.Errorf("Errors in parsing expected %d but got %d", NbCard, len(cardMap))
	}
}
