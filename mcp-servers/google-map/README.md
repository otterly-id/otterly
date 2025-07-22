# Google Maps MCP Server

A comprehensive Model Context Protocol (MCP) server that provides access to Google Maps services through various tools. This server enables AI applications to interact with Google Maps APIs for geocoding, places search, directions, and more.

## Features

This MCP server provides access to the following Google Maps services:

### 🗺️ Core Location Services
- **Geocoding**: Convert addresses to coordinates (latitude/longitude)
- **Reverse Geocoding**: Convert coordinates to human-readable addresses
- **Timezone Lookup**: Get timezone information for any location

### 🔍 Places & Search
- **Places Search**: Search for businesses and points of interest
- **Nearby Search**: Find places within a specified radius
- **Place Details**: Get detailed information about specific places (reviews, hours, contact info)

### 🚗 Navigation & Routing
- **Directions**: Get step-by-step directions between locations
- **Distance Matrix**: Calculate distances and travel times between multiple points
- **Multiple Travel Modes**: Support for driving, walking, transit, and bicycling

## Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd google-map
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

3. Set up your Google Maps API key:
```bash
export GOOGLE_MAPS_API_KEY="your_api_key_here"
```

## Google Maps API Setup

To use this MCP server, you need a Google Maps API key:

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the following APIs:
   - Maps JavaScript API
   - Geocoding API
   - Places API
   - Directions API
   - Distance Matrix API
   - Time Zone API
4. Create credentials (API Key)
5. Restrict the API key to the specific APIs above for security

### Required Google Maps APIs

- **Geocoding API**: For address-to-coordinates conversion
- **Places API**: For business and POI search
- **Directions API**: For route planning
- **Distance Matrix API**: For distance calculations
- **Time Zone API**: For timezone information

## Usage

Run the MCP server:

```bash
python src/google-map.py
```

## Available Tools

### 1. Geocoding
Convert an address or place name to coordinates.

**Parameters:**
- `query` (string): Address or place name to geocode

**Example:**
```json
{
  "query": "1600 Amphitheatre Parkway, Mountain View, CA"
}
```

### 2. Reverse Geocoding
Convert coordinates to a human-readable address.

**Parameters:**
- `latitude` (float): Latitude coordinate
- `longitude` (float): Longitude coordinate

**Example:**
```json
{
  "latitude": 37.4224082,
  "longitude": -122.0856086
}
```

### 3. Places Search
Search for places using text queries.

**Parameters:**
- `query` (string): Search query (e.g., "restaurants", "gas stations")
- `location` (string, optional): Location bias as "latitude,longitude"
- `radius` (int, optional): Search radius in meters (default: 5000)

**Example:**
```json
{
  "query": "pizza restaurants",
  "location": "37.4224082,-122.0856086",
  "radius": 2000
}
```

### 4. Get Directions
Get directions between two locations.

**Parameters:**
- `origin` (string): Starting location
- `destination` (string): Ending location
- `mode` (string): Travel mode (driving, walking, transit, bicycling)

**Example:**
```json
{
  "origin": "San Francisco, CA",
  "destination": "Los Angeles, CA",
  "mode": "driving"
}
```

### 5. Distance Matrix
Calculate distances between multiple origins and destinations.

**Parameters:**
- `origins` (array): List of origin locations
- `destinations` (array): List of destination locations
- `mode` (string): Travel mode (driving, walking, transit, bicycling)

**Example:**
```json
{
  "origins": ["San Francisco, CA", "Oakland, CA"],
  "destinations": ["Los Angeles, CA", "San Diego, CA"],
  "mode": "driving"
}
```

### 6. Nearby Search
Find places within a specified area.

**Parameters:**
- `location` (string): Center point as "latitude,longitude"
- `radius` (int): Search radius in meters (default: 1500, max: 50000)
- `place_type` (string, optional): Place type filter (e.g., "restaurant", "gas_station")

**Example:**
```json
{
  "location": "37.4224082,-122.0856086",
  "radius": 1000,
  "place_type": "restaurant"
}
```

### 7. Place Details
Get detailed information about a specific place.

**Parameters:**
- `place_id` (string): Unique place identifier from other searches

**Example:**
```json
{
  "place_id": "ChIJN1t_tDeuEmsRUsoyG83frY4"
}
```

### 8. Timezone Lookup
Get timezone information for a location.

**Parameters:**
- `location` (string): Coordinates as "latitude,longitude"
- `timestamp` (int, optional): UNIX timestamp (default: current time)

**Example:**
```json
{
  "location": "37.4224082,-122.0856086",
  "timestamp": 1609459200
}
```

## Security & Best Practices

1. **API Key Security**: Keep your Google Maps API key secure and never expose it in client-side code
2. **API Restrictions**: Restrict your API key to only the necessary Google Maps APIs
3. **Usage Monitoring**: Monitor your API usage in the Google Cloud Console
4. **Rate Limiting**: Be aware of Google Maps API rate limits and quotas

## Pricing

Google Maps APIs have usage-based pricing. You get $200 in free monthly usage, which covers:
- ~40,000 geocoding requests
- ~40,000 directions requests  
- ~100,000 places searches
- And more depending on the API

Check the [Google Maps Platform Pricing](https://developers.google.com/maps/billing/understanding-cost-of-use) for detailed information.

## Error Handling

The server includes comprehensive error handling for:
- Invalid API keys
- API quota exceeded
- Network errors
- Invalid parameters
- No results found

All errors are returned in a consistent format with descriptive messages.

## Testing

This project includes a comprehensive test suite that covers all tools and various edge cases.

### Running Tests

1. **Install test dependencies:**
   ```bash
   pip install -r test_requirements.txt
   ```

2. **Run all tests:**
   ```bash
   python run_tests.py
   ```

3. **Run tests with coverage:**
   ```bash
   python run_tests.py --coverage
   ```

### Test Coverage

The test suite includes:

- ✅ **Unit tests for all 8 tools**
- ✅ **Success scenario testing**
- ✅ **Error handling tests** (API errors, network errors, invalid inputs)
- ✅ **Edge case testing** (empty results, invalid modes, etc.)
- ✅ **Mock API responses** (no real API calls during testing)
- ✅ **Parameter validation**
- ✅ **HTTP error handling**

### Test Structure

- `test_google_map.py` - Main test suite
- `test_requirements.txt` - Test-specific dependencies
- `run_tests.py` - Test runner with coverage support

All tests use mocked HTTP responses, so you don't need a real Google Maps API key to run the tests.

## Contributing

Feel free to submit issues and enhancement requests!

When contributing:
1. Write tests for new features
2. Ensure all tests pass: `python run_tests.py`
3. Check code coverage: `python run_tests.py --coverage`
4. Follow the existing code style
