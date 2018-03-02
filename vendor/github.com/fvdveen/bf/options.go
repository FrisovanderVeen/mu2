package bf

import (
	"io"
)

// OptionFunc is a option for the bot
type OptionFunc func(*Bot) error

// EmptyOptionFunc is a OptionFunc that does nothing
// it is used when an error occurs
func EmptyOptionFunc(b *Bot) error {
	return nil
}

// Token sets the token of the bot
func Token(token string) OptionFunc {
	return func(b *Bot) error {
		b.Token = token
		return nil
	}
}

// Prefix sets the prefix of the bot
func Prefix(prefix string) OptionFunc {
	return func(b *Bot) error {
		b.Prefix = prefix
		return nil
	}
}

// ErrWriter sets the writer errors are written to
func ErrWriter(w io.Writer) OptionFunc {
	return func(b *Bot) error {
		b.ErrWriter = w
		return nil
	}
}

// ErrPrefix sets the function that will be called and prepended to the error
func ErrPrefix(prefunc func() string) OptionFunc {
	return func(b *Bot) error {
		b.ErrPrefix = prefunc
		return nil
	}
}
