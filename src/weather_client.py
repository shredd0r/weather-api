from logging import Logger
from typing import List
import requests
from .models import GetSunV3LocationSearchRequest, GetSunV3LocationSearchResponse

class BadRequestException(BaseException):
    def __init__(self):
        super().__init__("server return bad request status code")

class InternalErrorException(BaseException):
    def __init__(self):
        super().__init__("server return internal error status code")

class WeatherClient:
    r"""
        WeatherClient has methods for call 'The Weather Channel' backend
    """
    def __init__(self, logger: Logger):
        self.__logger = logger.getChild("weather-client")

    def get_sun_location_search(self, request: GetSunV3LocationSearchRequest) -> List[GetSunV3LocationSearchResponse]:
        r"""
            This method call weather api 'redux-dal' and map response to GetSunV3LocationSearchResponse.

            'redux-dal' is one endpoint for many apies. Method name 'getSunV3LocationSearchUrlConfig' returned 
            places with similar parameters like: city name, postal code, city address, etc.

            For example, you can write in 'place_detail' city name: "Kyiv" and api will returned all matches with name 'Kyiv'
        """
        method_name = "getSunV3LocationSearchUrlConfig"
        location_type = "locale"
        request_format = [
            {
                "name": method_name,
                "params": {
                    "language": request.language,
                    "locationType": location_type,
                    "query": request.place_detail
                }
            }
        ]

        self.__logger.debug("prepare request for search location using 'redux-dal'")
        response = requests.post("https://weather.com/api/v1/p/redux-dal", json=request_format)
        
        if response.status_code >= 400 and response.status_code < 500:
            self.__logger.error(f"http call returned status code: {response.status_code}")
            raise BadRequestException
        
        if response.status_code == 500:
            self.__logger.error('server has a problem, returned 500 status code')
            raise InternalErrorException
            
        self.__logger.debug("response successful received")
        response_data = response.json()
        place_infos = response_data["dal"][method_name][f"language:{request.language};locationType:{location_type};query:{request.place_detail}"]["data"]["location"]

        locations = []
        self.__logger.debug("start map received response to 'GetSunV3LocationSearchResponse'")
        for i in range(len(place_infos["address"])):
            self.__logger.debug(f"mapping place_info= {i}")
            locations.append(GetSunV3LocationSearchResponse(address=place_infos["address"][i], 
                                                            admin_district=place_infos["adminDistrict"][i],
                                                            city=place_infos["city"][i],
                                                            country=place_infos["country"][i],
                                                            country_code=place_infos["countryCode"][i],
                                                            display_name=place_infos["displayName"][i],
                                                            iana_time_zone=place_infos["ianaTimeZone"][i],
                                                            latitude=place_infos["latitude"][i],
                                                            longitude=place_infos["longitude"][i],
                                                            place_id=place_infos["placeId"][i],
                                                            postal_code=place_infos["postalCode"][i]))
        
        self.__logger.debug('return mapped locations')
        return locations