# Blog aggregator
## What you need 
You'll need a Postgresql database and Go installed to your machine to run this program.
To install the gator program simply run:
    go install

## Config file
Create a config file in your home dir
    ~/.gatorconfig.json
add 
    {"db_url":"postgres://postgres:@localhost:5432/gator?sslmode=disable"}

## Commands
login: logs in to a given user (e.g. blog-aggregator login marun)
register: creates a user with the given username and logs into it (e.g. blog-aggregator register marun)
users: displays all the users of the db (e.g. blog-aggregator users)
agg: creates posts with the feeds in db give a interval to go to next feed (e.g. blog-aggregator agg 20s)
addfeed: adds a feed and follows it with a given url (e.g. blog-aggregator addfeed example.com/rss)
feeds: displays all the feeds created (e.g. blog-aggregator feeds)
follow: current user follows a given feed url (e.g. blog-aggregator follow example.com/rss)
following: displays the feeds that the current user follows (e.g. blog-aggregator following)
unfollow: current user unfollows a given feed url (e.g. blog-aggregator unfollow example.com/rss)
browse: Browses through the posts of the current user, you can add a limit after command (e.g. blog-aggregator browse 10)
