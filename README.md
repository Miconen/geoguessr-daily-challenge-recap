# geoguessr-daily-challenge-recap
Geoguessr API &amp; Discord integration for recapping daily challenges.

## Setup and running the project
I personally run this project on Railway and handle the CRON job through there. But you can use any CRON job manager, or run manually.
You will also need to aquire your `_ncfa` cookie from your browser storage once logged into GeoGuessr.

### Running using command flags
`go run main.go -discord {token} -geoguessr {ncfa} -users {comma separated list of discord ids}`

### Running using environment variables
First set these in your environment.
- `NCFA_TOKEN`
- `DISCORD_TOKEN`
- `DISCORD_USERS` (Comma separated list of Discord IDs)
- `CHALLENGE_DATE` (Optional: `today` (default), `yesterday` or `YYYY-MM-DD`)

Then simply just run `go run main.go`

### Recapping a past day
If a run failed or was skipped, you can re-send the recap for an earlier challenge:

`go run main.go -date yesterday` or `go run main.go -date 2026-09-06`

Past days are fetched from the club leaderboard endpoint, so only club members' results are included.

## Results
Once setup, the application will output your GeoGuessr club members current placements on the daily challenge.
<img width="553" height="308" alt="Discord_MJl3xW7kAV" src="https://github.com/user-attachments/assets/a23252f5-2aff-4c66-9cfd-844cae5afbfe" />
