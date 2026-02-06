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

Then simply just run `go run main.go`

## Results
Once setup, the application will output your GeoGuessr club members current placements on the daily challenge.
<img width="553" height="308" alt="Discord_MJl3xW7kAV" src="https://github.com/user-attachments/assets/a23252f5-2aff-4c66-9cfd-844cae5afbfe" />
