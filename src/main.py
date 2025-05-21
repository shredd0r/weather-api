import logging
from typing import List
from fastapi import FastAPI
from models import CurrentWeatherForecast, DailyWeatherForecast, GetLocationSearchRequest, HourlyWeatherForecast, Location, UnitType, WeatherForecastRequest
from weather_service import WeatherService


app = FastAPI()
logger = logging.getLogger("uvicorn")
weather_service = WeatherService(logger)


@app.get("/weather/current/{place_id}")
def current_weather(place_id: str, localization: str, unit_type: int) -> CurrentWeatherForecast:
    request = WeatherForecastRequest(place_id= place_id, localization= localization, unit_type= unit_type)
    return weather_service.get_current_weather(request)


@app.get("/weather/hourly/{place_id}")
def hourly_weather(place_id: str, localization: str, unit_type: int) -> List[HourlyWeatherForecast]:
    request = WeatherForecastRequest(place_id= place_id, localization= localization, unit_type= unit_type)
    return weather_service.get_hourly_weather(request)

@app.get("/weather/daily/{city_id}")
def daily_weather(place_id: str, localization: str, unit_type: int) -> List[DailyWeatherForecast]:
    request = WeatherForecastRequest(place_id= place_id, localization= localization, unit_type= unit_type)
    return weather_service.get_daily_weather(request)

# TODO
@app.get("/location/search")
def search_location_id(request: GetLocationSearchRequest) -> List[Location]:
    pass