package websockets

var replaceTests = []struct {
	Name            string
	CaretStartDelta int
	CaretEndDelta   int
	DocDelta        int

	SenderCaret              Caret
	ReceiverCaret            Caret
	ExpectedReceiverNewCaret Caret
}{
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
