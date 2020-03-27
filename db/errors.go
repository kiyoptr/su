package db

import (
	"strings"

	"github.com/ShrewdSpirit/su/errors"
)

func ErrOpen(err error) error {
	return errors.Newis(1, err, "failed to open db connection")
}

func errFactory(msg string, err error, m Model, extra ...string) error {
	sb := strings.Builder{}
	for i := 0; i < len(extra); i++ {
		if i == 0 {
			sb.WriteString(": ")
		}

		sb.WriteString(extra[i])

		if i+1 != len(extra) {
			sb.WriteString(" ")
		}
	}

	return errors.Newfis(2, err, "%s model %s%s", msg, m.ModelName(), sb.String())
}

func ErrCreate(err error, m Model, extra ...string) error {
	return errFactory("failed to create", err, m, extra...)
}

func ErrUpdate(err error, m Model, extra ...string) error {
	return errFactory("failed to update", err, m, extra...)
}

func ErrDelete(err error, m Model, extra ...string) error {
	return errFactory("failed to delete", err, m, extra...)
}

func ErrQuery(err error, m Model, extra ...string) error {
	return errFactory("failed to query", err, m, extra...)
}
