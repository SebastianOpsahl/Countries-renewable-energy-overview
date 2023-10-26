# Overview
The service is a REST web application consisting of four resource root paths. The REST webservices we have used for the assignment are:

- REST Countries API (instance hosted for this course). Endpoint: http://129.241.150.113:8080/v3.1 Documentation: http://129.241.150.113:8080/

- Renewable Energy Dataset (Authors: Hannah Ritchie, Max Roser and Pablo Rosado (2022)) "Energy". Published online at OurWorldInData.org. Retrieved from: https://ourworldindata.org/energy

The services content retrieval and search operations are returned in JSON, while information about the service for online users are written in HTML.

# Deployment
The service will be deployed on openstack and will be available via the URL in the delivery system.
The service is deployed on the NTNU internal network, so to access the service you will either have to use the NTNU vpn or be logged onto the NTNU Wi-Fi.

# How to use
The user will have the option to use one of the four following resource root paths:
```
/energy/v1/renewables/current
/energy/v1/renewables/history
/energy/v1/notifications/
/energy/v1/status/
```

renewables/current and renewables/history is root paths for searching information about different countries's percentage of renewable energy.


To showcase our services we will use the following conventions for describing placeholders:
- {value} - mandatory value
- {value?} - optional value
- {?key=value} - mandatory parameter (key-value pair)
- {?key=value?} - optional parameter (key-value pair)

## Current percentage of renewables
Returns the current numbers (meaning 2021) of countries percentage of renewable energy. This will be done in the format:

Path: /energy/v1/renewables/current/{country?}/{neighbours=bool?}

Where "country" is either the country code or country name. We designed our application so the search has to be totally equal to the country name or the country code to get output.

Example requests:
```
/energy/v1/renewables/current/nor?neighbours=false
```
Response:
```
[{"name":"Norway","isoCode":"NOR","year":2021,"percentage":71.558365}]
```

```
/energy/v1/renewables/current/norway?neighbours=true
```
Response:
```
[{"name":"Norway","isoCode":"NOR","year":2021,"percentage":71.558365},
{"name":"Finland","isoCode":"FIN","year":2021,"percentage":34.61129},
{"name":"Sweden","isoCode":"SWE","year":2021,"percentage":50.924007},
{"name":"Russia","isoCode":"RUS","year":2021,"percentage":6.6202893}]
```


## Historical percentage of renewables
Returns all years of countries percentage of renewables as present in the data source. This will be done in the format:
Path: /energy/v1/renewables/history/{country?}{?begin=year&end=year?}{sortByValue=bool?}

Where country is either a countrycode or countryname, the **begin** year and **end** year can both be specified which prints out that interval. Begin can also only be specified which prints from that point and to current year, or only end year can be specified which prints from the start of renewable counting until the given end year. The service also provides the oppurtunity to sort the results in order via the sortByValue query.

Example request with sorting:
```
/energy/v1/renewables/history/norway?begin=2010&end=2020&sortByValue=true
```
Response:
```
[{"name":"Norway","isoCode":"NOR","year":2010,"percentage":65.47019},{"name":"Norway","isoCode":"NOR","year":2011,"percentage":66.30012},{"name":"Norway","isoCode":"NOR","year":2019,"percentage":67.08509},{"name":"Norway","isoCode":"NOR","year":2013,"percentage":67.50864},{"name":"Norway","isoCode":"NOR","year":2018,"percentage":68.85805},{"name":"Norway","isoCode":"NOR","year":2015,"percentage":68.87519},{"name":"Norway","isoCode":"NOR","year":2014,"percentage":68.88728},{"name":"Norway","isoCode":"NOR","year":2017,"percentage":69.260994},{"name":"Norway","isoCode":"NOR","year":2016,"percentage":69.86629},{"name":"Norway","isoCode":"NOR","year":2012,"percentage":70.095116},{"name":"Norway","isoCode":"NOR","year":2020,"percentage":70.96306}]
```

Example request without sorting:
```
/energy/v1/renewables/history/norway?begin=2010&end=2020&sortByValue=false
```
Response:
```
[{"name":"Norway","isoCode":"NOR","year":2010,"percentage":65.47019},{"name":"Norway","isoCode":"NOR","year":2011,"percentage":66.30012},{"name":"Norway","isoCode":"NOR","year":2012,"percentage":70.095116},{"name":"Norway","isoCode":"NOR","year":2013,"percentage":67.50864},{"name":"Norway","isoCode":"NOR","year":2014,"percentage":68.88728},{"name":"Norway","isoCode":"NOR","year":2015,"percentage":68.87519},{"name":"Norway","isoCode":"NOR","year":2016,"percentage":69.86629},{"name":"Norway","isoCode":"NOR","year":2017,"percentage":69.260994},{"name":"Norway","isoCode":"NOR","year":2018,"percentage":68.85805},{"name":"Norway","isoCode":"NOR","year":2019,"percentage":67.08509},{"name":"Norway","isoCode":"NOR","year":2020,"percentage":70.96306}]
```

## Notification endpoint
Is an endpoint where users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries (or any country) is invoked. Users can register multiple webhooks. Different from the other endpoints which all are GET the notification endpoint provides different methods dependent on what you wish to achieve.

### Registration of Webhook
Method: POST
Path: /energy/v1/notifications/

The request should contain the following:

The URL to be triggered upon event (the service that should be invoked)
the country for which the trigger applies (if empty, it applies to any invocation)
the number of invocations after which a notification is triggered (it should re-occur every number of invocations, i.e., if 5 is specified, it should occur after 5, 10, 15 invocation, and so on, unless the webhook is deleted).

Example of request:
```
{
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}
```

The response given will contain the ID for the registration that can be used to see detail information or to delete the webhook registration.

Response example:
```
{
    "webhook_id": "OIdksUDwveiwe"
}
```

### Deletion of Webhook
Method: DELETE
Path: /energy/v1/notifications/{id}

Where {id} is the ID returned during the webhook registration


### View registered webhook
Method: GET
Path: /energy/v1/notifications/{id}

The request should contain the {id} for the webhook registration

The response is similar to the POST request body, but further includes the ID assigned by the server upon adding the webhook.

Body example:
```
{
   "webhook_id": "OIdksUDwveiwe",
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}
```

### View all registered webhooks
Method: GET
Path: /energy/v1/notifications/

The response is a collection of all registered webhooks.

Body example:
```
[
   {
      "webhook_id": "OIdksUDwveiwe",
      "url": "https://localhost:8080/client/",
      "country": "NOR",
      "calls": 5
   },
   {
      "webhook_id": "DiSoisivucios",
      "url": "https://localhost:8081/anotherClient/",
      "country": "SWE",
      "calls": 2
   },
   ...
]
```

## Status endpoint
The status endpoint provides information about the services in the following format:
```
{
   "countries_api": "<http status code for *REST Countries API*>",
   "notification_db": "<http status code for *Notification DB* in Firebase>",
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```

# Retrieval
Our retrieval mechanism have two layers: cache and memory.

1. It checks if the input is nothing, meaning the user want all the countries to be written out. For the history handler it simply outputs the .csv file, but for the current handler it has a presaved datastructure with only current numbers which will be printed out. We decided to check this before caching, because it is not too costly (just an if statement) and may be an often used search.

2. If not the above it checks the cache (GET request). It does this by comparing the search URLs with the keys of the maps, which are the URL of the cached searches. If it finds a match it retrieves it and increases the "hit" variable by 1 (we will discuss what these means later). 

3. If it is not found it retrieves the data from memory, sends it to the user and saves it to the cache (SET request)

To ensure faster cache retrieval we have applied an algorithm which is based on number of hits (searches). If a country in the cache is being searched it hit counter increases by 1. If new data is being put on the stack it replaces the country with the least hits. We believe this is fitting for this application as the most popular is retrieved the fastest. To ensure that one search that is wildly popular at one time doesn't remain there forever as there can be loops of popularity. (A situation where it hit count is so high that the other countries are fighting instead and one search never will be dethroned). We have applied a purging mechanism that deletes cached data over a certain period (2 days). We find this period a good fit as it ensures fast retrieval for popular searches as well as resets to not have to high differences.