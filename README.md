# Simple banking app

## What it is:
It's just a simple app written in go to learn and understand the intricacies of using golang and how to correctly create transaccionts and create database connections

## Most inportant libraries used:
[1] SQLC -> library use to generate normaly hand made boilerplate go code for database usage using only SQL queries. It`s used via a CLI.

[2] go migrations package -> with help of this package we can seamlessly create database migrations and run the inside our database. 

[3] Standard libraries like fmt, sql for database connection and much others.  

## Ok! but, i want run the proyect
To run the proyect locally is recomended to use the makefile contained inside the same proyect,in order to use make file you need to have installed make on your machine. Feel free to find a way to add make on you machine Although there is not a version of the make package written por windows so in order to use make you should use WSL or similar applications. 

### Disclaimer:
This is not a proyect aimed to be a web server nor api, so the only command to make it work for the time beig is:

    make test