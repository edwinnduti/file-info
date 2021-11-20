package model

// response struct
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// describe config model
type Config struct {
	Host       string
	Dbport     string
	Dbusername string
	Dbname     string
	Passwd     string
}
