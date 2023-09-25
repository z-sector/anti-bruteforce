package models

type AuthCheck struct {
	Login    string
	Password string
	IP       string
}

type ResetBucketData struct {
	Login string
	IP    string
}
