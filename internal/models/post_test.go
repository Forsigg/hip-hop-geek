package models

import (
	"log"
	"strings"
	"testing"
	"time"

	"hip-hop-geek/internal/types"
)

func TestPost(t *testing.T) {
	post := Post{
		1, "21 Savage - American Dream", types.NewCustomDate(2024, time.January, 12),
	}

	t.Run("check Artist work correct", func(t *testing.T) {
		expected := "21 Savage"
		got := post.Artist()

		if got != expected {
			t.Errorf("incorrect artist, got '%s', want '%s'", got, expected)
		}
	})

	t.Run("check Title work correct", func(t *testing.T) {
		expected := "American Dream"
		got := post.Title()

		if got != expected {
			t.Errorf("incorrect title, got '%s', want '%s'", got, expected)
		}
	})

	t.Run("check Query work correct", func(t *testing.T) {
		expected := "21 Savage - American Dream"
		got := post.Query()

		if got != expected {
			t.Errorf("incorrect query, got '%s', want '%s'", got, expected)
		}
	})

	t.Run("check OutDate work correct", func(t *testing.T) {
		expected := types.NewCustomDate(2024, time.January, 12)
		got := post.ReleaseDate()

		if got != expected {
			t.Errorf("incorrect release date, got '%s', want '%s'", got, expected)
		}
	})

	t.Run("check Artist after Title", func(t *testing.T) {
		expectedArtist := "21 Savage"
		expectedTitle := "American Dream"

		gotArtist := post.Artist()
		gotTitle := post.Title()

		if gotArtist != expectedArtist || gotTitle != expectedTitle {
			t.Errorf(
				"incorrect artist or title: got '%s' artist, '%s' title, but want '%s' artist '%s' title",
				gotArtist,
				gotTitle,
				expectedArtist,
				expectedTitle,
			)
		}
	})

	t.Run("check – is not - case", func(t *testing.T) {
		postCheck := Post{
			1, "Dave East – Fortune Favors the Bold", types.NewCustomDate(2024, time.January, 1),
		}

		expectedArtist := "Dave East"
		expectedTitle := "Fortune Favors the Bold"

		log.Println(strings.Contains(postCheck.QueryField, "–"))

		gotArtist := postCheck.Artist()
		gotTitle := postCheck.Title()

		if gotArtist != expectedArtist || gotTitle != expectedTitle {
			t.Errorf(
				"incorrect artist or title: got '%s' artist, '%s' title, but want '%s' artist '%s' title",
				gotArtist,
				gotTitle,
				expectedArtist,
				expectedTitle,
			)
		}
	})
}
