package domain

const (
	headerSize = 12
)

/* header section offsets */
const (
	packetIdentifierOffset       = 0
	queryResponseIndicatorOffset = 2
	questionCountOffset          = 4
	answerRecordCountOffset      = 6
	authorityRecordCountOffset   = 8
	additionalRecordCountOffset  = 10
)

/* header fields masks */
const (
	queryResponseIndicatorMask = 0x8000
	operationCodeMask          = 0x7800
	authoritativeAnswerMask    = 0x0400
	truncationMask             = 0x0200
	recursionDesiredMask       = 0x0100
	recursionAvailableMask     = 0x0080
	reservedMask               = 0x0070
	responseCodeMask           = 0x000F
)

/* question section offsets */
const (
	recordClassOffset = 2
)
