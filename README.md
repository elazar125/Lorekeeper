# Lorekeeper

A bot for querying card data from https://spellslingerer.com

## To run yourself

You'll need to modify a few things:
1) Copy .env.sample to .env and supply your own bot key
2) Modify bd.sh to point to your own SSH configuration
3) Create a systemd unit for your and use it to run your app (refer to the correct name in bd.sh)

## Sample systemd config

I use this, modified from the [pocketbase.io docs](https://pocketbase.io/docs/going-to-production/)

```systemd
[Unit]
Description = Lorekeeper Discord bot

[Service]
Type           = simple
User           = root
Group          = root
LimitNOFILE    = 4096
Restart        = always
RestartSec     = 5s
StandardOutput = append:/root/lorekeeper/logs/out.log
StandardError  = append:/root/lorekeeper/logs/errors.log
WorkingDirectory = /root/lorekeeper/
ExecStart      = /root/lorekeeper/lorekeeper

[Install]
WantedBy = multi-user.target
```
