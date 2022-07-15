# Authentication proof of concept using Go and Angular

![Diagram](/docs/images/diagram.png)

## Backend
Web api developed with GO using gorilla-mux library as router handler. 

Code coverage: 91.6%

The application has a control of all active sessions linked to their actual refresh token, allowing the user to revoke any active session. 

The login call returns 2 different tokes, **access token** and **refresh token**, both tokens have a different expiration time, normally the access token would be valid for 2 min and 1 year for the refresh token, these tokens are sent to the user through 2 new cookies. The **access token** is also in the body response in order to read the user data at frontend side.

All calls to authenticated endpoints require a valid access token cookie, the call will return an http error 401 (Unauthorized) if the access token is expired, the user would have to call the refresh endpoint in order to get a new access token.

An user can be administrator, this flag allows the user to delete any users, revoke any session and create a new user.

### Endpoints

#### /auth/login (POST)
Requires a valid basic authentication header as input, returns 2 new cookies with the access and refresh cookie along with the access token payload in the body response.

#### /auth/refresh (POST)
 Requires a valid refreshToken cookie, returns a new access token cookie along with the access token payload in the body response.

#### /users (GET)
 Requires a valid accessToken cookie, returns a list of all users

#### /users (POST)
 Requires a valid accessToken cookie, the user must be administrator, it needs the user data to be created
 ` NewUserInput
{
    "name": "user1",
    "password": "user1",
    "IsAdmin": false
}
 `

Returns a valid response if the user has been created

#### /users/{id} (DELETE)
Requires a valid accessToken cookie, the user must be administrator, a valid user id must be provided on the url. Returns a valid response if the user has been deleted

#### /sessions (GET)
Requires a valid accessToken cookie, a normal user can only see their own sessions, an administrator will get a list of all sessions in the application.

#### /sessions/{id} (DELETE)
Requires a valid accessToken cookie, a normal user can only delete their own sessions, an administrator can delete any active session. Returns a valid response if the session has been deleted


## Frontend
Work in progress :)