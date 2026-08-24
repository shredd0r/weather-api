<a id="readme-top"></a>

<h3 align="center">Weather Channel API</h3>

  <p align="center">
    Microservice which can return today, hourly and daily weather forecast by http apies.
    <br />
    <a href="https://github.com/shredd0r/weather-api"><strong>Explore the docs »</strong></a>
    <br />
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

## About The Project

The service takes information about weather from 'The Weather Channel' and return the response by REST interface.
The service has 4 apies:
1. Get today weather forecast
2. Get hourly weather forecast (3 days)
3. Get daily weather forecast (10 days)
4. Get place id by your city details

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

* [![Python][Python]][Python-url]
* [![Playwright][Playwright]][Playwright-url]
* [![FastAPI][FastAPI]][FastAPI-url]
* [![Docker][Docker]][Docker-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Run by Docker

1. Build Dockerfile image
```sh
    docker build -t weather-api:v1.0.0 .
```
2. Start image
```sh
    docker run -p 8000:8000 weather-api:v1.0.0
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->

## Disclaimer

This project is created strictly for **educational and informational purposes only**.

* **Copyright & Intellectual Property:** All content, trademarks, and assets belong to their respective owners. No copyright infringement is intended.
* **Data Storage:** This project does not host, store, or distribute any copyrighted materials or personal data on its servers.
* **Limitation of Liability:** The author assumes no responsibility or liability for any misuse of this software or code.

<!-- Shields.io badges. You can a comprehensive list with many more badges at: https://github.com/inttter/md-badges -->
[Python]: https://img.shields.io/badge/python-000000?style=for-the-badge&logo=python&logoColor=white
[Python-url]: https://www.python.org/
[Playwright]: https://img.shields.io/badge/playwright-000000?style=for-the-badge&logo=Playwright&logoColor=white
[Playwright-url]: https://playwright.dev/
[FastAPI]: https://img.shields.io/badge/FastAPI-000000?style=for-the-badge&logo=fastapi&logoColor=white
[FastAPI-url]: https://fastapi.tiangolo.com/
[Docker]: https://img.shields.io/badge/Docker-000000?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/