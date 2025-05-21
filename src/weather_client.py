from logging import Logger
from typing import List
import requests
from models import GetSunV3LocationSearchRequest, GetSunV3LocationSearchResponse


class WeatherClient:
    r"""
        WeatherClient has methods for call 'The Weather Channel' backend
    """
    def __init__(self, logger: Logger):
        self.logger = logger

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

        self.logger.debug("prepare request for search location using 'redux-dal'")
        response = requests.post("https://weather.com/api/v1/p/redux-dal", json=request_format)
        
        if response.status_code != 200:
            self.logger.error(f"http call returned status code: {response.status_code}")
            raise Exception("request return error status code.")
            
        self.logger.debug("response successful received")
        response_data = response.json()
        place_infos = response_data["dal"][method_name][f"language:{request.language};locationType:{location_type};query:{request.place_detail}"]["data"]["location"]

        locations = []
        self.logger.debug("start map received response to 'GetSunV3LocationSearchResponse'")
        for i in range(len(place_infos["address"])):
            self.logger.debug(f"mapping place_info= {i}")
            locations.append(GetSunV3LocationSearchResponse(place_infos["address"][i], 
                                                          place_infos["adminDistrict"][i],
                                                          place_infos["city"][i],
                                                          place_infos["country"][i],
                                                          place_infos["countryCode"][i],
                                                          place_infos["displayName"][i],
                                                          place_infos["ianaTimeZone"][i],
                                                          place_infos["latitude"][i],
                                                          place_infos["longitude"][i],
                                                          place_infos["placeId"][i],
                                                          place_infos["postalCode"][i]))
        
        self.logger.debug('return mapped locations')
        return locations