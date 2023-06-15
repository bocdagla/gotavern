package quest_test

import (
	"testing"

	"github.com/bocdagla/tavern/adventurersguild/pkg/domain/quest"
	"github.com/google/uuid"
)

func TestNewQuest_ValidCase(t *testing.T) {
	holder := uuid.New()
	description := "Sample quest description"

	_, err := quest.NewQuest(holder, description)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestNewQuest_EmptyDescription(t *testing.T) {
	holder := uuid.New()
	description := ""

	q, err := quest.NewQuest(holder, description)
	if err != quest.ErrDescriptionEmpty {
		t.Errorf("Expected ErrDescriptionEmpty, but got %v", err)
	}

	if q != (quest.Quest{}) {
		t.Error("Expected empty quest, but got a non-empty quest")
	}
}

func TestNewQuest_InvalidHolder(t *testing.T) {
	holder := uuid.Nil
	description := "Sample quest description"

	q, err := quest.NewQuest(holder, description)
	if err != quest.ErrHolderInvalid {
		t.Errorf("Expected ErrHolderInvalid, but got %v", err)
	}

	if q != (quest.Quest{}) {
		t.Error("Expected empty quest, but got a non-empty quest")
	}
}
