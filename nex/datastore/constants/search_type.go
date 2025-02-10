package constants

type SearchType uint32

const (
	SearchTypeInvalid SearchType = iota

	// Search for objects from anybody (1)
	SearchTypePublic

	// unknown (2)
	SearchTypeSendFriend

	// unknown (3)
	SearchTypeSendSpecified

	// unknown (4)
	SearchTypeSendSpecifiedFriend

	// unknown (5)
	SearchTypeSend

	// Search for objects from friends (6)
	SearchTypeFriend

	// unknown (7)
	SearchTypeReceivedSpecified

	// unknown (8)
	SearchTypeReceived

	// unknown (9)
	SearchTypePrivate

	// unknown (10)
	SearchTypeOwn

	// unknown (was not set in c# constants?)
	SearchType11

	// unknown (12)
	SearchTypeOwnPending

	// unknown (13)
	SearchTypeOwnRejected

	// unknown (14)
	SearchTypeOwnAll
)
