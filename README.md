# 👀 togglwatcher

## Usage

Create `config.yaml`

```yaml
timezone: Asia/Tokyo
interval: 10
api:
  token: <toggl API token>
  dashboardid: <Dashboard ID>
notifications:
  - service: slack
    webhook_url: https://hooks.slack.com/services/xxx
    channel: general
users:
  - id: 8901234
    name: kimromi
```

Execute command

```
$ togglwatcher
```
