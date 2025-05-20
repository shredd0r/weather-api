from enum import Enum
from typing import List

from pydantic import BaseModel

class PrecipitationType(Enum):
    RAIN = "rain"
    SNOW = "snow"
    ICE = "ice"

class UnitType(Enum):
    METRIC = 1
    IMPERIAL = 2

class Wind(BaseModel):
    direction: str
    speed: float

class CurrentWeatherForecast(BaseModel):
    epoch_time: int
    visibility: float
    current_temperature: int 
    min_temperature: int 
    max_temperature: int
    feels_like_temperature: int
    icon_name: str 
    link: str

class HourlyWeatherForecast(BaseModel):
    epoch_time: int
    current_temperature: int
    feels_like_temperature: int
    uv_index: float
    probability_of_precipitation: float
    precipitation_type: PrecipitationType
    amount_of_precipitation: float
    wind: Wind
    icon: str
    link: str
        
class DailyWeatherForecast(BaseModel):
    epoch_time: int
    min_temperature: float
    max_temperature: float
    humidity: float
    uv_index: float
    sunrise_time: int
    sunset_time: int
    wind: Wind
    probability_of_precipitation: float
    precipitation_type: PrecipitationType
    icon_id: str
    link: str

class WeatherForecastRequest(BaseModel):
    city_id: str
    localization: str
    unit_type: UnitType

class GetLocationSearchRequest(BaseModel):
    language: str
    city_address: str
        
class Location(BaseModel):
    place_id: str
    address: str
    city: str
    country: str
    latitude: float
    longitude: float
    postal_code: str

class GetSunV3LocationSearchRequest(BaseModel):
    language: str
    place_detail: str

class GetSunV3LocationSearchResponse(BaseModel):
    address: str
    admin_district: str
    city: str
    country: str
    country_code: str
    display_name: str
    iana_time_zone: str
    latitude: float
    longitude: float
    place_id: str
    postal_code: str

    def _repr__(self):
        return f"""address={self.address}, admin_district={self.admin_district}, city={self.city}, country={self.country}, country_code={self.country_code}, 
        display_name={self.display_name}, iana_time_zone={self.iana_time_zone}, latitude={self.latitude}, longitude={self.longitude}, place_id={self.place_id}, postal_code={self.postal_code}"""

