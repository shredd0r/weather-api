# Weather API

### Requirements

- docker engine

This app get weather from open source and store the information about weather to cache.
Api realized by Graphql.

There are 4 api`s:
- findGeocoding
- currentWeather 
- hourlyWeather 
- dailyWeather

Example request current weather:
```graphql
query {
    currentWeather(input: {
        coords: {
            latitude: 50.026501
            longitude: 36.239391
        }
        locale: "uk-ua"
        unit: "metric"
        forecaster: "AccuWeather"
    }) {
        epochTime
        visibility
        currentTemperature
        minTemperature
        maxTemperature
        feelsLikeTemperature
        iconId
        mobileLink
        link
    }
}
```
```json
{
  "data": {
    "currentWeather": {
      "epochTime": 1722781620,
      "visibility": 24.1,
      "currentTemperature": 32.6,
      "minTemperature": null,
      "MaxTemperature": null,
      "FeelsLikeTemperature": 31.1,
      "weatherId": "01",
      "mobileLink": "http://www.accuweather.com/uk/tr/yagmur/1284006/current-weather/1284006",
      "Link": "http://www.accuweather.com/uk/tr/yagmur/1284006/current-weather/1284006"
    }
  }
}
```

### Start

1. Create file .env in app directory with variables from dev.env
2. In variables with api keys set your keys
3. run command in app directory: ``docker-compose up -d``

Component have sandbox where you can create requests for all apies. If you want enable this sandbox, add in ``.env`` file:

- SERVER_PORT=8080
- SERVER_PLAYGROUND_ENABLE=true

After start docker compose, follow the link ``localhost:8080`` 

### Features
- add another forecaster