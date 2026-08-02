package realtimetypes

import sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"

type YjsDocumentState = sharedtypes.YjsDocumentState
type YjsDocumentUpdate = sharedtypes.YjsDocumentUpdate

func MarshalYjsUpdateSequence(updateSequence int64) []byte {
	return sharedtypes.MarshalYjsUpdateSequence(updateSequence)
}
