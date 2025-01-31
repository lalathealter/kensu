# Running instructions

## Important preparations

For the purposes of this assignment I decided to use the PostgreSQL DBMS; Make sure you have it installed before proceeding.

To create a database you can use the PSQL tool:
```
psql -U postgres -d postgres -c "CREATE DATABASE $dbname;"
```

After doing so, you are to setup its structure using the `kensu.sql` dump file provided in `server` directory:
```
psql -U postgres $dbname < kensu.sql
```

On completing that please create an .env file containing the following fields and place it in the `server` directory:
```
host = address for the server app
port = port for the server app
dbport = port for the database (Postgres uses 5432 by default)
dbname = $dbname that you entered previously
dbuser = the preffered login for Postgres
dbpass = the password for that login
```

## Applications

**Please make sure that you have a golang compiler installed**

It is implied that you already cd'ed to go/carrier-pricing directory;

Both server app and CLI-client belong to their corresponding directories; You are not supposed to run or test a CLI-client before launching a server app.

To download the dependencies (a PostgreSQL driver and a .env utility):
```
go mod download
```

Command to run:
```
go run .
```

To build a binary:
```
go build .
```

To test:
```
go test ./...
```

***For a detailed assignment description address carrier-pricing.md***


# Take-Home Coding Exercises - General Overview

## What to expect?
We expect that the amount of effort to do any of these exercises is in the range of about 4-6 hours of actual work. 
We also understand that your time is valuable, and in anyone's busy schedule that constitutes a fairly substantial chunk of time, so we really appreciate any effort you put in to helping us build a solid team.

## What we are looking for?
**Keep it simple**. Really. 4-6 hours isn't a lot of time, and we really don't want you spending too much more time on it than that.

**Treat it like production code**. That is, develop your software in the same way that you would for any code that is intended 
to be deployed to production. These may be toy exercises, but we really would like to get an idea of how you build code on a day-to-day basis.

## How to submit?
You can do this however you see fit - you can email us a tarball, a pointer to download your code from somewhere or just a link to a source control repository.
Make sure your submission includes a **README**, documenting any assumptions, simplifications and/or choices you made, 
as well as a short description of how to run the code and/or tests. Finally, to help us review your code, 
please split your commit history in sensible chunks (at least separate the initial provided code from your personal additions).

## GO role exercises

### [Backend Developer] Built an API for managing prices
The complete specification for this exercise can be found in the [carrier-data.md](go/carrier-pricing/carrier-data.md).

## Scala role exercises

### [Junior Backend Developer] Built Json Transformer API
The complete specification for this exercise can be found in the [junior-challenge.md](scala/junior/junior-challenge.md).

### [Backend Developer] Build an API for managing users
The complete specification for this exercise can be found in the [UsersAPI.md](scala/users/UsersAPI.md).

### [Senior Backend Developer] Build a local proxy for currency exchange rates

The complete specification for this exercise can be found in the [Forex.md](scala/forex/Forex.md).

## F.A.Q.
1) _Is it OK to share your solutions publicly?_
Yes, the questions are not prescriptive, the process and discussion around the code is the valuable part. 
You do the work, you own the code. Given we are asking you to give up your time, it is entirely reasonable for you to keep and use your solution as you see fit.

2) _Should I do X?_
For any value of X, it is up to you, we intentionally leave the problem a little open-ended and will leave it up to you
 to provide us with what you see as important. Just remember the rough time frame of the project. 
 If it is going to take you a couple of days, it isn't essential.

3) _Something is ambiguous, and I don't know what to do?_
The first thing is: don't get stuck. We really don't want to trip you up intentionally, we are just attempting to see
 how you approach problems. That said, there are intentional ambiguities in the specifications, mainly to 
 see how you fill in those gaps, and how you make design choices. If you really feel stuck, our first preference is 
 for you to make a decision and document it with your submission - in this case there is really no wrong answer. 
 If you feel it is not possible to do this, just send us an email and we will try to clarify or correct the question for you.

Cheers!

## Kensu Project Details

### BUSINESS CRITICALITY

low

### ENVIRONMENT

internal

### LIFECYCLE STAGE

development

### PROJECT OWNER

michele-pinto-kensu
