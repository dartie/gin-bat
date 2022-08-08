# gowebf

# TODO

- [ ] Registration/Request Access
- [ ] Embed Files

- [ ] Admin
    - [ ] Create user
    - [ ] Update user
    - [ ] Delete user
    - [ ] Reset user password
    - [ ] Create auth-token
    - [ ] Display auth-token

- [ ] Command line
    - [ ] Create user
    - [ ] Update user
    - [ ] Delete user
    - [ ] Reset user password
    - [ ] Create auth-token
    - [ ] Display auth-token
    - [ ] ---------------------
    - [ ] Project management
    - [ ] command line administration
    - [ ] Bootstrap/Vanilla javascript selection

- [ ] API
    - [ ] Create user
    - [ ] Update user
    - [ ] Delete user
    - [ ] Reset user password
    - [ ] Create auth-token


- [ ] CSRF


# Notes

## DB
NULL values are not handled, therefore the database has to be created using `NOT NULL` for all fields, to make sure no exceptions are raised in case of queries.


# Troubleshooting

## could not import github.com/mattn/go-sqlite3 (no required module provides package "github.com/mattn/go-sqlite3")
1. Download `mattn/go-sqlite` version `1.14.14`
    ```bash
    go get github.com/mattn/go-sqlite3@v1.14.14
    ```

    as mentioned in this thread : [issue-755](https://github.com/mattn/go-sqlite3/issues/755)

