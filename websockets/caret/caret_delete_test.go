package caret

var deleteTests = []struct {
	Name            string
	CaretStartDelta int
	CaretEndDelta   int
	DocDelta        int

	SenderCaret              Caret
	ReceiverCaret            Caret
	ExpectedReceiverNewCaret Caret
}{
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
}
