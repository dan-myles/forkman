# [INSERT BOT NAME]

> WARNING: [INSERT BOT NAME] is still in early development, and is not yet
ready for use in production. BE CAREFUL!

[INSERT BOT NAME] is a self-hosted, open source Discord bot that is designed
for enterprise use in mind. As Discord becomes more of a social media platform
than an IRC clone, Discord bots have filled crucial needs in moderation,
analytics, security and more. The problem with today's Discord bots is that
if they get too big they end up offering NFT advertisements to your users.
This bot aims to be simple to understand, and easy to self host.

[INSERT BOT NAME] is *just* a Go binary that serves a REST API for a web dashboard,
and also opens a websocket connection to Discord to interact there. The authentication
for our web dashboard is hand rolled and only accepts Discord logins. The server uses
sqlite as its just the simplest database to understand. The web dashboard is just a
React SPA, that interacts with that REST API. The long term goal here is that you can
throw [INESRT BOT NAME] on a server of your choice and forget about it.

*Features*
- None yet lol

## Development

**Roadmap:**
- web dashboard
- client-server sync for realtime updates on commands to dashboard
- aggregate logs & open rest endpoint to collect them
- integrate sqlite
- setup case system for moderation actions
- extend base commands
- roll up a good ol' auth system (session based 🙅 no jwts)

**TODO:**
- pick web router (stdlib/echo/chi?)
- add graphite to repo for easy PRs


