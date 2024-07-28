# Sport-Companion
A CLI to help sports enthusiasts stay up to date with upcoming games for their favourite teams from the MLB, NBA, NFL and NHL!

# Setup
Firstly, go to [API-SPORTS](https://api-sports.io) and create a free account. A free account allows up to 100 API calls per day per sport.

After that, you will want to create a .env file and store your API key under "SPORT_API_KEY".\
```
SPORT_API_KEY=XxXxXxXxXxXxXxXx
```

If you don't have Go installed then there is an executable for this CLI called `Sport-Companion` in the repo.\
Otherwise, you can build the application using:\
```
go build -o <exec name> github.com/tremerj/Sport-Companion
```

If you want to use this command anywhere in your file system, make sure to add it to your PATH.

# Usage
There is a help command that will give you general instructions on how to use this CLI.
**Reminder to always put team names in double quotes!!!**
eg.
```
Sport-Companion add NFL "Miami Dolphins"
Sport-Companion remove MLB "New York Yankees"
```

### Adding teams
To add teams to your watchlist run this command:\
```
Sport-Companion add <league> <team>
```

### Removing teams
To remove teams from your watchlist run this command:\
```
Sport-Companion remove <league> <team>
```

### Clearing watchlist
To clear your entire watchlist run this command:\
```
Sport-Companion clear
```

### List watchlist
To list all the teams currently in your watchlist run this command:\
```
Sport-Companion list
```

### See schedule
To see the next 5 games of all of your teams in your watchlist run this command:\
```
Sport-Companion schedule
```

# Thanks
Thank you very much for using this CLI if there are any bugs please report them!!! :)


