# Fantasy Football Archive

Go-based RESTful API for accessing historical fantasy football data. The API lets you query past fantasy football league data, including league standings, team performance, drafts, and more.

Checkout the background and motivations for this projcect [here](https://litts.me/projects/2025/first/)

This data was retrieved from ESPN's API using [cwendt94's espn-api](https://github.com/cwendt94/espn-api) and stored in a PostgreSQL database. This wouldn't be possible without his amazing work!

### Prerequisites

- Go (version 1.23 or higher)
- PostgreSQL database with the exported ESPN data (see my other project for instructions on how to export the data, coming soon)