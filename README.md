# young_astrologer

## Installation
1. Rename the `.env_example` file to `.env`:
2. Build the Docker containers by running the following command:
   ```shell
   make compose-up

## Usage


````
Endpoints
Get All Astrology Data
Method: POST
URL: http://localhost:8080/api/get_all
Description: Retrieve all available astrology data.

Get Astrology Data by Day
Method: POST
URL: http://localhost:8080/api/get
Description: Retrieve astrology data for a specific day.

Get Astrology Image
Method: GET
URL: http://localhost:8080/image.jpg
Description: Get an astrology-related image.