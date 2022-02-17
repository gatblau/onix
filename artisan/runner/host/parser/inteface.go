package parser

type Event interface {
    GetObjectDownloadURL() (string, error)
}