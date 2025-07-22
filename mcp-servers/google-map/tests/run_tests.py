#!/usr/bin/env python3
"""
Test runner for Google Maps MCP Server.

This script runs all tests with proper configuration and provides
clear output about test results.
"""

import sys
import subprocess
import os

def run_tests():
    """Run the test suite with proper configuration."""
    print("🧪 Running Google Maps MCP Server Tests")
    print("=" * 50)
    
    # Check if pytest is installed
    try:
        import pytest
    except ImportError:
        print("❌ Error: pytest not found!")
        print("Please install test requirements:")
        print("pip install -r test_requirements.txt")
        return False
    
    # Set up test environment
    test_env = os.environ.copy()
    test_env["GOOGLE_MAPS_API_KEY"] = "test_api_key"  # Mock API key for tests
    
    # Run tests with verbose output
    test_args = [
        sys.executable, "-m", "pytest",
        "src/test_google_map.py",
        "-v",  # Verbose output
        "--tb=short",  # Short traceback format
        "--color=yes",  # Colored output
        "-ra"  # Show summary of all results
    ]
    
    try:
        result = subprocess.run(test_args, env=test_env, check=True)
        print("\n✅ All tests passed!")
        return True
    except subprocess.CalledProcessError as e:
        print(f"\n❌ Tests failed with exit code {e.returncode}")
        return False
    except KeyboardInterrupt:
        print("\n⚠️  Tests interrupted by user")
        return False

def run_coverage():
    """Run tests with coverage reporting."""
    print("📊 Running Tests with Coverage")
    print("=" * 50)
    
    try:
        import coverage
    except ImportError:
        print("❌ Coverage not installed. Installing...")
        subprocess.run([sys.executable, "-m", "pip", "install", "coverage"], check=True)
    
    # Set up test environment
    test_env = os.environ.copy()
    test_env["GOOGLE_MAPS_API_KEY"] = "test_api_key"
    
    # Run tests with coverage
    coverage_args = [
        sys.executable, "-m", "coverage", "run",
        "-m", "pytest", "test_google_map.py", "-v"
    ]
    
    try:
        subprocess.run(coverage_args, env=test_env, check=True)
        
        # Generate coverage report
        print("\n📈 Coverage Report:")
        subprocess.run([sys.executable, "-m", "coverage", "report"], check=True)
        
        # Generate HTML coverage report
        subprocess.run([sys.executable, "-m", "coverage", "html"], check=True)
        print("\n📄 HTML coverage report generated: htmlcov/index.html")
        return True
        
    except subprocess.CalledProcessError as e:
        print(f"\n❌ Coverage run failed with exit code {e.returncode}")
        return False

def main():
    """Main test runner."""
    if len(sys.argv) > 1 and sys.argv[1] == "--coverage":
        success = run_coverage()
    else:
        success = run_tests()
    
    if not success:
        sys.exit(1)
    
    print("\n🎉 Testing completed successfully!")

if __name__ == "__main__":
    main() 