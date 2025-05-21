import re
import logging
from types import NoneType
from typing import List
import requests
from selectolax.parser import HTMLParser
import time

from models import CurrentWeatherForecast, DailyWeatherForecast, GetLocationSearchRequest, HourlyWeatherForecast, Location, PrecipitationType, UnitType, WeatherForecastRequest, Wind

# localization = standard 'BCP 47': uk-UA, 
# interval = today / hourbyhour / tenday, 
# place_id = 0b8697c01baca04214b4abd206319d3eba60db5fb05c191012c4e02f6bdb23a4, 
# unit = m / e / h
WEATHER_URL_FORMAT = "https://weather.com/{localization}/weather/{interval}/l/{place_id}?unit={unit}"

# unit type on the weather channel backend
METRIC_UNIT_TYPE = "m"
IMPERIA_UNIT_TYPE = "e"
HYBRID_UNIT_TYPE = "h"

CURRENT_FORECAST_INTERVAL = "today"
HOULY_FORECAST_INTERVAL = "hourbyhour"
DAILY_FORECAST_INTERVAL = "tenday"

SELECTOR_FOR_CURRENT_FEELS_LIKE_TEMPERATURE = 'span[class*="TodayDetailsCard--feelsLikeTempValue--"][data-testid="TemperatureValue"]'
SELECTOR_FOR_CURRENT_CURRENT_TEMPERATURE = 'span[class*="CurrentConditions--tempValue--"][data-testid="TemperatureValue"]'
SELECTOR_FOR_CURRENT_MAX_TEMPERATURE = 'div[class*="CurrentConditions--tempHiLoValue"] > span[data-testid="TemperatureValue"]:first-child'
SELECTOR_FOR_CURRENT_MIN_TEMPERATURE = 'div[class*="CurrentConditions--tempHiLoValue"] > span[data-testid="TemperatureValue"]:last-child'
SELECTOR_FOR_CURRENT_ICON_NAME = 'span[class*="Icon--iconWrapper"] > svg[class*="CurrentConditions--wxIcon"] > title'
SELECTOR_FOR_CURRENT_VISIBILITY = 'span[data-testid="VisibilityValue"] > span'

SELECTOR_FOR_HOURLY_BLOCK_SUMMARY_FORECAST = 'div[class*="DetailsSummary--DetailsSummary--"][data-testid="DetailsSummary"]' #
SELECTOR_FOR_HOURLY_BLOCK_DETAIL_TABLE_FORECAST = 'ul[class*=DetailsTable--DetailsTable--]'
SELECTOR_FOR_HOURLY_CURRENT_TEMPERATURE = 'span[class*="DetailsSummary--tempValue--"]'
SELECTOR_FOR_HOURLY_FEELS_LIKE_TEMPERATURE = 'span[class*=DetailsTable--value--][data-testid="TemperatureValue"]'
SELECTOR_FOR_HOURLY_UV_INDEX = 'span[class*=DetailsTable--value--][data-testid="UVIndexValue"]'
SELECTOR_FOR_HOURLY_PROBABILITY_OF_PRECIPITATION = 'div[class*="DetailsSummary--precip--"] > span[data-testid="PercentageValue"]'
SELECTOR_FOR_HOURLY_AMOUNT_OF_PRECIPITATION = 'li[class*="DetailsTable--listItem--"][data-testid="AccumulationSection"] > div > span > span:first-child'
SELECTOR_FOR_HOURLY_WIND = 'span[class*="Wind--windWrapper--"][data-testid="Wind"] > span'
SELECTOR_FOR_HOURLY_ICON_NAME = 'svg[class*="DetailsSummary--wxIcon--"] > title'

SELECTOR_FOR_NOT_FOUND = 'div[class*="NotFound--notFound--"]'


class WeatherService:
    r"""
        WeatherService has methods for return weather forecast by interval or search placeId using apies 'The Weather Channel'
    """

    def __init__(self, logger: logging.Logger):
        self.logger = logger.getChild("weather-service")
        self.logger.setLevel(logging.DEBUG)

    def get_current_weather(self, request: WeatherForecastRequest) -> CurrentWeatherForecast:
        r"""
            This method return current weather forecast from site 'The Weather Channel'.
            Data is obtained by parsing html page using selectolax.HTMLParser

            Request:
                WeatherForecastRequest:
                    place_id: str - id which you can get from WeatherClient.get_sun_location_search()
                    localization: str - localization string in standard 'BCP 47'. For example: 'uk-UA'
                    unit_type: UnitType
        """

        self.logger.info("start getting current weather forecast")
        url = self._get_url_for_get_weather_page(CURRENT_FORECAST_INTERVAL, request)
        html_parser = self._get_html_parser_for_weather_page(url)
        
        epoch_time = int(time.time())
        visibility = self._get_current_visibility(html_parser)
        feels_like_temperature = self._get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_FEELS_LIKE_TEMPERATURE)
        current_temperature = self._get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_CURRENT_TEMPERATURE)
        min_temperature = self._get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_MIN_TEMPERATURE)
        max_temperature = self._get_temperature_by(html_parser, SELECTOR_FOR_CURRENT_MAX_TEMPERATURE)
        icon_name = self._get_current_icon_name(html_parser)
        
        self.logger.info("current weather forecast successful formed")
        return CurrentWeatherForecast(epoch_time= epoch_time,
                                      visibility= visibility,
                                      current_temperature= current_temperature,
                                      min_temperature= min_temperature, 
                                      max_temperature= max_temperature, 
                                      feels_like_temperature= feels_like_temperature, 
                                      icon_name= icon_name, 
                                      link= url)
       
    def get_hourly_weather(self, request: WeatherForecastRequest) -> List[HourlyWeatherForecast]:
        r"""
            This method return hourly weather forecast from site 'The Weather Channel'.
            Data is obtained by parsing html page using selectolax.HTMLParser

            Request:
                WeatherForecastRequest:
                    place_id: str - id which you can get from WeatherClient.get_sun_location_search()
                    localization: str - localization string in standard 'BCP 47'. For example: 'uk-UA'
                    unit_type: UnitType
        """

        self.logger.debug("start getting hourly weather forecast")
        url = self._get_url_for_get_weather_page(HOULY_FORECAST_INTERVAL, request)
        html_parser = self._get_html_parser_for_weather_page(url)

        hourly_weather_forecasts = []
        hourly_summary_nodes = html_parser.css(SELECTOR_FOR_HOURLY_BLOCK_SUMMARY_FORECAST)
        hourly_detail_nodes = html_parser.css(SELECTOR_FOR_HOURLY_BLOCK_DETAIL_TABLE_FORECAST)

        for i in range(len(hourly_summary_nodes)):
            epoch_time=0 # TODO
            current_temperature = self._get_temperature_by(hourly_summary_nodes[i], SELECTOR_FOR_HOURLY_CURRENT_TEMPERATURE)
            feels_like_temperature = self._get_temperature_by(hourly_detail_nodes[i], SELECTOR_FOR_HOURLY_FEELS_LIKE_TEMPERATURE)
            uv_index = self._get_uv_index(hourly_detail_nodes[i], SELECTOR_FOR_HOURLY_UV_INDEX)
            probability_of_precipitation = self._get_percent_by(hourly_summary_nodes[i], SELECTOR_FOR_HOURLY_PROBABILITY_OF_PRECIPITATION)
            amount_of_precipitation = self._get_amount_of_precipitation(hourly_detail_nodes[i])
            wind = self._get_wind_by(hourly_detail_nodes[i], SELECTOR_FOR_HOURLY_WIND)
            icon_name = self._get_icon_name(hourly_summary_nodes[i])

            hourly_weather_forecasts.append(HourlyWeatherForecast(epoch_time=epoch_time,
                                                                  current_temperature=current_temperature,
                                                                  feels_like_temperature=feels_like_temperature,
                                                                  uv_index=uv_index,
                                                                  probability_of_precipitation=probability_of_precipitation,
                                                                  precipitation_type=PrecipitationType.RAIN,
                                                                  amount_of_precipitation =amount_of_precipitation,
                                                                  wind= wind,
                                                                  icon= icon_name,
                                                                  link= url))
            
        return hourly_weather_forecasts

    def get_daily_weather(self, request: WeatherForecastRequest) -> List[DailyWeatherForecast]:
        r"""
            This method return daily weather forecast from site 'The Weather Channel'.
            Data is obtained by parsing html page using selectolax.HTMLParser

            Request:
                WeatherForecastRequest:
                    place_id: str - id which you can get from WeatherClient.get_sun_location_search()
                    localization: str - localization string in standard 'BCP 47'. For example: 'uk-UA'
                    unit_type: UnitType
        """
        pass

    def search_location_id(self, request: GetLocationSearchRequest) -> List[Location]:
        r"""
            This method return similar places by place details. 
            Information gets from WeatherClient.get_sun_location_search()

            Request:
                GetLocationSearchRequest:
                    localization: str - localization string in standard 'BCP 47'. For example: 'uk-UA'
                    place_details: str - string with detail information about your place, for example: city name, postcode, address etc

        """
        pass

    def _get_wind_by(self, html_parser: HTMLParser, selector: str) -> Wind:
        r"""
            This selector return spans likes:
                <span> WSW </span>
                <span> 6 </span>
                <span> km/h </span>
            So, second element is weather speed
        """
        wind_elements = html_parser.css(selector)
       
        # 0 - direction, 1 - speed
        direction = wind_elements[0].text()[:-1]
        speed = wind_elements[1].text()
        self.logger.debug(f"got direction: {direction}, speed: {speed}")
        return Wind(direction= direction, 
                    speed= speed)

    def _get_url_for_get_weather_page(self, interval: str, request: WeatherForecastRequest) -> str:
        return WEATHER_URL_FORMAT.format(localization= request.localization, 
                                        interval= interval, 
                                        place_id= request.place_id,
                                        unit= self._map_unit_type(request.unit_type))
        
    def _get_html_parser_for_weather_page(self, url: str) -> HTMLParser:
        self.logger.debug(f"prepared request for get weather page, url: {url}")
        weather_forecast_page = requests.get(url).text
        html_parser = HTMLParser(weather_forecast_page)

        if self._page_is_not_found(html_parser):
            self.logger.debug("server return page with 'not_found' error")
            raise requests.RequestException("invalid request volumes")

        self.logger.debug("got html parser for weather page")

        return html_parser

    def _map_unit_type(self, unit_type: UnitType) -> str:
        if unit_type == UnitType.METRIC:
            return METRIC_UNIT_TYPE
        
        return IMPERIA_UNIT_TYPE
            

    def _page_is_not_found(self, html_parser: HTMLParser) -> bool:
        not_found_node = html_parser.css_first(SELECTOR_FOR_NOT_FOUND)
        return type(not_found_node) != NoneType

    def _get_temperature_by(self, html_parser: HTMLParser, selector: str):
        self.logger.debug(f"try return temperature by selector:{selector}")
        temperature_element = html_parser.css_first(selector)

        # Remove degree symbol from text
        temperature = temperature_element.text().replace("°", "")
        self.logger.debug(f"got temperature: {temperature}")
        return temperature

    def _get_percent_by(self, html_parser: HTMLParser, selector: str) -> int:
        perсent_element = html_parser.css_first(selector)

        percent = perсent_element.text().replace("%", "")
        self.logger.debug(f"got percent: {percent}")
        return percent
        
    def _get_uv_index(self, html_parser: HTMLParser, selector: str) -> int:
        uv_index_element = html_parser.css_first(selector).text()

        uv_index = re.search("^\\d{1,2}", uv_index_element).group(0)
        self.logger.debug(f"got uv_index: {uv_index}")
        return uv_index
    
    def _get_current_icon_name(self, html_parser: HTMLParser):
        icon_name = html_parser.css_first(SELECTOR_FOR_CURRENT_ICON_NAME).text()
        self.logger.debug(f"got icon_name: {icon_name}")
        return icon_name
    
    def _get_amount_of_precipitation(self, html_parser: HTMLParser):
        amount_of_precipitation = html_parser.css_first(SELECTOR_FOR_HOURLY_AMOUNT_OF_PRECIPITATION).text()
        self.logger.debug(f'got amount_of_precipitation: {amount_of_precipitation}')
        return amount_of_precipitation

    def _get_icon_name(self, html_parser: HTMLParser):
        icon_name = html_parser.css_first(SELECTOR_FOR_HOURLY_ICON_NAME).text()
        self.logger.debug(f'got icon_name: {icon_name}')
        return icon_name

    def _get_current_visibility(self, html_parser: HTMLParser):
        self.logger.debug("try get visibility node for current weather")
        visibility_node = html_parser.css_first(SELECTOR_FOR_CURRENT_VISIBILITY)

        if type(visibility_node) == NoneType:
            self.logger.debug("visibility is unlimited, return '0'")
            # Searching by selector SELECTOR_FOR_VISIBILITY return NoneType, because in the web, if visibility unlimited, display string `Unlimited`
            return 0
        
        # If visibility is limited, 'visibility_text' contains some float volume
        return visibility_node.text()
