# Forkman

> WARNING: Forkman is still in early development, and is not yet
ready for use in production. BE CAREFUL!

Forkman is a self-hosted, open source Discord bot that is designed
for enterprise use in mind. As Discord becomes more of a social media platform
than an IRC clone, Discord bots have filled crucial needs in moderation,
analytics, security and more. The problem with today's Discord bots is that
if they get too big they end up offering NFT advertisements to your users.
This bot aims to be simple to understand, and easy to self host.

Forkman is *just* a Go binary that serves a REST API for a web dashboard,
and also opens a websocket connection to Discord to interact there. The authentication
for our web dashboard is hand rolled and only accepts Discord logins. The server uses
sqlite as its just the simplest database to understand. The web dashboard is just a
React SPA, that interacts with that REST API. The long term goal here is that you can
throw Forkman on a server of your choice and forget about it.

*Features*
- None yet lol

## Development

**Roadmap:**
- web dashboard
- client-server sync for realtime updates on commands to dashboard
- aggregate logs & open rest endpoint to collect them
- setup case system for moderation actions
- extend base commands
- roll up a good ol' auth system (session based 🙅 no jwts)
- backup system (with some sort of zip file like pocketbase)

**TODO:**
- [x] pick web router (stdlib/echo/chi?)
- [ ] setup dev script for hot-reloading vite & go
- [ ] add graphite to repo for easy PRs
- [ ] move from config.json to env variables
- [ ] store module state separately from env vars (prob still json?)
- [ ] move to github issues instead of this todo-list lol
- [ ] pick db (prob just sqlite!)
- [ ] write a dockerfile that binds a volume to preserve DB state (need to store all state in some folder)

## Deps
- go
- air
- pnpm

