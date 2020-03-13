package websockets

import (
	"fmt"
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
	}{
		// ADD TESTS
		{
			Name:            "Add::Single::Receiver Cursor::Sender After Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Add::Single::Receiver Cursor::Sender Before Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 2,
				End:   2,
			},
		},
		{
			Name:            "Add::Single::Receiver Cursor::Same Position",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Add::Multiple::Receiver Cursor::Sender After Receiver",
			CaretStartDelta: 2,
			CaretEndDelta:   2,
			DocDelta:        2,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Add::Multiple::Receiver Cursor::Sender Before Receiver",
			CaretStartDelta: 2,
			CaretEndDelta:   2,
			DocDelta:        2,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 3,
				End:   3,
			},
		},
		{
			Name:            "Add::Multiple::Receiver Cursor::Same Position",
			CaretStartDelta: 2,
			CaretEndDelta:   2,
			DocDelta:        2,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Add::Single::Receiver Range::Sender After Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 3,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Add::Single::Receiver Range::Sender Before Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 2,
				End:   4,
			},
		},
		{
			Name:            "Add::Single::Receiver Range::Sender Same as Receiver Start",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   3,
			},
		},
		{
			Name:            "Add::Single::Receiver Range::Sender Same as Receiver End",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 2,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Add::Single::Receiver Range::Sender After Receiver Start & Before Receiver End",
			CaretStartDelta: 1,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   3,
			},
		},

		// DELETE
		{
			Name:            "Delete::Sender Cursor::Backward::Single::Receiver Cursor::Sender After Receiver",
			CaretStartDelta: -1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Cursor:Backward::Single::Receiver Cursor::Sender Before Receiver",
			CaretStartDelta: -1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Single::Receiver Cursor::Same Position",
			CaretStartDelta: -1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Cursor::Deleted Content After Receiver",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 4,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 2,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Cursor:Backward::Multiple::Receiver Cursor::Deleted Content Before Receiver",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 3,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 3,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Cursor::Deleted Content Overlaps Receiver",
			CaretStartDelta: -3,
			CaretEndDelta:   -3,
			DocDelta:        -3,
			SenderCaret: Caret{
				Start: 3,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Range::Deleted Content After Receiver",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 5,
				End:   5,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   3,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Range::Deleted Content Before Receiver",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 3,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 3,
				End:   5,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   3,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Range::Deleted Content Overlaps Receiver End",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 4,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Range::Deleted Content Overlaps Receiver Start",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 2,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Range::Deleted Content Overlaps Receiver Start & End",
			CaretStartDelta: -4,
			CaretEndDelta:   -4,
			DocDelta:        -4,
			SenderCaret: Caret{
				Start: 4,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Range::Deleted Content Same as Receiver Start & End",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 3,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Backward::Multiple::Receiver Range::Deleted Content Surrounded by Receiver",
			CaretStartDelta: -2,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 3,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Range::Backward::Multiple::Receiver Cursor::Sender After Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 2,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Range::Backward::Multiple::Receiver Cursor::Sender Before Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Range::Backward::Multiple::Receiver Cursor::Deleted Content Overlaps Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Range::Backward::Multiple::Receiver Range::Deleted Content After Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 2,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Range::Backward::Multiple::Receiver Range::Deleted Content Before Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Range::Backward::Multiple::Receiver Range::Deleted Content Overlaps Receiver End",
			CaretStartDelta: 0,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 1,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Range::Backward::Multiple::Receiver Range::Deleted Content Overlaps Receiver Start",
			CaretStartDelta: 0,
			CaretEndDelta:   -2,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Single::Receiver Cursor::Sender After Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Single::Receiver Cursor::Sender Before Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Single::Receiver Cursor::Same Position",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Cursor::Sender After Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 2,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Cursor::Sender Before Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 3,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Cursor::Same Position",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Cursor::Deleted Content Overlaps Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Range::Deleted Content After Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 2,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Range::Deleted Content Before Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Range::Deleted Content Overlaps Receiver End",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 2,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   2,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Range::Deleted Content Overlaps Receiver Start",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Range::Deleted Content Overlaps Receiver Start & End",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -4,
			SenderCaret: Caret{
				Start: 0,
				End:   0,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Range::Deleted Content Same as Receiver Start & End",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -3,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   1,
			},
		},
		{
			Name:            "Delete::Sender Cursor::Forward::Multiple::Receiver Range::Deleted Content Surrounded by Receiver",
			CaretStartDelta: 0,
			CaretEndDelta:   0,
			DocDelta:        -2,
			SenderCaret: Caret{
				Start: 1,
				End:   1,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},

		// REPLACE
		{
			Name:            "Replace::Receiver Cursor::Smaller::Sender After Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 2,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   0,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Replace::Receiver Cursor::Smaller::Sender Before Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Replace::Receiver Cursor::Smaller::Deleted Content Overlaps Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   1,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Replace::Receiver Range::Smaller::Sender After Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 2,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   2,
			},
		},
		{
			Name:            "Replace::Receiver Range::Smaller::Sender Before Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 2,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   3,
			},
		},
		{
			Name:            "Replace::Receiver Range::Smaller::Replaced Content Overlaps Receiver End",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 1,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   1, // FUDGE
			},
		},
		{
			Name:            "Replace::Receiver Range::Smaller::Replaced Content Overlaps Receiver Start",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 1,
				End:   2,
			},
		},
		{
			Name:            "Replace::Receiver Range::Smaller::Replaced Content Overlaps Receiver Start & End",
			CaretStartDelta: 1,
			CaretEndDelta:   -3,
			DocDelta:        -3,
			SenderCaret: Caret{
				Start: 0,
				End:   4,
			},
			ReceiverCaret: Caret{
				Start: 1,
				End:   3,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Replace::Receiver Range::Smaller::Replaced Content Same as Receiver Start & End",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 0,
				End:   2,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   2,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   0,
			},
		},
		{
			Name:            "Replace::Receiver Range::Smaller::Replaced Content Surrounded by Receiver",
			CaretStartDelta: 1,
			CaretEndDelta:   -1,
			DocDelta:        -1,
			SenderCaret: Caret{
				Start: 1,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   3,
			},
		},
		{
			Name:            "Replace::Receiver Range::Same Size::Replaced Content Surrounded by Receiver",
			CaretStartDelta: 2,
			CaretEndDelta:   0,
			DocDelta:        0,
			SenderCaret: Caret{
				Start: 1,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   4,
			},
		},
		{
			Name:            "Replace::Receiver Range::Bigger::Replaced Content Surrouned by Receiver",
			CaretStartDelta: 3,
			CaretEndDelta:   1,
			DocDelta:        1,
			SenderCaret: Caret{
				Start: 1,
				End:   3,
			},
			ReceiverCaret: Caret{
				Start: 0,
				End:   4,
			},
			ExpectedReceiverNewCaret: Caret{
				Start: 0,
				End:   5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fmt.Println(test.Name)

			newCaret := shiftCaret(
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
