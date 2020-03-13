package protocol

import (
	"testing"
)

func TestPostConversationsHandler(t *testing.T) {
	tests := []struct {
		Name            string
		CaretStartDelta int
		CaretEndDelta   int
		DocDelta        int

		SenderCaret              Caret
		ReceiverCaret            Caret
		ExpectedReceiverNewCaret Caret
	}{}

	tests = append(tests, addTests...)
	tests = append(tests, deleteTests...)
	tests = append(tests, replaceTests...)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			newCaret := ShiftCaret(
				test.ReceiverCaret,
				test.SenderCaret,
				test.CaretStartDelta,
				test.CaretEndDelta,
				test.DocDelta,
			)

			if test.ExpectedReceiverNewCaret != newCaret {
				t.Errorf(
					"Updated caret is wrong. Expected: %+v. Actual: %+v.",
					test.ExpectedReceiverNewCaret,
					newCaret,
				)
			}
		})
	}
}
