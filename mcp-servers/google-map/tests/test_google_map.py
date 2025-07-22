"""
Comprehensive test suite for Otterly's Google Maps MCP Server.

This test suite covers all tools and handles various scenarios including
successful API calls, error cases, and edge cases.
"""

import pytest
import httpx
import os
import sys
from unittest.mock import patch, AsyncMock
import importlib.util

# Import the google-map module
spec = importlib.util.spec_from_file_location("google_map", "src/google-map.py")
if spec is None:
    raise ImportError("Failed to load google_map module")
if spec.loader is None:
    raise ImportError("Failed to load google_map module")

google_map = importlib.util.module_from_spec(spec)
sys.modules["google_map"] = google_map
spec.loader.exec_module(google_map)

# Test fixtures and mock data
@pytest.fixture
def mock_api_key():
    """Mock API key for testing."""
    with patch.dict(os.environ, {"GOOGLE_MAPS_API_KEY": "test_api_key"}):
        yield "test_api_key"

@pytest.fixture
def sample_geocode_response():
    """Sample successful geocoding API response."""
    return {
        "status": "OK",
        "results": [
            {
                "formatted_address": "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA",
                "geometry": {
                    "location": {
                        "lat": 37.4224082,
                        "lng": -122.0856086
                    }
                },
                "place_id": "ChIJtYuu0V25j4ARwu5e4wwRYgE",
                "types": ["street_address"]
            }
        ]
    }

@pytest.fixture
def sample_places_response():
    """Sample successful places search API response."""
    return {
        "status": "OK",
        "results": [
            {
                "name": "Tony's Little Star Pizza",
                "formatted_address": "846 Divisadero St, San Francisco, CA 94117, USA",
                "rating": 4.5,
                "place_id": "ChIJtYuu0V25j4ARwu5e4wwRYgE",
                "types": ["restaurant", "food", "point_of_interest"],
                "geometry": {
                    "location": {
                        "lat": 37.776,
                        "lng": -122.437
                    }
                }
            }
        ]
    }

@pytest.fixture
def sample_directions_response():
    """Sample successful directions API response."""
    return {
        "status": "OK",
        "routes": [
            {
                "legs": [
                    {
                        "start_address": "San Francisco, CA, USA",
                        "end_address": "Los Angeles, CA, USA",
                        "distance": {"text": "383 mi", "value": 617000},
                        "duration": {"text": "6 hours 30 mins", "value": 23400},
                        "steps": [
                            {
                                "html_instructions": "Head south on US-101 S",
                                "distance": {"text": "50 mi", "value": 80467},
                                "duration": {"text": "1 hour", "value": 3600}
                            }
                        ]
                    }
                ],
                "overview_polyline": {"points": "encoded_polyline_data"}
            }
        ]
    }

@pytest.fixture
def sample_distance_matrix_response():
    """Sample successful distance matrix API response."""
    return {
        "status": "OK",
        "origin_addresses": ["San Francisco, CA, USA"],
        "destination_addresses": ["Los Angeles, CA, USA"],
        "rows": [
            {
                "elements": [
                    {
                        "status": "OK",
                        "distance": {"text": "383 mi", "value": 617000},
                        "duration": {"text": "6 hours 30 mins", "value": 23400}
                    }
                ]
            }
        ]
    }

@pytest.fixture
def sample_timezone_response():
    """Sample successful timezone API response."""
    return {
        "status": "OK",
        "timeZoneId": "America/Los_Angeles",
        "timeZoneName": "Pacific Standard Time",
        "dstOffset": 0,
        "rawOffset": -28800
    }

class TestGoogleMapServer:
    """Test cases for Google Maps MCP Server."""

    def test_api_key_missing(self):
        """Test that missing API key raises ValueError."""
        with patch.dict(os.environ, {}, clear=True):
            with pytest.raises(ValueError, match="GOOGLE_MAPS_API_KEY environment variable is required"):
                google_map.get_api_key()

    def test_api_key_present(self, mock_api_key):
        """Test that API key is correctly retrieved."""
        assert google_map.get_api_key() == "test_api_key"

    @pytest.mark.asyncio
    async def test_geocode_success(self, mock_api_key, sample_geocode_response):
        """Test successful geocoding."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_geocode_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.geocode("1600 Amphitheatre Parkway, Mountain View, CA")
            
            assert result["formatted_address"] == "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA"
            assert result["latitude"] == 37.4224082
            assert result["longitude"] == -122.0856086
            assert result["place_id"] == "ChIJtYuu0V25j4ARwu5e4wwRYgE"

    @pytest.mark.asyncio
    async def test_geocode_no_results(self, mock_api_key):
        """Test geocoding with no results."""
        response = {"status": "OK", "results": []}
        
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.geocode("NonexistentPlace12345")
            
            assert "error" in result
            assert result["error"] == "No results found"

    @pytest.mark.asyncio
    async def test_geocode_api_error(self, mock_api_key):
        """Test geocoding with API error."""
        response = {"status": "INVALID_REQUEST", "error_message": "Invalid request"}
        
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.geocode("")
            
            assert "error" in result
            assert "Invalid request" in result["error"]

    @pytest.mark.asyncio
    async def test_reverse_geocode_success(self, mock_api_key, sample_geocode_response):
        """Test successful reverse geocoding."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_geocode_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.reverse_geocode(37.4224082, -122.0856086)
            
            assert result["formatted_address"] == "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA"
            assert result["place_id"] == "ChIJtYuu0V25j4ARwu5e4wwRYgE"

    @pytest.mark.asyncio
    async def test_places_search_success(self, mock_api_key, sample_places_response):
        """Test successful places search."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_places_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.places_search("pizza restaurants")
            
            assert "results" in result
            assert len(result["results"]) == 1
            assert result["results"][0]["name"] == "Tony's Little Star Pizza"
            assert result["results"][0]["rating"] == 4.5

    @pytest.mark.asyncio
    async def test_places_search_with_location(self, mock_api_key, sample_places_response):
        """Test places search with location bias."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_places_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.places_search(
                "restaurants", 
                location="37.7749,-122.4194", 
                radius=1000
            )
            
            assert "results" in result
            assert result["results"][0]["name"] == "Tony's Little Star Pizza"

    @pytest.mark.asyncio
    async def test_get_directions_success(self, mock_api_key, sample_directions_response):
        """Test successful directions request."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_directions_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.get_directions("San Francisco, CA", "Los Angeles, CA")
            
            assert result["start_address"] == "San Francisco, CA, USA"
            assert result["end_address"] == "Los Angeles, CA, USA"
            assert result["distance"] == "383 mi"
            assert result["duration"] == "6 hours 30 mins"
            assert len(result["steps"]) == 1

    @pytest.mark.asyncio
    async def test_get_directions_invalid_mode(self, mock_api_key):
        """Test directions with invalid travel mode."""
        result = await google_map.get_directions(
            "San Francisco, CA", 
            "Los Angeles, CA", 
            mode="invalid_mode"
        )
        
        assert "error" in result
        assert "Invalid mode" in result["error"]

    @pytest.mark.asyncio
    async def test_distance_matrix_success(self, mock_api_key, sample_distance_matrix_response):
        """Test successful distance matrix request."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_distance_matrix_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.distance_matrix(
                ["San Francisco, CA"], 
                ["Los Angeles, CA"]
            )
            
            assert "results" in result
            assert len(result["results"]) == 1
            assert result["results"][0]["origin"] == "San Francisco, CA, USA"
            assert result["results"][0]["destination"] == "Los Angeles, CA, USA"
            assert result["results"][0]["distance"] == "383 mi"

    @pytest.mark.asyncio
    async def test_distance_matrix_invalid_mode(self, mock_api_key):
        """Test distance matrix with invalid travel mode."""
        result = await google_map.distance_matrix(
            ["San Francisco, CA"], 
            ["Los Angeles, CA"], 
            mode="invalid_mode"
        )
        
        assert "error" in result
        assert "Invalid mode" in result["error"]

    @pytest.mark.asyncio
    async def test_nearby_search_success(self, mock_api_key, sample_places_response):
        """Test successful nearby search."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_places_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.nearby_search("37.7749,-122.4194", radius=1000)
            
            assert "results" in result
            assert result["results"][0]["name"] == "Tony's Little Star Pizza"

    @pytest.mark.asyncio
    async def test_nearby_search_with_type(self, mock_api_key, sample_places_response):
        """Test nearby search with place type filter."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_places_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.nearby_search(
                "37.7749,-122.4194", 
                radius=1000, 
                place_type="restaurant"
            )
            
            assert "results" in result

    @pytest.mark.asyncio
    async def test_place_details_success(self, mock_api_key):
        """Test successful place details request."""
        details_response = {
            "status": "OK",
            "result": {
                "name": "Tony's Little Star Pizza",
                "formatted_address": "846 Divisadero St, San Francisco, CA 94117, USA",
                "formatted_phone_number": "(415) 441-1118",
                "rating": 4.5,
                "website": "http://www.tonylittlestar.com/",
                "opening_hours": {"open_now": True},
                "reviews": [{"rating": 5, "text": "Great pizza!"}],
                "geometry": {"location": {"lat": 37.776, "lng": -122.437}}
            }
        }
        
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = details_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.place_details("ChIJtYuu0V25j4ARwu5e4wwRYgE")
            
            assert result["name"] == "Tony's Little Star Pizza"
            assert result["phone_number"] == "(415) 441-1118"
            assert result["rating"] == 4.5
            assert result["website"] == "http://www.tonylittlestar.com/"

    @pytest.mark.asyncio
    async def test_timezone_lookup_success(self, mock_api_key, sample_timezone_response):
        """Test successful timezone lookup."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_timezone_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.timezone_lookup("37.4224082,-122.0856086")
            
            assert result["timezone_id"] == "America/Los_Angeles"
            assert result["timezone_name"] == "Pacific Standard Time"
            assert result["dst_offset"] == 0
            assert result["raw_offset"] == -28800

    @pytest.mark.asyncio
    async def test_timezone_lookup_with_timestamp(self, mock_api_key, sample_timezone_response):
        """Test timezone lookup with custom timestamp."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_timezone_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.timezone_lookup("37.4224082,-122.0856086", timestamp=1609459200)
            
            assert result["timezone_id"] == "America/Los_Angeles"

    @pytest.mark.asyncio
    async def test_http_error_handling(self, mock_api_key):
        """Test HTTP error handling."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.raise_for_status.side_effect = httpx.HTTPStatusError(
                "Client error", request=httpx.Request("GET", "https://example.com"), response=httpx.Response(status_code=400)
            )
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            with pytest.raises(httpx.HTTPStatusError):
                await google_map.geocode("test address")

    @pytest.mark.asyncio
    async def test_distance_matrix_element_error(self, mock_api_key):
        """Test distance matrix with element-level errors."""
        error_response = {
            "status": "OK",
            "origin_addresses": ["San Francisco, CA, USA"],
            "destination_addresses": ["InvalidDestination"],
            "rows": [
                {
                    "elements": [
                        {"status": "NOT_FOUND"}
                    ]
                }
            ]
        }
        
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = error_response
            mock_response.raise_for_status.return_value = None
            mock_client.return_value.__aenter__.return_value.get.return_value = mock_response

            result = await google_map.distance_matrix(
                ["San Francisco, CA"], 
                ["InvalidDestination"]
            )
            
            assert "results" in result
            assert result["results"][0]["error"] == "NOT_FOUND"

    @pytest.mark.asyncio
    async def test_nearby_search_radius_limit(self, mock_api_key, sample_places_response):
        """Test that nearby search respects maximum radius limit."""
        with patch('httpx.AsyncClient') as mock_client:
            mock_response = AsyncMock()
            mock_response.json.return_value = sample_places_response
            mock_response.raise_for_status.return_value = None
            
            # Capture the actual request parameters
            captured_params = {}
            async def capture_get(url, params=None):
                captured_params.update(params or {})
                return mock_response
            
            mock_client.return_value.__aenter__.return_value.get = capture_get

            # Test with radius larger than max (50000)
            await google_map.nearby_search("37.7749,-122.4194", radius=100000)
            
            # Should be limited to 50000
            assert captured_params["radius"] == "50000"


if __name__ == "__main__":
    pytest.main([__file__, "-v"]) 