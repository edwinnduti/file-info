# File-info API
An API that takes in a file and gives back info about the file.

The data is stored in a mysql database but the file is discarded.

### Requirements
* Golang.
* Mysql.
* Postman or REST CLIENT from Vscode for testing.


To run it locally, <b>edit the .env file</b>, set up myql connection, create database called fileinfodb

To create table property, run the command:
```
$ mysql -u yourusername -p dbnamehere < database/fileinfodb.sql
```

e.g file - hello.py

request URI locally:
http://localhost:8080/upload

#### Response:
```json
{
    "id":7,
    "name":"hello",
    "extension":".py",
    "size":128,
    "type":"application/octet-stream"
}
```


 Size is in bytes always.

 ### Available paths:

| function              |   path                    |   method  |
|   ----                |   ----                    |   ----    |
| welcome           |   /			|	    GET    |
| upload file       |   /upload			|	POST     |


 Happy coding.
