package constants

const (
	LOAN_STATUS_NEW         = 1
	LOAN_STATUS_IN_PROGRESS = 2
	LOAN_STATUS_DELINQUENT  = 3
	LOAN_STATUS_COMPLETE    = 4
)

var (
	MAP_LOAN_STATUS = map[int16]string{
		LOAN_STATUS_NEW:         "New",
		LOAN_STATUS_IN_PROGRESS: "In Progress",
		LOAN_STATUS_DELINQUENT:  "Delinquent",
		LOAN_STATUS_COMPLETE:    "Complete",
	}
)
