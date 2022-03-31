slack-[pensieve](https://en.wikipedia.org/wiki/Magical_objects_in_Harry_Potter#Pensieve)

This is a hastily thrown together program to pull down all DMs that you don't get from Slack's official export tool.

## To use it

1. Create a slack app
2. Grant the following _user_ oAuth scopes
- `channels:history`
- `channels:read`
- `groups:history`
- `groups:read`
- `im:history`
- `im:read`
- `mpim:history`
- `mpim:read`
- `users:read`
3. Populate `.env` with values for the two variables defined in `.env.sample`
