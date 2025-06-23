import logging
from typing import List
from fastapi import FastAPI
from .models import CurrentWeatherForecast, DailyWeatherForecast, GetLocationSearchRequest, HourlyWeatherForecast, Location, UnitType, WeatherForecastRequest
from .weather_client import WeatherClient
from .weather_service import WeatherService


app = FastAPI()
logger = logging.getLogger("uvicorn")
logger.setLevel(logging.DEBUG) # TODO add configuration for logger

weather_client = WeatherClient(logger)
weather_service = WeatherService(logger, weather_client)


@app.get("/weather/current/{place_id}")
def current_weather(place_id: str, language: str, unit_type: int) -> CurrentWeatherForecast:
    r"""
        Method: GET,
        Routing: /weather/current/{place_id}?language={language}&unit_type={unit_type}
        
        REST API for get current weather forecast by place id.
        Data is obtained by parsing html page using selectolax.HTMLParser from site The Weather Channel.

        Request:
            - place_id    - required string thats describe place. See api /location/search
            - localizaion - localization string in standard 'BCP 47'. For example: 'uk-UA'
            - unit_type   - type of unit format. See enum UnitType on models.py
        
        Response:
            - epoch_time:             int
            - visibility:             float
            - current_temperature:    int 
            - min_temperature:        int 
            - max_temperature:        int
            - feels_like_temperature: int
            - icon_name:              str 
            - link:                   str

    """

    request = WeatherForecastRequest(place_id= place_id, language= language, unit_type= unit_type)
    return weather_service.get_current_weather(request)


@app.get("/weather/hourly/{place_id}")
def hourly_weather(place_id: str, language: str, unit_type: int) -> List[HourlyWeatherForecast]:
    r"""
        Method: GET,
        Routing: /weather/hourly/{place_id}?language={language}&unit_type={unit_type}
        
        REST API for get hourly weather forecast by place id.
        Data is obtained by parsing html page using selectolax.HTMLParser from site The Weather Channel.

        Request:
            - place_id    - required string thats describe place. See api /location/search
            - language    - language string in standard 'BCP 47'. For example: 'uk-UA'
            - unit_type   - type of unit format. See enum UnitType on models.py
        
        Response:
        List of models:
            - epoch_time:                    int
            - current_temperature:           int
            - feels_like_temperature:        int
            - uv_index:                      float
            - probability_of_precipitation:  float
            - precipitation_type:            PrecipitationType (see in models.py)
            - amount_of_precipitation:       float
            - wind:                          Wind
            - icon:                          str
            - link:                          str
    """
    request = WeatherForecastRequest(place_id= place_id, language= language, unit_type= unit_type)
    return weather_service.get_hourly_weather(request)

@app.get("/weather/daily/{place_id}")
def daily_weather(place_id: str, language: str, unit_type: int) -> List[DailyWeatherForecast]:
    r"""
        Method: GET,
        Routing: /weather/daily/{place_id}?language={language}&unit_type={unit_type}
        
        REST API for get daily weather forecast by place id.
        Data is obtained by parsing html page using selectolax.HTMLParser from site The Weather Channel.

        Request:
            - place_id    - required string thats describe place. See api /location/search
            - localizaion - required localization string in standard 'BCP 47'. For example: 'uk-UA'
            - unit_type   - required type of unit format. See enum UnitType on models.py
        
        Response:
        List of models:
            - epoch_time:  int
            - day:         DailyWeatherForecastDetail | None (can be null if current day period is night)
            - night:       DailyWeatherForecastDetail
            - link:        str

            DailyWeatherForecastDetail:
                - temperature:                   int
                - humidity:                      float
                - wind:                          Wind
                - rise_time:                     int
                - set_time:                      int
                - probability_of_precipitation:  float
                - precipitation_type:            PrecipitationType
                - icon_name:                     str
    """
    request = WeatherForecastRequest(place_id= place_id, language= language, unit_type= unit_type)
    return weather_service.get_daily_weather(request)

@app.get("/location/search")
def search_location_id(language: str, place_details: str) -> List[Location]:
    r"""
        Method: GET,
        Routing: /location/search?language={language}&place_details={place_details}
        
        REST API for search your place by some details like city name, postcode etc. The api will return similar places by details. 

        Request:
            - localization: str - localization string in standard 'BCP 47'. For example: 'uk-UA'
            - place_details: str - string with detail information about your place, for example: city name, postcode, address etc

        Response:
        List of models:
            - place_id:     str
            - address:      str
            - city:         str
            - country:      str
            - latitude:     float
            - longitude:    float
            - postal_code:  str
    """
    request = GetLocationSearchRequest(localization=language, place_details=place_details)
    return weather_service.search_location_id(request)