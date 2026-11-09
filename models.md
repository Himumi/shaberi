# Models

## Postgresql

### User

- ID: UUID
- Username: String Unique
- Email: String
- Hash: String
- Created At: Time Stamp
- Updated At: Time Stamp
- Deleted At: Time Stamp

### Profile

- ID: Integer
- Name: String
- Status: String
- Logo: Link String
- User ID: UUID
- Created At: Time Stamp
- Updated At: Time Stamp
- Deleted At: Time Stamp

### User Profile Response

- Profile
- Username: String

## Redis

### Username (SISMEMBER)

- "users:taken_usernames": "actual_usernames", etc

### Email (SISMEMBER)

- "users:taken_emails": "actual_emails", etc

### Increment Login Tries (Key-Value)

- "auth:login:input:username_or_email_in_lowercase:ip:target_ip_address": Number (Max 3) (24 hours expiration)

### Remove Login Tries

- "auth:login:input:username_or_email_in_lowercase:ip:target_ip_address"

### Lock Login
