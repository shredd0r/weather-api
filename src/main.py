import logging
from typing import List
from fastapi import FastAPI
from models import CurrentWeatherForecast, DailyWeatherForecast, GetLocationSearchRequest, HourlyWeatherForecast, Location, UnitType, WeatherForecastRequest
from weather_client import WeatherClient
from weather_service import WeatherService


app = FastAPI()
logger = logging.getLogger("uvicorn")
logger.setLevel(logging.DEBUG) # TODO add configuration for logger

weather_client = WeatherClient(logger)
weather_service = WeatherService(logger, weather_client)


@app.get("/weather/current/{place_id}")
def current_weather(place_id: str, localization: str, unit_type: int) -> CurrentWeatherForecast:
    request = WeatherForecastRequest(place_id= place_id, localization= localization, unit_type= unit_type)
    return weather_service.get_current_weather(request)


@app.get("/weather/hourly/{place_id}")
def hourly_weather(place_id: str, localization: str, unit_type: int) -> List[HourlyWeatherForecast]:
    request = WeatherForecastRequest(place_id= place_id, localization= localization, unit_type= unit_type)
    return weather_service.get_hourly_weather(request)

@app.get("/weather/daily/{place_id}")
def daily_weather(place_id: str, localization: str, unit_type: int) -> List[DailyWeatherForecast]:
    request = WeatherForecastRequest(place_id= place_id, localization= localization, unit_type= unit_type)
    return weather_service.get_daily_weather(request)

@app.get("/location/search")
def search_location_id(language: str, place_details: str) -> List[Location]:
    request = GetLocationSearchRequest(localization=language, place_details=place_details)
    return weather_service.search_location_id(request)