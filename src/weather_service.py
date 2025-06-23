import datetime
import re
import logging
from types import NoneType
from typing import List
import requests
from selectolax.parser import HTMLParser
from selectolax.parser import Node
import time

from .models import Celestial, CurrentWeatherForecast, DailyWeatherForecast, DailyWeatherForecastDetail, GetLocationSearchRequest, GetSunV3LocationSearchRequest, HourlyWeatherForecast, Location, PrecipitationType, UnitType, WeatherForecastRequest, Wind
from .weather_client import WeatherClient

# language = standard 'BCP 47': uk-UA, 
# interval = today / hourbyhour / tenday, 
# place_id = 0b8697c01baca04214b4abd206319d3eba60db5fb05c191012c4e02f6bdb23a4, 
# unit = m / e / h
WEATHER_URL_FORMAT = "https://weather.com/{language}/weather/{interval}/l/{place_id}?unit={unit}"

# unit type on the weather channel backend
METRIC_UNIT_TYPE = "m"
IMPERIA_UNIT_TYPE = "e"
HYBRID_UNIT_TYPE = "h"

CURRENT_FORECAST_INTERVAL = "today"
HOULY_FORECAST_INTERVAL = "hourbyhour"
DAILY_FORECAST_INTERVAL = "tenday"

r"""
    This selectors using for get volumes from weather forecast page.
    Way for get all volumes:
        1. Get block where stored all weather forecast info
        2. Get volumes using relative selector by block info
        3. If need, convert or calculter received volumes

    For example, get probability of precipitation from daily weather forecast page.
    First get block with detail info about day weather forecast. After that get from this block needed volume

    detail_block = daily_page.css_first(SELECTOR_FOR_DAILY_PROBABILITY_OF_PRECIPITATION)
    pop = detail_block.css_first(SELECTOR_FOR_HOURLY_PROBABILITY_OF_PRECIPITATION) 
"""
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
SELECTOR_FOR_HOURLY_ICON_NAME = 'svg[class*="DetailsSummary--wxIcon--"] > title'
SELECTOR_FOR_HOURLY_PRECIPITATION_TYPE = 'li[class*="DetailsTable--listItem--"][data-testid="AccumulationSection"] > span > svg[class*="DetailsTable--icon--"] > title'

SELECTOR_FOR_DAILY_LAST_UPDATED_TIME = 'div[class*="DailyForecast--timestamp--"]'
SELECTOR_FOR_DAILY_BLOCK_SUMMARY_FORECAST = 'div[class*="DetailsSummary--DetailsSummary--"]'
SELECTOR_FOR_DAILY_BLOCK_DETAIL_TABLE_FORECAST = 'div[class*="DaypartDetails--Content--"]'
SELECTOR_FOR_DAILY_DETAILS_TABLE = 'div[class*="DaypartDetails--DetailsTable--"]'
SELECTOR_FOR_DAILY_TEMPERATURE = 'span[class*="DailyContent--temp--"][data-testid="TemperatureValue"]'
SELECTOR_FOR_DAILY_PROBABILITY_OF_PRECIPITATION = 'span[class*="DailyContent--value--"][data-testid="PercentageValue"]'
SELECTOR_FOR_DAILY_HUMIDITY = 'div[class*="DetailsTable--field--"] > span[class*="DetailsTable--value--"][data-testid="PercentageValue"]'
SELECTOR_FOR_DAILY_UV_INDEX = 'div[class*="DetailsTable--field--"] > span[class*="DetailsTable--value--"][data-testid="UVIndexValue"]'
SELECTOR_FOR_DAILY_RISE_SET_TIME = 'div[class*="DetailsTable--field--"] > span[class*="DetailsTable--value--"]:not([data-testid=PercentageValue]):not([data-testid="UVIndexValue"])'
SELECTOR_FOR_DAILY_ICON_NAME = 'svg[class*="DailyContent--weatherIcon--"] > title'
SELECTOR_FOR_DAILY_PRECIPITATION_TYPE = 'svg[class*="DailyContent--precipIcon--"] > title'
SELECTOR_FOR_DAILY_DAY = 'span[class*="DailyContent--daypartDate--"]'

SELECTOR_FOR_WIND = 'span[class*="Wind--windWrapper--"][data-testid="Wind"]'
SELECTOR_FOR_NOT_FOUND = 'div[class*="NotFound--notFound--"]'

INDEX_FOR_FIRST_DAILY_BLOCK = 0
INDEX_FOR_SECOND_DAILY_BLOCK = 1

SECONDS_IN_DAY = 24 * 60 * 60

class NotFoundException(BaseException):
    def __init__(self):
        super().__init__("server return not found page")

class WeatherService:
    r"""
        WeatherService has methods for return weather forecast by interval or search placeId using apies 'The Weather Channel'
    """

    def __init__(self, logger: logging.Logger, weather_client: WeatherClient):
        self._logger = logger.getChild("weather-service")
        self._weather_client = weather_client

    def get_current_weather(self, request: WeatherForecastRequest) -> CurrentWeatherForecast:
        r"""
            This method return current weather forecast from site 'The Weather Channel'.
            Data is obtained by parsing html page using selectolax.HTMLParser

            Request:
                WeatherForecastRequest:
                    place_id: str - id which you can get from WeatherClient.get_sun_location_search()
                    language: str - language string in standard 'BCP 47'. For example: 'uk-UA'
                    unit_type: UnitType
        """

        self._logger.info("start getting current weather forecast")
        url = self._get_url_for_get_weather_page(CURRENT_FORECAST_INTERVAL, request)
        current_weather_page = self._get_html_parser_for_weather_page(url)
        
        epoch_time = int(time.time())
        visibility = self._get_current_visibility(current_weather_page)
        feels_like_temperature = self._get_temperature_by(current_weather_page, SELECTOR_FOR_CURRENT_FEELS_LIKE_TEMPERATURE)
        current_temperature = self._get_temperature_by(current_weather_page, SELECTOR_FOR_CURRENT_CURRENT_TEMPERATURE)
        min_temperature = self._get_temperature_by(current_weather_page, SELECTOR_FOR_CURRENT_MIN_TEMPERATURE)
        max_temperature = self._get_temperature_by(current_weather_page, SELECTOR_FOR_CURRENT_MAX_TEMPERATURE)
        icon_name = self._get_current_icon_name(current_weather_page)
        
        self._logger.info("current weather forecast successful formed")
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
                    language: str - language string in standard 'BCP 47'. For example: 'uk-UA'
                    unit_type: UnitType
        """

        self._logger.debug("start getting hourly weather forecast")
        url = self._get_url_for_get_weather_page(HOULY_FORECAST_INTERVAL, request)
        hourly_weather_page = self._get_html_parser_for_weather_page(url)

        hourly_weather_forecasts = []
        hourly_summary_nodes = hourly_weather_page.css(SELECTOR_FOR_HOURLY_BLOCK_SUMMARY_FORECAST)
        hourly_detail_nodes = hourly_weather_page.css(SELECTOR_FOR_HOURLY_BLOCK_DETAIL_TABLE_FORECAST)

        for i in range(len(hourly_summary_nodes)):
            epoch_time=0 # TODO
            current_temperature = self._get_temperature_by(hourly_summary_nodes[i], SELECTOR_FOR_HOURLY_CURRENT_TEMPERATURE)
            feels_like_temperature = self._get_temperature_by(hourly_detail_nodes[i], SELECTOR_FOR_HOURLY_FEELS_LIKE_TEMPERATURE)
            uv_index = self._get_uv_index(hourly_detail_nodes[i], SELECTOR_FOR_HOURLY_UV_INDEX)
            probability_of_precipitation = self._get_percent_by(hourly_summary_nodes[i], SELECTOR_FOR_HOURLY_PROBABILITY_OF_PRECIPITATION)
            precipitation_type = self._which_precipitation_type(hourly_detail_nodes[i], SELECTOR_FOR_HOURLY_PRECIPITATION_TYPE)
            amount_of_precipitation = self._get_amount_of_precipitation(hourly_detail_nodes[i])
            wind = self._get_wind_by(hourly_detail_nodes[i])
            icon_name = self._get_icon_name(hourly_summary_nodes[i], SELECTOR_FOR_HOURLY_ICON_NAME)

            hourly_weather_forecasts.append(HourlyWeatherForecast(epoch_time=epoch_time,
                                                                  current_temperature=current_temperature,
                                                                  feels_like_temperature=feels_like_temperature,
                                                                  uv_index=uv_index,
                                                                  probability_of_precipitation=probability_of_precipitation,
                                                                  precipitation_type=precipitation_type,
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
                    language: str - language string in standard 'BCP 47'. For example: 'uk-UA'
                    unit_type: UnitType
            
            after 15:00 local time, first element daily weather forecast return only for night 
        """

        self._logger.debug("start getting daily weather forecast")
        url = self._get_url_for_get_weather_page(DAILY_FORECAST_INTERVAL, request)
        daily_weather_page = self._get_html_parser_for_weather_page(url)

        daily_weather_forecasts = []
        daily_summary_nodes = daily_weather_page.css(SELECTOR_FOR_DAILY_BLOCK_SUMMARY_FORECAST)
        daily_detail_nodes = daily_weather_page.css(SELECTOR_FOR_DAILY_BLOCK_DETAIL_TABLE_FORECAST)

        r"""
            For first element need check local time, because backend return only night block between time 15:00 - 03:00

            I write copy-paste code because I dont want check every iteration "is first day?" then "is night?"
            In this work flow just once check "is night?".
        """

        start_index = 0
        day = self._get_day_from_daily_detail(daily_detail_nodes[0])
        last_updated_time = self._get_local_last_updated_time(daily_weather_page, SELECTOR_FOR_DAILY_LAST_UPDATED_TIME)
        epoch_time = self._get_local_datetime_by(day, last_updated_time)
        local_time_seconds = self._get_seconds_by_time(last_updated_time)
        if self._is_night(local_time_seconds):
            night = self._get_daily_detail_forecast(daily_summary_nodes[0], daily_detail_nodes[0], INDEX_FOR_FIRST_DAILY_BLOCK)

            daily_weather_forecasts.append(
                 DailyWeatherForecast(
                    epoch_time=epoch_time,
                    day=None,
                    night=night,
                    link=url
                )
            )
            start_index = 1


        for i in range(start_index, len(daily_summary_nodes)):
            day = self._get_daily_detail_forecast(daily_summary_nodes[i], daily_detail_nodes[i], INDEX_FOR_FIRST_DAILY_BLOCK)
            night = self._get_daily_detail_forecast(daily_summary_nodes[i], daily_detail_nodes[i], INDEX_FOR_SECOND_DAILY_BLOCK)

            daily_weather_forecasts.append(
                DailyWeatherForecast(
                    epoch_time=epoch_time,
                    day=day,
                    night=night,
                    link=url
                )
            )

            epoch_time += SECONDS_IN_DAY
        
        return daily_weather_forecasts


    def search_location_id(self, request: GetLocationSearchRequest) -> List[Location]:
        r"""
            This method return similar places by place details. 
            Information gets from WeatherClient.get_sun_location_search()

            Request:
                GetLocationSearchRequest:
                    localization: str - localization string in standard 'BCP 47'. For example: 'uk-UA'
                    place_details: str - string with detail information about your place, for example: city name, postcode, address etc

        """
        list_of_found_locations = self._weather_client.get_sun_location_search(GetSunV3LocationSearchRequest(language=request.localization,
                                                                                                             place_detail=request.place_details))
        locations_for_response = []
        for location in list_of_found_locations:
            locations_for_response.append(
                Location(
                    place_id=location.place_id,
                    address=location.address,
                    city=location.city,
                    country=location.country,
                    latitude=location.latitude,
                    longitude=location.longitude,
                    postal_code=location.postal_code))
        
        return locations_for_response
    
    def _get_local_datetime_by(self, day: int, time: datetime.time) -> datetime.datetime:
        self._logger.debug(f'start formatting local datetime by: day={day}, time={time}')
        current_datetime = datetime.datetime.now()
        local_datetime = datetime.datetime(year= current_datetime.year, 
                                           month=current_datetime.month,
                                           day=day,
                                           hour=time.hour,
                                           minute=time.minute).timestamp() * 1000

        self._logger.debug(f"formated local datetime: {local_datetime}")
        return local_datetime 
    
    def _get_local_last_updated_time(self, html_page: HTMLParser, selector) -> datetime.time:
        # Will be returned like: 'As of 03:48 GMT-03:00'
        self._logger.debug("start get local last updated time")
        last_updated_str = html_page.css_first(selector).text()
        
        # This regex will be return '03:48' from string: 'As of 03:48 GMT-03:00'
        self._logger.debug(f"found last updated time: {last_updated_str}")
        time_str = re.search(r'(?<=\s)\d{1,2}:\d{2}', last_updated_str).group(0)
        
        self._logger.debug(f"regex return: {time_str} time")
        return  datetime.datetime.strptime(time_str, "%H:%M").time()

    def _get_day_from_daily_detail(self, detail_node: Node) -> int:
        self._logger.debug('start get day from daily detail')
        day_with_week_name = detail_node.css_first(SELECTOR_FOR_DAILY_DAY).text()
        self._logger.debug(f'returned day with week: {day_with_week_name}')
        # Will be return '01' from 'Thu 01'
        day = re.search(r"\d{2}", day_with_week_name).group(0)
        self._logger.debug(f'found day number: {day}')
        return int(day)
        
    def _is_night(self, time_in_seconds: int) -> bool:
        evening = self._get_seconds_by_str_time("15:00")
        morning = self._get_seconds_by_str_time("03:00")

        # Backend return only night block between time 15:00 - 03:00
        return time_in_seconds >= evening or time_in_seconds <= morning

    def _which_precipitation_type(self, html_parser: HTMLParser, selector: str, index_of_element: int = 0) -> PrecipitationType:
        precipitation_type_title = html_parser.css(selector)[index_of_element]

        if type(precipitation_type_title) == NoneType:
            self._logger.debug("title of precipitation type not found, return type 'NONE'")
            return PrecipitationType.NONE
        
        precipitation_type_text = precipitation_type_title.text()

        if "Rain" in precipitation_type_text:
            self._logger.debug("title of precipitation type is rain, return 'RAIN'")
            return PrecipitationType.RAIN
        
        if "Mixed" in precipitation_type_text:
            self._logger.debug("title of precipitation type is mixed, return 'MIXED'")
            return PrecipitationType.MIXED
        
        if "Snowflake" == precipitation_type_text:
            self._logger.debug("title of precipitation type is mixed, return 'SNOW'")
            return PrecipitationType.SNOW

        return PrecipitationType.ICE


    def _get_wind_by(self, html_parser: HTMLParser, index_of_element: int = 0) -> Wind:
        r"""
            Wind block have similar selector, Parser return element like that:
            <span> # element
                <span> WSW </span>  # volume
                <span> 6 </span>    #
                <span> km/h </span> #
            <span>
            So, second element is weather speed
        """
        wind_elements = html_parser.css(SELECTOR_FOR_WIND)
        wind_element = wind_elements[index_of_element]
       
        # 0 - direction, 1 - speed
        wind_volumes = wind_element.css(f"span:not({SELECTOR_FOR_WIND})")
        direction = wind_volumes[0].text()[:-1]
        speed = wind_volumes[1].text()
        self._logger.debug(f"got direction: {direction}, speed: {speed}")
        return Wind(direction= direction, 
                    speed= speed)

    def _get_daily_detail_forecast(self, summary_block: Node, detail_block: Node, index_of_block: int):
        temperature = self._get_temperature_by(detail_block, SELECTOR_FOR_DAILY_TEMPERATURE, index_of_block)
        humidity = self._get_percent_by(detail_block, SELECTOR_FOR_DAILY_HUMIDITY, index_of_block)
        wind = self._get_wind_by(detail_block, index_of_block)
        celestial = self._get_celestial_by(detail_block, index_of_block)
        probability_of_precipitation = self._get_percent_by(detail_block, SELECTOR_FOR_DAILY_PROBABILITY_OF_PRECIPITATION, index_of_block)
        precipitation_type = self._which_precipitation_type(detail_block, SELECTOR_FOR_DAILY_PRECIPITATION_TYPE, index_of_block)
        icon_name = self._get_icon_name(detail_block, SELECTOR_FOR_DAILY_ICON_NAME, index_of_block)

        return DailyWeatherForecastDetail(
            temperature= temperature,
            humidity= humidity,
            wind= wind,
            rise_time= celestial.rise_time,
            set_time= celestial.set_time,
            probability_of_precipitation= probability_of_precipitation,
            precipitation_type= precipitation_type,
            icon_name= icon_name)

    def _get_url_for_get_weather_page(self, interval: str, request: WeatherForecastRequest) -> str:
        return WEATHER_URL_FORMAT.format(language= request.language, 
                                         interval= interval, 
                                         place_id= request.place_id,
                                         unit= self._map_unit_type(request.unit_type))
        
    def _get_html_parser_for_weather_page(self, url: str) -> HTMLParser:
        self._logger.debug(f"prepared request for get weather page, url: {url}")
        weather_forecast_page = requests.get(url).text
        html_parser = HTMLParser(weather_forecast_page)

        if self._page_is_not_found(html_parser):
            self._logger.debug("server return page with 'not_found' error")
            raise NotFoundException

        self._logger.debug("got html parser for weather page")

        return html_parser

    def _map_unit_type(self, unit_type: UnitType) -> str:
        if unit_type == UnitType.METRIC:
            return METRIC_UNIT_TYPE
        
        return IMPERIA_UNIT_TYPE

    def _get_seconds_by_str_time(self, time_str: str) -> int:
        time = datetime.datetime.strptime(time_str, "%H:%M").time()
        return self._get_seconds_by_time(time)

    def _get_seconds_by_time(self, time: datetime.time) -> int:
        return time.hour * 3600 + time.minute * 60

    def _get_celestial_by(self, node: Node, index_of_element: int = 0) -> Celestial:
        r"""
            This method will be return array with 4 nodes where:
                <span>06:00</span>  # index of list [0] = sunrise
                <span>19:00</span>  # index of list [1] = sunset
                <span>19:00</span>  # index of list [2] = moonrise
                <span>05:00</span>  # index of list [3] = moonset
        """
        time_str_elements = node.css(SELECTOR_FOR_DAILY_RISE_SET_TIME)

        if len(time_str_elements) == 0:
            self._logger.debug("selector for celestial not found elements")
            raise TypeError("selector for celestial not found elements")


        if index_of_element == 0:
            return Celestial(
                rise_time=self._get_seconds_by_str_time(time_str_elements[0].text()),
                set_time=self._get_seconds_by_str_time(time_str_elements[1].text())
            )
        
        return Celestial(
            rise_time=self._get_seconds_by_str_time(time_str_elements[0].text()),
            set_time=self._get_seconds_by_str_time(time_str_elements[1].text())
        )

    def _page_is_not_found(self, html_parser: HTMLParser) -> bool:
        not_found_node = html_parser.css_first(SELECTOR_FOR_NOT_FOUND)
        return type(not_found_node) != NoneType

    def _get_temperature_by(self, html_parser: HTMLParser, selector: str, element_index: int = 0):
        self._logger.debug(f"try return temperature by selector:{selector}")
        temperature_elements = html_parser.css(selector)
        self._logger.debug(f"temperatures: {temperature_elements[0].text()}")

        if len(temperature_elements) == 0:
            self._logger.debug("list with temperature element is empty")
            raise TypeError("selector for temperature not found elements")

        # Remove degree symbol from text
        temperature = temperature_elements[element_index].text().replace("°", "")
        self._logger.debug(f"got temperature: {temperature}")
        return temperature

    def _get_percent_by(self, html_parser: HTMLParser, selector: str, index_of_element: int = 0) -> int:
        self._logger.debug(f"try return percent by selector:{selector}")
        perсent_elements = html_parser.css(selector)

        if len(perсent_elements) == 0:
            self._logger.debug("list with percent element is empty")
            raise TypeError("selector for percent not found elements")
        
        percent = perсent_elements[index_of_element].text().replace("%", "")
        self._logger.debug(f"got percent: {percent}")
        return percent
        
    def _get_uv_index(self, html_parser: HTMLParser, selector: str, index_of_element: int = 0) -> int:
        uv_index_elements = html_parser.css(selector)

        if len(uv_index_elements) == 0:
            self._logger.debug("list with uv index elements is empty")
            raise TypeError("selector for uv index not found element")
        
        uv_index = re.search("^\\d{1,2}", uv_index_elements[index_of_element].text()).group(0)
        self._logger.debug(f"got uv_index: {uv_index}")
        return uv_index
    
    def _get_current_icon_name(self, current_weather_page: HTMLParser):
        icon_name = current_weather_page.css_first(SELECTOR_FOR_CURRENT_ICON_NAME).text()
        self._logger.debug(f"got icon_name: {icon_name}")
        return icon_name
    
    def _get_amount_of_precipitation(self, html_parser: HTMLParser):
        amount_of_precipitation = html_parser.css_first(SELECTOR_FOR_HOURLY_AMOUNT_OF_PRECIPITATION).text()
        self._logger.debug(f'got amount_of_precipitation: {amount_of_precipitation}')
        return amount_of_precipitation

    def _get_icon_name(self, html_parser: HTMLParser, selector: str, index_of_element: int = 0):
        icon_name_element = html_parser.css(selector)[index_of_element]
        icon_name = icon_name_element.text()
        self._logger.debug(f'got icon_name: {icon_name}')
        return icon_name

    def _get_current_visibility(self, current_weather_page: HTMLParser):
        self._logger.debug("try get visibility node for current weather")
        visibility_node = current_weather_page.css_first(SELECTOR_FOR_CURRENT_VISIBILITY)

        if type(visibility_node) == NoneType:
            self._logger.debug("visibility is unlimited, return '0'")
            # Searching by selector SELECTOR_FOR_VISIBILITY return NoneType, because in the web, if visibility unlimited, display string `Unlimited`
            return 0
        
        # If visibility is limited, 'visibility_text' contains some float volume
        return visibility_node.text()
