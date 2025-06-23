import logging
from typing import List
import unittest
from unittest.mock import MagicMock, Mock, patch

from src.models import GetSunV3LocationSearchRequest, GetSunV3LocationSearchResponse
from src.weather_client import BadRequestException, InternalErrorException, WeatherClient

logger = logging.getLogger("test")
weather_client = WeatherClient(logger)

class WeatherClientTests(unittest.TestCase):

    def get_default_request(self) -> GetSunV3LocationSearchRequest:
        return GetSunV3LocationSearchRequest(
            language="uk-UA",
            place_detail="Kharkiv")

    def test_successful_response_from_server(self):
        expected = self.get_expected_response()
        actual = weather_client.get_sun_location_search(self.get_default_request())
        
        self.assertEqual(actual, expected)

    @patch("src.weather_client.requests")
    def test_server_return_internal_error(self, mock_request: Mock):
        self.case_for_check_error_status_code(mock_request, 500, InternalErrorException)


    @patch("src.weather_client.requests")
    def test_server_return_bad_request_error(self, mock_request: Mock):
        self.case_for_check_error_status_code(mock_request, 400, BadRequestException)

    def case_for_check_error_status_code(self, mock_request: Mock, status_code: int, exception: BaseException):
        mock_response = MagicMock()
        mock_response.status_code= status_code

        mock_request.post.return_value= mock_response

        with self.assertRaises(exception):
            weather_client.get_sun_location_search(self.get_default_request())

    def get_expected_response(self) -> List[GetSunV3LocationSearchResponse]:
        return [
            GetSunV3LocationSearchResponse(address='Харків, Харківська область, Україна',
                                           admin_district='Харківська область', 
                                           city='Харків', 
                                           country='Україна', 
                                           country_code='UA', 
                                           display_name='Харків', 
                                           iana_time_zone='Europe/Kiev', 
                                           latitude=49.993, 
                                           longitude=36.232, 
                                           place_id='a3fc3a9546b3303604bb2e67c47280cd13107b83a2da060d38c4e0941f3105d2', 
                                           postal_code='61057'), 
            GetSunV3LocationSearchResponse(address='Харьков, Ширак, Вірменія', 
                                           admin_district='Ширак', 
                                           city='Харьков', 
                                           country='Вірменія', 
                                           country_code='AM', 
                                           display_name='Харьков', 
                                           iana_time_zone='Asia/Yerevan', 
                                           latitude=40.506,
                                           longitude=43.594, 
                                           place_id='f98241bf4b2c860885e8dcd3b57d884a055bbee803ac1ac8ece45b33999e5b0d', 
                                           postal_code='29'), 
            GetSunV3LocationSearchResponse(address='Kharki, Baghbor, Барпета, Ассам, Індія', 
                                           admin_district='Ассам', 
                                           city='Baghbor', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=26.213, 
                                           longitude=90.99, 
                                           place_id='6894e5d442cf3037b0a45b0b0eb39d11254a3188989a53ca38e57990116f53a5', 
                                           postal_code='781308'), 
            GetSunV3LocationSearchResponse(address='Kharki, Барасат, Северные 24 парганы, Західний Бенгал, Індія', 
                                           admin_district='Західний Бенгал', 
                                           city='Барасат', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=22.741, 
                                           longitude=88.564, 
                                           place_id='44e53a2c08fb8597e52954b13f7fdf50e6086a5af7e3a0f3513cf27fcc4acca1', 
                                           postal_code='743294'), 
            GetSunV3LocationSearchResponse(address='Kharki, Mathurapur, Південні 24 паргани, Західний Бенгал, Індія', 
                                           admin_district='Західний Бенгал', 
                                           city='Mathurapur', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=22.077, 
                                           longitude=88.433, 
                                           place_id='03621928c4c7e99c7fcd587e439f99b71d00f737b4fe6c4dc9bf38563c9b954b', 
                                           postal_code='743349'), 
            GetSunV3LocationSearchResponse(address='Kharki, Bishungarh, Хазарибагх, Джхаркханд, Індія', 
                                           admin_district='Джхаркханд', 
                                           city='Bishungarh', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=23.909, 
                                           longitude=85.865, place_id='524a35f25c929bc38e10c77947f922baac6d82d094c678a37bb18cd01f5c8bc9', postal_code='825312'), 
            GetSunV3LocationSearchResponse(address='Kharki, Kisko, Лохардага, Джхаркханд, Індія', 
                                           admin_district='Джхаркханд', 
                                           city='Kisko', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=23.522, 
                                           longitude=84.657, 
                                           place_id='a281e3663fb083778464fd2f83951b64b062ea4f53308e32d06ec9230909f924', 
                                           postal_code='835302'), 
            GetSunV3LocationSearchResponse(address='Kharki Khurd, Adhaura, Каймур, Біхар, Індія', 
                                           admin_district='Біхар', 
                                           city='Adhaura', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki Khurd', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=24.777, 
                                           longitude=83.607, 
                                           place_id='501df5e155b6bfe8ec98ee0cdbbc3cfc27006099ce1ee2c6854ee4e1ed9c66ad', 
                                           postal_code='821102'), 
            GetSunV3LocationSearchResponse(address='Kharki, Kodinga, Набарангпур, Одіша, Індія', 
                                           admin_district='Одіша', 
                                           city='Kodinga', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=19.295, 
                                           longitude=82.447, 
                                           place_id='a0c1133cc294c900c8121d2b912af24f18e7456dcd2dba41fa3824426a904e1e', 
                                           postal_code='764071'), 
            GetSunV3LocationSearchResponse(address='Kharki, Ganai Gangoli, Пітхораґарх, Уттаракханд, Індія', 
                                           admin_district='Уттаракханд', 
                                           city='Ganai Gangoli', 
                                           country='Індія', 
                                           country_code='IN', 
                                           display_name='Kharki', 
                                           iana_time_zone='Asia/Kolkata', 
                                           latitude=29.708, 
                                           longitude=79.969, 
                                           place_id='bb590dfe4b2f4fb3fd64662aa71a10a3788d9267ba10ce9592fb553790d3ad75', 
                                           postal_code='262532')]