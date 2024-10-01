# envy 

This is a Work in progress!

Envy is a small app to manage environment variables and provide them to all your applications. A lightweight HTTP API allows for easy retrieval.

### Features

- Set and retrieve your variables from arbitrary scopes

Simple!

### Usage

```
Usage of envy.exe:
  -cfg string
        The configuration file to use (default "app.yaml")
  -logfile string
        The log file to write (default "envy.log")
```

If the configuration file isn't found, fallback values are used.

If the log file doesn't exist, it's created. If the log file is set to an empty string, no log file is written.

### Client implementations
- [Go: envyclient](https://github.com/KaiserWerk/envyclient)
