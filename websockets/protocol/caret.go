package protocol

// Caret represents a user's Start and End position in the document
type Caret struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// processAdd takes a recceiverCaret and returns a shifted Caret according to a
// specified delta.
func (receiverCaret Caret) processAdd(senderCaret Caret, delta int) Caret {
	if senderCaret.End < receiverCaret.End {
		receiverCaret.End += delta
		if senderCaret.End <= receiverCaret.Start {
			receiverCaret.Start += delta
		}
	}
	return receiverCaret
}

// ShiftCaret takes a receiverCaret and returns a shifted Caret according to
// the delta values and the location of the senderCaret.
func (receiverCaret Caret) ShiftCaret(senderCaret Caret, delta Delta) Caret {
	var rangeStart, rangeEnd int
	if senderCaret.Start != senderCaret.End {
		rangeStart = senderCaret.Start
		rangeEnd = senderCaret.End
	} else if *delta.Doc > 0 {
		rangeStart = senderCaret.Start
		rangeEnd = senderCaret.Start + *delta.CaretStart
	} else if *delta.CaretStart == 0 {
		rangeStart = senderCaret.Start
		rangeEnd = senderCaret.Start - *delta.Doc
	} else {
		rangeStart = senderCaret.Start + *delta.CaretStart
		rangeEnd = senderCaret.Start
	}

	if *delta.Doc > 0 && senderCaret.Start == senderCaret.End {
		return receiverCaret.processAdd(senderCaret, *delta.Doc)
	}

	rangeDelta := rangeStart - rangeEnd

	if rangeEnd < receiverCaret.End {
		receiverCaret.End += rangeDelta
		if rangeEnd <= receiverCaret.Start {
			receiverCaret.Start += rangeDelta
		} else if rangeStart < receiverCaret.Start {
			receiverCaret.Start = rangeStart
		}
	} else if rangeEnd == receiverCaret.End {
		receiverCaret.Start = rangeEnd + rangeDelta
		receiverCaret.End = rangeEnd + rangeDelta
	} else if rangeStart < receiverCaret.End {
		receiverCaret.End = rangeStart
		if rangeStart < receiverCaret.Start {
			receiverCaret.Start = rangeStart
		}
	}

	if *delta.CaretStart > 0 {
		tmpSenderCaret := Caret{
			Start: senderCaret.Start,
			End:   senderCaret.Start,
		}
		return receiverCaret.processAdd(tmpSenderCaret, *delta.CaretStart)
	}
	return receiverCaret
}
