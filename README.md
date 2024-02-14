split/bind packages taken from [Jesse Duffield's horcrux](https://github.com/jesseduffield/horcrux), which adapts shamir code from [Hashicorp's Vault](https://github.com/hashicorp/vault)

requires a `secrets.json` file with a discord bot `TOKEN` and `CHANNEL_ID`:

```
{
    "TOKEN": "token123",
    "CHANNEL_ID": "123"
}
```

currently in a minimal working state but needs improvement  

to-do:  
- lots of general refactoring
- handle many more input cases/errors
- tests/pipeline
- maybe a UI
