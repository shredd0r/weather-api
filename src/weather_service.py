import logging
from types import NoneType
from typing import List
import requests
from selectolax.parser import HTMLParser
import time

from models import CurrentWeatherForecast, DailyWeatherForecast, GetLocationSearchRequest, HourlyWeatherForecast, Location, UnitType, WeatherForecastRequest

# localization = standard 'BCP 47': uk-UA, 
# interval = today / hourbyhour / tenday, 
# city_id = 0b8697c01baca04214b4abd206319d3eba60db5fb05c191012c4e02f6bdb23a4, 
# unit = m / e / h
WEATHER_URL_FORMAT = "https://weather.com/{localization}/weather/{interval}/l/{city_id}?unit={unit}"

# unit type on the weather channel backend
METRIC_UNIT_TYPE = "m"
IMPERIA_UNIT_TYPE = "e"
HYBRID_UNIT_TYPE = "h"

SELECTOR_FOR_CURRENT_FEELS_LIKE_TEMPERATURE = 'span[class*="TodayDetailsCard--feelsLikeTempValue--"][data-testid="TemperatureValue"]'
SELECTOR_FOR_CURRENT_CURRENT_TEMPERATURE = 'span[class*="CurrentConditions--tempValue--"][data-testid="TemperatureValue"]'
SELECTOR_FOR_CURRENT_MAX_TEMPERATURE = 'div[class*="CurrentConditions--tempHiLoValue"] > span[data-testid="TemperatureValue"]:first-child'
SELECTOR_FOR_CURRENT_MIN_TEMPERATURE = 'div[class*="CurrentConditions--tempHiLoValue"] > span[data-testid="TemperatureValue"]:last-child'
SELECTOR_FOR_CURRENT_ICON_NAME = 'span[class*="Icon--iconWrapper"] > svg[class*="CurrentConditions--wxIcon"] > title'
SELECTOR_FOR_CURRENT_VISIBILITY = 'span[data-testid="VisibilityValue"] > span'

SELECTOR_FOR_HOURLY_FEELS_LIKE_TEMPERATURE = ''


SELECTOR_FOR_NOT_FOUND = 'div[class*="NotFound--notFound--"]'


class WeatherService:
    def __init__(self, logger: logging.Logger):
        self.logger = logger.getChild("weather-service")
        self.logger.setLevel(logging.DEBUG)

    r"""
        #TODO
    """
    def get_current_weather(self, request: WeatherForecastRequest) -> CurrentWeatherForecast:
        self.logger.info("start getting current weather forecast")
        url = WEATHER_URL_FORMAT.format(localization= request.localization, 
                                          interval= "today", 
                                          city_id= request.city_id,
                                          unit= self.__map_unit_type(request.unit_type))

        self.logger.debug(f"prepare request for get current weather page, url: {url}")
        current_weather_page = requests.get(url).text
        html_parser = HTMLParser(current_weather_page)

        if self.__page_is_not_found(html_parser):
            self.logger.debug("server return page with 'not_found' error")
            raise requests.RequestException("invalid request volumes")

        self.logger.debug("get html parser for current weather page")
        
        epoch_time = int(time.time())
        visibility = self.__get_current_visibility(html_parser)
        feels_like_temperature = self.__get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_FEELS_LIKE_TEMPERATURE)
        current_temperature = self.__get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_CURRENT_TEMPERATURE)
        min_temperature = self.__get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_MIN_TEMPERATURE)
        max_temperature = self.__get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_MAX_TEMPERATURE)
        icon_name = self.__get_current_icon_name(html_parser)
        
        self.logger.info("current weather forecast successful formed")
        return CurrentWeatherForecast(epoch_time= epoch_time,
                                      visibility= visibility,
                                      current_temperature= current_temperature,
                                      min_temperature= min_temperature, 
                                      max_temperature= max_temperature, 
                                      feels_like_temperature= feels_like_temperature, 
                                      icon_name= icon_name, 
                                      link= url)
    
    r"""
        #TODO
    """    
    def get_hourly_weather(self, request: WeatherForecastRequest) -> List[HourlyWeatherForecast]:
        pass

    r"""
        #TODO
    """
    def get_daily_weather(self, request: WeatherForecastRequest) -> List[DailyWeatherForecast]:
        pass

    r"""
        #TODO
    """
    def search_location_id(self, request: GetLocationSearchRequest) -> List[Location]:
        pass

    def __map_unit_type(self, unit_type: UnitType) -> str:
        if unit_type == UnitType.METRIC:
            return METRIC_UNIT_TYPE
        
        return IMPERIA_UNIT_TYPE
            

    def __page_is_not_found(self, html_parser: HTMLParser) -> bool:
        not_found_node = html_parser.css_first(SELECTOR_FOR_NOT_FOUND)
        return type(not_found_node) != NoneType

    def __get_temperature_by(self, html_parser: HTMLParser, selector: str):
        self.logger.debug(f"try return temperature by selector:{selector}")
        temperature_element = html_parser.css_first(selector)

        # Remove degree symbol from text
        temperature = temperature_element.text().replace("°", "")
        return temperature

    def __get_current_icon_name(self, html_parser: HTMLParser):
        return html_parser.css_first(SELECTOR_FOR_CURRENT_ICON_NAME).text()

    def __get_current_visibility(self, html_parser: HTMLParser):
        self.logger.debug("try get visibility node for current weather")
        visibility_node = html_parser.css_first(SELECTOR_FOR_CURRENT_VISIBILITY)

        if type(visibility_node) == NoneType:
            self.logger.debug("visibility is unlimited, return '0'")
            # Searching by selector SELECTOR_FOR_VISIBILITY return NoneType, because in the web, if visibility unlimited, display string `Unlimited`
            return 0
        
        # If visibility is limited, 'visibility_text' contains some float volume
        return visibility_node.text()
