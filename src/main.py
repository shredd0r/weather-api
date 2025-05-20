import logging
from typing import List
from fastapi import FastAPI
from models import CurrentWeatherForecast, DailyWeatherForecast, GetLocationSearchRequest, HourlyWeatherForecast, Location, UnitType, WeatherForecastRequest
from weather_service import WeatherService


app = FastAPI()
logger = logging.getLogger("uvicorn")
weather_service = WeatherService(logger)


@app.get("/weather/current/{city_id}")
def current_weather(city_id: str, localization: str, unit_type: int) -> CurrentWeatherForecast:
    request = WeatherForecastRequest(city_id= city_id, localization= localization, unit_type= unit_type)
    return weather_service.get_current_weather(request)


@app.get("/weather/hourly")
def hourly_weather(request: WeatherForecastRequest) -> List[HourlyWeatherForecast]:
    pass

@app.get("/weather/daily")
def daily_weather(request: WeatherForecastRequest) -> List[DailyWeatherForecast]:
    pass

@app.get("/location/search")
def search_location_id(request: GetLocationSearchRequest) -> List[Location]:
    pass