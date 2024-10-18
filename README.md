# Project 2 for CMPS 3390

## Usage:

After cloning the repository to your target directory,

Set up environment path parameters and then run the project like this:

A runnable target directory is the folder that holds the .go file with 'func main()' and is also in 'package main'

```
go run <target-directory>
```

Alternatively, you can run the project using a .json configuration file:

```
go run <target-directory> <configuration-path>
```

* Reference resources/config.json
The address by default in config.json should be localhost:8080 if you chose to do this.
Set a header with the key 'X-API-KEY' with the reference key you defined in your environment or the configuration file.

* You can download an API platform such as www.postman.com to simulate API queries as well as the localhost agent to run it locally.

## Configuration:
### Variables:
Path is the file path to the SQLite database
Address is the address listened to by the router
Reference key is the key tested against when querying with 'X-API-KEY' in the header of a http request

### Setting environment variables:
Linux:
```
export <key>=<value>
```
Windows/Shell:
```
$env:<key> = <value>
```
