# Sport-Companion
A CLI to help sport enthusiasts stay up to date with scores and players from their favourite teams!

To use this CLI, make a free account on [api-sports.io](https://api-sports.io/) and replace XXXXX in the `req.Header.Add("x-rapidapi-key", "XXXXX")` line of each file in the `/internal/api` folder with your own personal API key!

The free account comes with 100 requests per day for each sport. This means that if I make a request for all Premier League teams then all Serie A teams, that would count as 2 requests.
