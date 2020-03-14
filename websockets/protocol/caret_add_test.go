package protocol

var addTests = []struct {
	Name            string
	CaretStartDelta int
	CaretEndDelta   int
	DocDelta        int

	SenderCaret              Caret
	ReceiverCaret            Caret
	ExpectedReceiverNewCaret Caret
}{
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
}
