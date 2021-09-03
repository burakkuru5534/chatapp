Chat Application

Technology:

- Golang
- Vue js
- Sqlite3
- docker
- docker-compose

* How Its working

Room,msg,user,userfriend tables were created with sqlite. With the room table, the channel information that users will message is kept. Users and passwords registered in the user table were hashed and kept. The user friend table was created to manage who users can message and who they block. The msg table was used to hold messages and associate them with users.

Registering, logging in is done by sending a request to the relevant APIs via postman.

A simple frontend design has been created for websocket communication. Here, we open a chat window by entering the token we obtained from the login api and the username of the person we want to msg. Then we send our message.

In order to add friends, block someone, see messages, see the friend list, we send requests to the relevant APIs via postman.

API's

* Register

- EndPoint: http://localhost:8080/api/register?username=burakkuru&password=test123
- Attributes: username is unique. So, if username already taken, user should try something else.

Success Response: 
{
    "ID": "6fa792f4-08af-45e1-8e1e-fead5486af7f",
    "Success": true,
    "Message": "",
    "Data": "5269481e-ddc8-4604-b576-fac3b0d2a7a7"
}

Fail Response: 
{
    "ID": "c475de81-4974-429c-9ff4-55e7733e1b55",
    "Success": false,
    "Message": "username already exist error",
    "Data": null
}

* Login

- EndPoint: http://localhost:8080/api/login
- Attributes: Request Body should include username and password.
Sample Request Body:
{
    "username":"burakkuru",
    "password":"test123"
}

We generate a token with this api.
Success Response: 
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE2MzEyODY2NTIsIklkIjoiOTY3MzEwMmItZjFkZi00YjZkLTlhMjMtOWJmZTQ0NTIzNGY2IiwiTmFtZSI6ImJ1cmFra3VydSJ9.nPkgd59jPqKdDWosiqRLr1wy1GBp2qA5cj28oGA871w

Failed Response: 
{
    "ID": "cd2fa29a-bc67-4412-adde-cc830f8f0370",
    "Success": false,
    "Message": "find user by username error",
    "Data": null
}

* Add Friend

-EndPoint: http://localhost:8080/api/friend/add?fUserName=burakozcelik&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE2MzEyODU4MDUsIklkIjoiOTY3MzEwMmItZjFkZi00YjZkLTlhMjMtOWJmZTQ0NTIzNGY2IiwiTmFtZSI6ImJ1cmFra3VydSJ9.ijelimRpC0aSzQvUPgw1-KQaeEmlbPRcIggizZbi2I0

- Attributes: we should use our token when we send a request to this api. also we should specify friend name to add him/her.
Query String Parameters: ?fUserName=sampleuser&token=ourtoken

* Ban Friend

- EndPoint: http://localhost:8080/api/friend/ban?fUserName=burakozcelik&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE2MzEyODU4MDUsIklkIjoiOTY3MzEwMmItZjFkZi00YjZkLTlhMjMtOWJmZTQ0NTIzNGY2IiwiTmFtZSI6ImJ1cmFra3VydSJ9.ijelimRpC0aSzQvUPgw1-KQaeEmlbPRcIggizZbi2I0

- Attributes: we should use our token when we send a request to this api. also we should specify friend name to ban him/her.
Query String Parameters: ?fUserName=sampleuser&token=ourtoken

* Friends

- EndPoint: http://localhost:8080/api/friends?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE2MzEyODgyNTksIklkIjoiOTY3MzEwMmItZjFkZi00YjZkLTlhMjMtOWJmZTQ0NTIzNGY2IiwiTmFtZSI6ImJ1cmFra3VydSJ9.9flX6FhArFXQel8lRKqzW1jkGG2Fo43q36etZhI16Og

- Attributes: we should use our token when we send a request to this api to see our friends list.

* Show User Sended Messages to Someone:

EndPoint: http://localhost:8080/api/user/sended/messages?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE2MzEyODgyNTksIklkIjoiOTY3MzEwMmItZjFkZi00YjZkLTlhMjMtOWJmZTQ0NTIzNGY2IiwiTmFtZSI6ImJ1cmFra3VydSJ9.9flX6FhArFXQel8lRKqzW1jkGG2Fo43q36etZhI16Og

- Attributes: we should use our token when we send a request to this api to see sended messages and also toName.
QueryStringParemeters: ?toName=sampletoname&token=ourtoken

* Show User Received Messages from Someone:

EndPoint: http://localhost:8080/api/user/sended/messages?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE2MzEyODgyNTksIklkIjoiOTY3MzEwMmItZjFkZi00YjZkLTlhMjMtOWJmZTQ0NTIzNGY2IiwiTmFtZSI6ImJ1cmFra3VydSJ9.9flX6FhArFXQel8lRKqzW1jkGG2Fo43q36etZhI16Og

- Attributes: we should use our token when we send a request to this api to see received messages and also fromName.
QueryStringParemeters: ?fromName=sampletoname&token=ourtoken












