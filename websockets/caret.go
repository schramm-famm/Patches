package websockets

type Caret struct {
	Start int `json:"caret_start"`
	End   int `json:"caret_end"`
}

func processAdd(
	receiverCaret Caret,
	senderCaret Caret,
	caretStartDelta int,
	caretEndDelta int,
	docDelta int,
) Caret {
	if senderCaret.End < receiverCaret.End && caretEndDelta != 0 {
		receiverCaret.End += docDelta
		if senderCaret.End <= receiverCaret.Start {
			receiverCaret.Start += docDelta
		}
	}
	return receiverCaret
}

func shiftCaret(
	receiverCaret Caret,
	senderCaret Caret,
	caretStartDelta int,
	caretEndDelta int,
	docDelta int,
) Caret {
	var rangeStart, rangeEnd int
	if senderCaret.Start != senderCaret.End {
		rangeStart = senderCaret.Start
		rangeEnd = senderCaret.End
	} else if docDelta > 0 {
		rangeStart = senderCaret.Start
		rangeEnd = senderCaret.Start + caretStartDelta
	} else if caretStartDelta == 0 {
		rangeStart = senderCaret.Start
		rangeEnd = senderCaret.Start - docDelta
	} else {
		rangeStart = senderCaret.Start + caretStartDelta
		rangeEnd = senderCaret.Start
	}

	if docDelta > 0 && senderCaret.Start == senderCaret.End {
		return processAdd(receiverCaret, senderCaret, caretStartDelta, caretEndDelta, docDelta)
	}

	delta := rangeStart - rangeEnd

	if rangeEnd < receiverCaret.End {
		receiverCaret.End += delta
		if rangeEnd <= receiverCaret.Start {
			receiverCaret.Start += delta
		} else {
			if rangeStart < receiverCaret.Start {
				receiverCaret.Start = rangeStart
			}
		}
	} else {
		if rangeEnd == receiverCaret.End {
			receiverCaret.Start = rangeEnd + delta
			receiverCaret.End = rangeEnd + delta
		} else if rangeStart < receiverCaret.End {
			receiverCaret.End = rangeStart
			if rangeStart < receiverCaret.Start {
				receiverCaret.Start = rangeStart
			}
		}
	}

	if caretStartDelta > 0 {
		tmpSenderCaret := Caret{
			Start: senderCaret.Start,
			End:   senderCaret.Start,
		}
		return processAdd(receiverCaret, tmpSenderCaret, caretStartDelta, caretStartDelta, caretStartDelta)
	}
	return receiverCaret
}
