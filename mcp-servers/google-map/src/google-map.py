import httpx
import os
from typing import Optional, Dict, Any, List
from mcp.server.fastmcp import FastMCP
from mcp.types import Tool, TextContent

# Initialize the MCP server
mcp = FastMCP("google-maps")

# Google Maps API base URL
GOOGLE_MAPS_BASE_URL = "https://maps.googleapis.com/maps/api"

def get_api_key() -> str:
    """Get Google Maps API key from environment variable."""
    api_key = os.getenv("GOOGLE_MAPS_API_KEY")
    if not api_key:
        raise ValueError("GOOGLE_MAPS_API_KEY environment variable is required")
    return api_key

@mcp.tool()
async def geocode(query: str) -> Dict[str, Any]:
    """
    Convert an address or place name into geographic coordinates (latitude/longitude).
    
    Args:
        query: The address or place name to geocode
        
    Returns:
        Dictionary containing geocoding results with coordinates and formatted address
    """
    api_key = get_api_key()
    
    params = {
        "address": query,
        "key": api_key
    }
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/geocode/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Geocoding failed: {data.get('error_message', data['status'])}"}
        
        if not data["results"]:
            return {"error": "No results found"}
        
        result = data["results"][0]
        return {
            "formatted_address": result["formatted_address"],
            "latitude": result["geometry"]["location"]["lat"],
            "longitude": result["geometry"]["location"]["lng"],
            "place_id": result.get("place_id"),
            "types": result.get("types", [])
        }

@mcp.tool()
async def reverse_geocode(latitude: float, longitude: float) -> Dict[str, Any]:
    """
    Convert geographic coordinates into a human-readable address.
    
    Args:
        latitude: The latitude coordinate
        longitude: The longitude coordinate
        
    Returns:
        Dictionary containing reverse geocoding results with address information
    """
    api_key = get_api_key()
    
    params = {
        "latlng": f"{latitude},{longitude}",
        "key": api_key
    }
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/geocode/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Reverse geocoding failed: {data.get('error_message', data['status'])}"}
        
        if not data["results"]:
            return {"error": "No results found"}
        
        result = data["results"][0]
        return {
            "formatted_address": result["formatted_address"],
            "address_components": result.get("address_components", []),
            "place_id": result.get("place_id"),
            "types": result.get("types", [])
        }

@mcp.tool()
async def places_search(query: str, location: Optional[str] = None, radius: Optional[int] = 5000) -> Dict[str, Any]:
    """
    Search for places using Google Places API.
    
    Args:
        query: The search query (e.g., "restaurants", "gas stations")
        location: Optional location bias as "latitude,longitude"
        radius: Optional search radius in meters (default: 5000)
        
    Returns:
        Dictionary containing places search results
    """
    api_key = get_api_key()
    
    params = {
        "query": query,
        "key": api_key
    }
    
    if location:
        params["location"] = location
        params["radius"] = str(radius)
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/place/textsearch/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Places search failed: {data.get('error_message', data['status'])}"}
        
        results = []
        for place in data.get("results", [])[:10]:  # Limit to first 10 results
            results.append({
                "name": place.get("name"),
                "formatted_address": place.get("formatted_address"),
                "rating": place.get("rating"),
                "place_id": place.get("place_id"),
                "types": place.get("types", []),
                "geometry": place.get("geometry", {}).get("location", {})
            })
        
        return {"results": results, "total_results": len(data.get("results", []))}

@mcp.tool()
async def get_directions(origin: str, destination: str, mode: str = "driving") -> Dict[str, Any]:
    """
    Get directions between two locations.
    
    Args:
        origin: Starting location (address or coordinates)
        destination: Ending location (address or coordinates)
        mode: Travel mode (driving, walking, transit, bicycling)
        
    Returns:
        Dictionary containing directions and route information
    """
    api_key = get_api_key()
    
    valid_modes = ["driving", "walking", "transit", "bicycling"]
    if mode not in valid_modes:
        return {"error": f"Invalid mode. Must be one of: {', '.join(valid_modes)}"}
    
    params = {
        "origin": origin,
        "destination": destination,
        "mode": mode,
        "key": api_key
    }
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/directions/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Directions failed: {data.get('error_message', data['status'])}"}
        
        if not data["routes"]:
            return {"error": "No routes found"}
        
        route = data["routes"][0]
        leg = route["legs"][0]
        
        return {
            "start_address": leg["start_address"],
            "end_address": leg["end_address"],
            "distance": leg["distance"]["text"],
            "duration": leg["duration"]["text"],
            "steps": [
                {
                    "instruction": step["html_instructions"],
                    "distance": step["distance"]["text"],
                    "duration": step["duration"]["text"]
                }
                for step in leg["steps"]
            ],
            "overview_polyline": route["overview_polyline"]["points"]
        }

@mcp.tool()
async def distance_matrix(origins: List[str], destinations: List[str], mode: str = "driving") -> Dict[str, Any]:
    """
    Calculate distance and time between multiple origins and destinations.
    
    Args:
        origins: List of origin locations
        destinations: List of destination locations
        mode: Travel mode (driving, walking, transit, bicycling)
        
    Returns:
        Dictionary containing distance matrix results
    """
    api_key = get_api_key()
    
    valid_modes = ["driving", "walking", "transit", "bicycling"]
    if mode not in valid_modes:
        return {"error": f"Invalid mode. Must be one of: {', '.join(valid_modes)}"}
    
    params = {
        "origins": "|".join(origins),
        "destinations": "|".join(destinations),
        "mode": mode,
        "key": api_key
    }
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/distancematrix/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Distance matrix failed: {data.get('error_message', data['status'])}"}
        
        results = []
        for i, row in enumerate(data["rows"]):
            for j, element in enumerate(row["elements"]):
                if element["status"] == "OK":
                    results.append({
                        "origin": data["origin_addresses"][i],
                        "destination": data["destination_addresses"][j],
                        "distance": element["distance"]["text"],
                        "duration": element["duration"]["text"]
                    })
                else:
                    results.append({
                        "origin": data["origin_addresses"][i],
                        "destination": data["destination_addresses"][j],
                        "error": element["status"]
                    })
        
        return {"results": results}

@mcp.tool()
async def nearby_search(location: str, radius: int = 1500, place_type: Optional[str] = None) -> Dict[str, Any]:
    """
    Find places within a specified area.
    
    Args:
        location: The latitude/longitude around which to search (format: "lat,lng")
        radius: The search radius in meters (default: 1500, max: 50000)
        place_type: Optional place type to filter results (e.g., "restaurant", "gas_station")
        
    Returns:
        Dictionary containing nearby places
    """
    api_key = get_api_key()
    
    params = {
        "location": location,
        "radius": min(radius, 50000),  # Max radius is 50km
        "key": api_key
    }
    
    if place_type:
        params["type"] = place_type
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/place/nearbysearch/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Nearby search failed: {data.get('error_message', data['status'])}"}
        
        results = []
        for place in data.get("results", [])[:20]:  # Limit to first 20 results
            results.append({
                "name": place.get("name"),
                "place_id": place.get("place_id"),
                "rating": place.get("rating"),
                "vicinity": place.get("vicinity"),
                "types": place.get("types", []),
                "geometry": place.get("geometry", {}).get("location", {}),
                "price_level": place.get("price_level"),
                "opening_hours": place.get("opening_hours", {}).get("open_now")
            })
        
        return {"results": results, "total_results": len(data.get("results", []))}

@mcp.tool()
async def place_details(place_id: str) -> Dict[str, Any]:
    """
    Get detailed information about a specific place.
    
    Args:
        place_id: The unique identifier for a place
        
    Returns:
        Dictionary containing detailed place information
    """
    api_key = get_api_key()
    
    params = {
        "place_id": place_id,
        "fields": "name,formatted_address,formatted_phone_number,rating,opening_hours,website,reviews,photos,geometry",
        "key": api_key
    }
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/place/details/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Place details failed: {data.get('error_message', data['status'])}"}
        
        result = data.get("result", {})
        return {
            "name": result.get("name"),
            "formatted_address": result.get("formatted_address"),
            "phone_number": result.get("formatted_phone_number"),
            "rating": result.get("rating"),
            "website": result.get("website"),
            "opening_hours": result.get("opening_hours"),
            "reviews": result.get("reviews", [])[:5],  # Limit to first 5 reviews
            "geometry": result.get("geometry", {}).get("location", {})
        }

@mcp.tool()
async def timezone_lookup(location: str, timestamp: Optional[int] = None) -> Dict[str, Any]:
    """
    Get timezone information for a location.
    
    Args:
        location: The latitude/longitude (format: "lat,lng")
        timestamp: Optional UNIX timestamp (default: current time)
        
    Returns:
        Dictionary containing timezone information
    """
    api_key = get_api_key()
    
    import time
    if timestamp is None:
        timestamp = int(time.time())
    
    params = {
        "location": location,
        "timestamp": timestamp,
        "key": api_key
    }
    
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{GOOGLE_MAPS_BASE_URL}/timezone/json", params=params)
        response.raise_for_status()
        data = response.json()
        
        if data["status"] != "OK":
            return {"error": f"Timezone lookup failed: {data.get('error_message', data['status'])}"}
        
        return {
            "timezone_id": data.get("timeZoneId"),
            "timezone_name": data.get("timeZoneName"),
            "dst_offset": data.get("dstOffset"),
            "raw_offset": data.get("rawOffset")
        }

if __name__ == "__main__":
    # Run the MCP server
    mcp.run()