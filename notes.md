# Let's Go

## Structure
- The `cmd` directory will contain the application-specific code for the executable applications in the project.
- The `internal` directory will contain the ancillary non-application specific code used in the project, We'll use it to hold potentially reusable code like validation helpers and SQL database models for the project.
    packages under internal cannot be imported by code outside of our project.
- The ui directory will contain the user-interface assets used by the web application. Specifically, the `ui/html` contains HTML template and ui/static directory will contain statis files(like css, and images)
