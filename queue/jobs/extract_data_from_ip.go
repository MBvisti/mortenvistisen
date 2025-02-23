package jobs

import "github.com/google/uuid"

const extractDataFromIPJobKind string = "ExtractDataFromIPJob"

type ExtractDataFromIPJob struct {
	IPAddress string    `json:"ip_address"`
	SessionID uuid.UUID `json:"session_id"`
}

func (ExtractDataFromIPJob) Kind() string { return extractDataFromIPJobKind }
