#!/usr/bin/env python3
"""
Example script demonstrating Google Maps MCP Server usage.

This script shows how to use various Google Maps services through the MCP server.
Make sure to set your GOOGLE_MAPS_API_KEY environment variable before running.
"""

import asyncio
import os
import sys
import os
sys.path.append(os.path.join(os.path.dirname(__file__), 'src'))

# Import from the google-map module
import importlib.util
spec = importlib.util.spec_from_file_location("google_map", "src/google-map.py")
google_map = importlib.util.module_from_spec(spec)
spec.loader.exec_module(google_map)

# Import the functions we need
geocode = google_map.geocode
reverse_geocode = google_map.reverse_geocode
places_search = google_map.places_search
get_directions = google_map.get_directions
distance_matrix = google_map.distance_matrix
nearby_search = google_map.nearby_search
place_details = google_map.place_details
timezone_lookup = google_map.timezone_lookup

async def demo_geocoding():
    """Demonstrate geocoding functionality."""
    print("🗺️  GEOCODING DEMO")
    print("=" * 50)
    
    # Geocode an address
    result = await geocode("1600 Amphitheatre Parkway, Mountain View, CA")
    print(f"Geocoding result: {result}")
    
    if "latitude" in result and "longitude" in result:
        lat, lng = result["latitude"], result["longitude"]
        
        # Reverse geocode the coordinates
        reverse_result = await reverse_geocode(lat, lng)
        print(f"Reverse geocoding result: {reverse_result}")
    
    print("\n")

async def demo_places():
    """Demonstrate places functionality."""
    print("🔍 PLACES DEMO")
    print("=" * 50)
    
    # Search for restaurants
    places_result = await places_search("pizza restaurants", location="37.4224082,-122.0856086", radius=2000)
    print(f"Places search result: {places_result}")
    
    # Get nearby places
    nearby_result = await nearby_search("37.4224082,-122.0856086", radius=1000, place_type="restaurant")
    print(f"Nearby search result: {nearby_result}")
    
    # Get details for the first place if available
    if "results" in places_result and places_result["results"]:
        place_id = places_result["results"][0].get("place_id")
        if place_id:
            details_result = await place_details(place_id)
            print(f"Place details result: {details_result}")
    
    print("\n")

async def demo_directions():
    """Demonstrate directions and distance functionality."""
    print("🚗 DIRECTIONS & DISTANCE DEMO")
    print("=" * 50)
    
    # Get directions
    directions_result = await get_directions(
        "San Francisco, CA", 
        "Los Angeles, CA", 
        mode="driving"
    )
    print(f"Directions result: {directions_result}")
    
    # Calculate distance matrix
    distance_result = await distance_matrix(
        origins=["San Francisco, CA", "Oakland, CA"],
        destinations=["Los Angeles, CA", "San Diego, CA"],
        mode="driving"
    )
    print(f"Distance matrix result: {distance_result}")
    
    print("\n")

async def demo_timezone():
    """Demonstrate timezone functionality."""
    print("🕒 TIMEZONE DEMO")
    print("=" * 50)
    
    # Get timezone for coordinates
    timezone_result = await timezone_lookup("37.4224082,-122.0856086")
    print(f"Timezone result: {timezone_result}")
    
    print("\n")

async def main():
    """Run all demonstrations."""
    print("🚀 Google Maps MCP Server Demo")
    print("=" * 60)
    
    # Check if API key is set
    if not os.getenv("GOOGLE_MAPS_API_KEY"):
        print("❌ Error: GOOGLE_MAPS_API_KEY environment variable not set!")
        print("Please set your Google Maps API key:")
        print("export GOOGLE_MAPS_API_KEY='your_api_key_here'")
        return
    
    print("✅ API key found, running demos...\n")
    
    try:
        await demo_geocoding()
        await demo_places()
        await demo_directions()
        await demo_timezone()
        
        print("🎉 All demos completed successfully!")
        
    except Exception as e:
        print(f"❌ Error during demo: {e}")
        print("Make sure your API key is valid and the required APIs are enabled.")

if __name__ == "__main__":
    asyncio.run(main()) 