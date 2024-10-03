# godmail
A go program to send mail with file attachment capabilities



# Usage

**NOTE: Please add a config file at ~/.config/godmail/config.json with password field**

```json
{
    "password" : "email_password"
}
`

`godmail -from=from_address -to=to_address -body=body -subject=subject -files="file1,file2,..."`
