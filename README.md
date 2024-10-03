# godmail
A go program to send mail with file attachment capabilities

# Configuration


Configuration file stores the default values for email, password, subject. These fields are accessed when no flag for the corresponding field is passed to `godmail`.

Config should be placed at `~/.config/godmail/config.json`

```json
{
    "email" : "email_address", // If this is not provided, then address has to be provided with the -from flag
    "password" : "email_password" // Same for this
}
```

# Usage

* `godmail -from=from_address -to=to_address -body=body -subject=subject -files="file1,file2,..." -password="My not so super secure password"`

* `godmail -to=someone@example.com -body="This is a body" -subject="Exciting Subject" -files="~/example.pdf"`

In this case, `from address` and `password` is used from the config file (if it exists).


# TODOs

- [x] Attach text
- [x] Attach files
- [ ] Delayed sending

