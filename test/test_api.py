#!/usr/bin/env python3
"""
Python API Test Suite for Shop Project
快速测试脚本，无需启动 Go 服务，可独立运行
"""

# Prefer the 'requests' library but provide a small urllib-based fallback so
# the script can run even when 'requests' is not installed (fixes static
# analysis / missing dependency errors).
try:
    import requests
except Exception:
    import urllib.request
    import urllib.error
    import urllib.parse
    import json as _json

    class _Response:
        def __init__(self, status_code: int, text: str, content: bytes = b""):
            self.status_code = status_code
            self.text = text
            self.content = content

        def json(self):
            return _json.loads(self.text) if self.text else None

    class Session:
        def __init__(self):
            self.headers = {}

        def _request(self, method: str, url: str, json: dict = None):
            data = None
            headers = dict(self.headers or {})
            if json is not None:
                data = _json.dumps(json).encode("utf-8")
                headers.setdefault("Content-Type", "application/json")
            req = urllib.request.Request(url, data=data, method=method, headers=headers)
            try:
                with urllib.request.urlopen(req) as resp:
                    content = resp.read()
                    text = content.decode("utf-8") if content else ""
                    status = resp.getcode()
                    return _Response(status, text, content)
            except urllib.error.HTTPError as e:
                try:
                    text = e.read().decode("utf-8")
                except Exception:
                    text = str(e)
                return _Response(e.code if hasattr(e, "code") else 500, text, b"")
            except Exception as e:
                return _Response(0, str(e), b"")

        def get(self, url: str):
            return self._request("GET", url)

        def post(self, url: str, json: dict = None):
            return self._request("POST", url, json=json)

        def patch(self, url: str, json: dict = None):
            return self._request("PATCH", url, json=json)

        def delete(self, url: str, json: dict = None):
            return self._request("DELETE", url, json=json)

    # mimic the requests module interface used in this script
    requests = type("requests_module", (), {"Session": Session})

import json
import sys
from typing import Dict, Any, Optional

BASE_URL = "http://localhost:8080/api/v2"
ORDER_BASE_URL = "http://localhost:8080/api/v1"

class APITester:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})
        self.created_ids = {}
    
    def log(self, message: str, level: str = "INFO"):
        """Print formatted log message"""
        colors = {
            "INFO": "\033[94m",
            "SUCCESS": "\033[92m",
            "ERROR": "\033[91m",
            "WARNING": "\033[93m",
            "RESET": "\033[0m"
        }
        color = colors.get(level, colors["INFO"])
        print(f"{color}[{level}] {message}{colors['RESET']}")
    
    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None) -> Optional[Dict]:
        """Make HTTP request and return JSON response"""
        url = f"{self.base_url}{endpoint}"
        try:
            if method == "GET":
                response = self.session.get(url)
            elif method == "POST":
                response = self.session.post(url, json=data)
            elif method == "PATCH":
                response = self.session.patch(url, json=data)
            elif method == "DELETE":
                response = self.session.delete(url, json=data)
            else:
                self.log(f"Unknown method: {method}", "ERROR")
                return None
            
            if response.status_code in [200, 201]:
                self.log(f"{method} {endpoint} -> {response.status_code}", "SUCCESS")
                return response.json() if response.text else None
            else:
                self.log(f"{method} {endpoint} -> {response.status_code}: {response.text}", "ERROR")
                return None
        except Exception as e:
            self.log(f"Request failed: {e}", "ERROR")
            return None
    
    # ============= Shop Tests =============
    
    def test_create_shop(self) -> Optional[int]:
        """Create a test shop"""
        self.log("\n--- Creating Shop ---")
        data = {
            "name": "Python Test Shop",
            "description": "Created by Python test script",
            "owner_id": 1
        }
        result = self.make_request("POST", "/shops", data)
        if result and "ID" in result:
            shop_id = result["ID"]
            self.created_ids["shop"] = shop_id
            self.log(f"Shop created with ID: {shop_id}", "SUCCESS")
            return shop_id
        return None
    
    def test_list_shops(self):
        """List all shops"""
        self.log("\n--- Listing Shops ---")
        result = self.make_request("GET", "/shops")
        if result:
            self.log(f"Found shops", "SUCCESS")
            return result
        return None
    
    def test_get_shop(self, shop_id: int):
        """Get shop details"""
        self.log(f"\n--- Getting Shop {shop_id} ---")
        result = self.make_request("GET", f"/shops/{shop_id}")
        if result:
            self.log(f"Shop name: {result.get('name')}", "SUCCESS")
        return result
    
    def test_update_shop(self, shop_id: int):
        """Update shop"""
        self.log(f"\n--- Updating Shop {shop_id} ---")
        data = {
            "name": "Updated Python Test Shop",
            "description": "Updated by test script"
        }
        result = self.make_request("PATCH", f"/shops/{shop_id}", data)
        return result
    
    # ============= Product Tests =============
    
    def test_create_product(self, shop_id: int) -> Optional[int]:
        """Create a test product"""
        self.log(f"\n--- Creating Product in Shop {shop_id} ---")
        data = {
            "name": "Python Test Product",
            "description": "Test product from Python",
            "price": 99.99,
            "stock": 50,
            "product_img": "https://example.com/test-product.jpg"
        }
        result = self.make_request("POST", f"/shops/{shop_id}/products", data)
        if result and "ID" in result:
            product_id = result["ID"]
            self.created_ids["product"] = product_id
            self.log(f"Product created with ID: {product_id}", "SUCCESS")
            return product_id
        return None
    
    def test_list_products(self, shop_id: int):
        """List products in a shop"""
        self.log(f"\n--- Listing Products in Shop {shop_id} ---")
        result = self.make_request("GET", f"/shops/{shop_id}/products")
        if result:
            count = len(result) if isinstance(result, list) else 1
            self.log(f"Found {count} products", "SUCCESS")
        return result
    
    def test_search_product(self, shop_id: int, name: str):
        """Search product by name"""
        self.log(f"\n--- Searching Product: {name} in Shop {shop_id} ---")
        result = self.make_request("GET", f"/shops/{shop_id}/products/search?name={name}")
        if result:
            self.log(f"Search completed", "SUCCESS")
        return result
    
    def test_get_product(self, product_id: int):
        """Get product details"""
        self.log(f"\n--- Getting Product {product_id} ---")
        result = self.make_request("GET", f"/products/{product_id}")
        if result:
            self.log(f"Product name: {result.get('name')}", "SUCCESS")
        return result
    
    def test_update_product(self, product_id: int):
        """Update product"""
        self.log(f"\n--- Updating Product {product_id} ---")
        data = {
            "name": "Updated Python Test Product",
            "price": 149.99,
            "stock": 40
        }
        result = self.make_request("PATCH", f"/products/{product_id}", data)
        return result
    
    def test_delete_product(self, product_id: int):
        """Delete product"""
        self.log(f"\n--- Deleting Product {product_id} ---")
        result = self.make_request("DELETE", f"/products/{product_id}")
        return result
    
    # ============= Order Tests =============
    
    def test_create_order(self) -> Optional[int]:
        """Create a test order"""
        self.log("\n--- Creating Order ---")
        data = {
            "user_id": 999,
            "items": [
                {
                    "product_id": 1,
                    "product_name": "Test Product 1",
                    "price": 99.99,
                    "quantity": 2
                },
                {
                    "product_id": 2,
                    "product_name": "Test Product 2",
                    "price": 49.99,
                    "quantity": 1
                }
            ]
        }
        result = self.make_request("POST", "/orders", data, base_url=ORDER_BASE_URL)
        if result and "ID" in result:
            order_id = result["ID"]
            self.created_ids["order"] = order_id
            self.log(f"Order created with ID: {order_id}", "SUCCESS")
            return order_id
        return None
    
    def test_list_orders(self):
        """List all orders"""
        self.log("\n--- Listing Orders ---")
        result = self.make_request("GET", "/orders", base_url=ORDER_BASE_URL)
        if result:
            count = len(result) if isinstance(result, list) else 1
            self.log(f"Found {count} orders", "SUCCESS")
        return result
    
    def test_get_order(self, order_id: int):
        """Get order details"""
        self.log(f"\n--- Getting Order {order_id} ---")
        result = self.make_request("GET", f"/orders/{order_id}", base_url=ORDER_BASE_URL)
        if result:
            self.log(f"Order status: {result.get('status')}", "SUCCESS")
        return result
    
    def test_update_order_status(self, order_id: int):
        """Update order status"""
        self.log(f"\n--- Updating Order {order_id} Status ---")
        data = {"status": "shipped"}
        result = self.make_request("PATCH", f"/orders/{order_id}/status", data, base_url=ORDER_BASE_URL)
        return result
    
    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None, base_url: str = None) -> Optional[Dict]:
        """Make HTTP request with optional custom base URL"""
        url = f"{base_url or self.base_url}{endpoint}"
        try:
            if method == "GET":
                response = self.session.get(url)
            elif method == "POST":
                response = self.session.post(url, json=data)
            elif method == "PATCH":
                response = self.session.patch(url, json=data)
            elif method == "DELETE":
                response = self.session.delete(url, json=data)
            else:
                self.log(f"Unknown method: {method}", "ERROR")
                return None
            
            if response.status_code in [200, 201]:
                self.log(f"{method} {endpoint} -> {response.status_code}", "SUCCESS")
                return response.json() if response.text else None
            else:
                self.log(f"{method} {endpoint} -> {response.status_code}: {response.text}", "ERROR")
                return None
        except Exception as e:
            self.log(f"Request failed: {e}", "ERROR")
            return None
    
    def run_all_tests(self):
        """Run all test cases"""
        print("\n" + "="*50)
        print("      API Integration Tests (Python)")
        print("="*50)
        
        # Test Shop
        shop_id = self.test_create_shop()
        if shop_id:
            self.test_list_shops()
            self.test_get_shop(shop_id)
            self.test_update_shop(shop_id)
            
            # Test Product
            product_id = self.test_create_product(shop_id)
            if product_id:
                self.test_list_products(shop_id)
                self.test_search_product(shop_id, "Python")
                self.test_get_product(product_id)
                self.test_update_product(product_id)
                # Don't delete yet, might need it for orders
        
        # Test Order
        order_id = self.test_create_order()
        if order_id:
            self.test_list_orders()
            self.test_get_order(order_id)
            self.test_update_order_status(order_id)
        
        print("\n" + "="*50)
        print("      All Tests Completed")
        print("="*50)
        print(f"\nCreated IDs: {json.dumps(self.created_ids, indent=2)}")


def main():
    """Main entry point"""
    tester = APITester(BASE_URL)
    try:
        tester.run_all_tests()
    except KeyboardInterrupt:
        print("\n\nTests interrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"\nFatal error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()

# test_api.py 示例
shop_data = {
    "name": "Apple Store",
    "description": "Official Apple Products",
    "owner_id": 1,
    "products": [
        {"name": "iPhone 15", "description": "Latest iPhone", "price": 999.99, "stock": 50},
        {"name": "MacBook Pro", "description": "Powerful laptop", "price": 1999.99, "stock": 20}
    ]
}

# 查询测试
search_url = "http://localhost:8080/api/v2/shops/1/products/search?name=iPhone%2015"
