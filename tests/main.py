import unittest
from tests.weather_service_test import CurrentWeatherForecastTests, HourlyWeatherForecastTests, DailyWeatherForecastTests
from tests.weather_client_test import WeatherClientTests

if __name__ == '__main__':
    suite = unittest.TestSuite()
    
    suite.addTest(unittest.makeSuite(CurrentWeatherForecastTests))
    suite.addTest(unittest.makeSuite(HourlyWeatherForecastTests))
    suite.addTest(unittest.makeSuite(DailyWeatherForecastTests))
    suite.addTest(unittest.makeSuite(WeatherClientTests))

    runner = unittest.TextTestRunner()
    runner.run(suite)