# Services

## Auth

- Register(Username, Email, Password String)
  - Hash User Password using bycrypt
    - Return Err_Password_Encryption error if hashing is failed
  - Add User to database
    - Return Err_Username_Taken error if the Username has taken
    - Return Err_Email_Taken error if the Email has taken
    - Return error if it failed
  (I can consider a success process if it reached this step)
  - Register Username and Email using AuthService.Credential service
    - Return Err_Credential_Registration error if it failed
  - Add the created User to the cache
  - Generate JWT token
    - Return Err_JWT_Sign error if it failed
  - Return JWT token

- Login(Username, Email, Password, IP String)
  - Get login_tries from redis
    - Return Err_Too_Many_Attempts error if login_tries is more than 3
    - Return error if it failed
  - Get User by Username
    - Return Err_Not_Found error if it is not found
    - Return error if it failed
  - Compare the User hash with User password
    - If Password is wrong
      - Increment login_tries by one
      - Return Err_Invalid_Password error
    - Return error if it failed
  - Generate JWT token
    - Return Err_JWT_Sign error if it failed
  - Add the fetched User to the cache
  - Remove login_tries from redis
    - Return error if it failed
  - Return JWT token

- Credential(Username, Email String)
  - Pipeline redis
    - Add username:"username" key to redis
    - Add email:"<email@email.com>" key to redis
  - Return Err_Credential_Registration error if pipeline failed

- Authorize(User_ID UUID) (I can use UserService.GetByID)
  - Get User by ID from the cache
    - Return User if it is found
  - Get User by ID from the database
    - Return Err_Not_Found error if it is not found
    - Return error if it failed
  - Add the fetched User to the cache
  - Return the fetched User

## User

- Get By ID(ID UUID)
  - Get User by ID from the cache
    - Return the User if it is found
  - Get User by ID from the database
    - Return Err_Not_Found error if the user is not found
    - Return error if it failed
  - Add the User to the cache
  - Return User

- Get Profile By ID(ID UUID)
  - Get User Profile from the cache
    - Return the User Profile if it is found
  - Join User and User Profile tables
  - Get User Profile with username by User ID
    - Return Err_Not_Found error if the user is not found
    - Return error if it failed
  - Add The User Profile to the cache
  - Return User Profile

- List Profile By ID(IDs UUID[])
  - Get User Profiles from the cache
    - Return List if all User Profiles are found
  - Collect User Profile User_IDs were not found
  - Join User and User Profile tables
  - Get User profiles with username by User IDs
    - Return Err_Not_Found error if they are not found
    - Return error if it failed
  - Add The User Profiles from the database to the cache
  - Merge User Profiles
  - Return User Profiles

- Has Username(Username String) Boolean
  - Check Username in the cache
    - If it is not found
      - Check Username in database
        - If it is found
          - Add Username to the cache
          - Return true
        - Return false if it is not found
        - Return error if it failed
    - Return error if it failed
    - Return true if it is found

- Has Email(Email String) Boolean
  - Check Email in the cache
    - If it is not found
      - Check Email in database
        - If it is found
          - Add Email to cache
          - Return true
        - Return false if it is not found
        - Return error if it failed
    - Return error if it failed
    - Return true if it is found

- Update Password(User_ID UUID, Password, New_Password String)
  - Get User from database by User ID
    - Return Err_Not_Found error if it is not found
    - Return error if it failed
  - Compare User's hashed to password
    - Return Err_Invalid_Password error if the User's password is wrong
  - Hash new password using bycrypt
    - Return Err_Password_Encryption error if hashing is failed
  - Update User's password
    - Return error if it failed
  - Remove the User from the cache

- Update Profile(User_ID UUID, Name, Status, Logo String)
  - Update User's profile (Need to do it in one query)
    - Return Err_Not_Found error if it is not found
    - Return error if it failed
  - Get joined User's profile from database
    - Return Err_Not_Found error if it is not found
    - Return error if it failed
  - Remove The User from the cache
  - Return updated user's profile
