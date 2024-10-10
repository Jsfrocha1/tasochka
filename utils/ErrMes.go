package utils

type TErrMessages struct {
	InvalidData             string
	InvalidRequest          string
	MethodNotAllowed        string
	MissingFields           string
	RecordNotFound          string
	FailedToCreatePerson    string
	FailedToRetrieveRecords string
	InvalidIDFormat         string
}

var ErrMessages TErrMessages = TErrMessages{
	InvalidData:             "Invalid data",
	InvalidRequest:          "Invalid Message",
	MethodNotAllowed:        "Method not allowed",
	MissingFields:           "Missing fields",
	RecordNotFound:          "Record not found",
	FailedToCreatePerson:    "Failed to create person",
	FailedToRetrieveRecords: "Failed to retrieve records",
	InvalidIDFormat:         "Invalid ID format",
} 