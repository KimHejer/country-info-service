## Country Information Service


This API provides data for countries using their ISO2 country code. The information is delivered through GET request
and is fetched from two external API's:

* CountriesNow API: 

Endpoint: http://129.241.150.113:3500/api/v0.1/  

Documentation: https://documenter.getpostman.com/view/1134062/T1LJjU52

* REST Countries API:

Endpoint: http://129.241.150.113:8080/v3.1/ 

Documentation: http://129.241.150.113:8080/

# Features:
1. Get general country information by country code
2. Get population data from specified country with country code
3. Get API status


# API Endpoints:

1. Get country info:
```bash
GET /countryinfo/v1/info/{countryCode}?limit={cityCount}
```
If not provided, default limit is set to 10 cities.

Example request:
```bash
GET /countryinfo/v1/info/no?limit=5
```
Response:
```JSON
{
    "name": "Norway",
    "continent": "Europe",
    "population": 5379475,
    "languages": {
        "nno": "Norwegian Nynorsk",
        "nob": "Norwegian Bokmål",
        "smi": "Sami"
    },
    "borders": [
        "FIN",
        "SWE",
        "RUS"
    ],
    "flag": "https://flagcdn.com/no.svg",
    "capital": "Oslo",
    "cities": [
        "Abelvaer",
        "Adalsbruk",
        "Adland",
        "Agotnes",
        "Agskardet"
    ]
}
```

2. Get population data:
```bash
GET /countryinfo/v1/population/{countryCode}?limit={startYear-endYear}
```
If limit is not provided, all documentet years will be fetched.

Example request:
```bash
GET /countryinfo/v1/population/in?limit=2010-2015
```
Response:
```json
{
    "mean": 1272825900,
    "values": [
        {
            "year": 2010,
            "value": 1234281170
        },
        {
            "year": 2011,
            "value": 1250288729
        },
        {
            "year": 2012,
            "value": 1265782790
        },
        {
            "year": 2013,
            "value": 1280846129
        },
        {
            "year": 2014,
            "value": 1295604184
        },
        {
            "year": 2015,
            "value": 1310152403
        }
    ]
}
```

3. Get API status:
```bash
GET /countryinfo/v1/status
```
Response:
```json
{
    "countriesnowapi": "200",
    "restcountriesapi": "200",
    "version": "v1",
    "uptime": 850
}
```

# Possible responses:
200 - OK, succesfull request and valid data returned

400 - Bad request, missing parameters or invalid input

404 - Not found, the requested resource was not found

500 - Internal server error, something went wrong on the server