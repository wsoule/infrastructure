const API_GATEWAY_URL = import.meta.env.VITE_API_GATEWAY_URL || 'http://localhost:8080';

// Auto-load data when page loads
window.addEventListener('DOMContentLoaded', () => {
  loadHealthStatus();
  loadUsers();
  loadProducts();
});

// Generic API call function
async function apiCall(endpoint, method = 'GET', body = null) {
  try {
    const options = {
      method,
      headers: { 'Content-Type': 'application/json' }
    };
    if (body) {
      options.body = JSON.stringify(body);
    }

    const response = await fetch(`${API_GATEWAY_URL}${endpoint}`, options);
    const data = await response.json();
    return { success: response.ok, status: response.status, data };
  } catch (error) {
    return { success: false, error: error.message };
  }
}

// Health Check
async function loadHealthStatus() {
  const result = await apiCall('/health');
  const el = document.getElementById('health-status');
  if (result.success) {
    el.innerHTML = '<strong>Gateway Status:</strong> ' + result.data.status + '<br>' +
                   '<strong>Services:</strong><pre>' + JSON.stringify(result.data.services, null, 2) + '</pre>';
  } else {
    el.innerHTML = '<strong>Error:</strong> ' + (result.error || 'Failed to check health');
  }
}

// Users
async function loadUsers() {
  const result = await apiCall('/api/users');
  const el = document.getElementById('users-list');
  if (result.success) {
    if (result.data && result.data.length > 0) {
      el.innerHTML = '<pre>' + JSON.stringify(result.data, null, 2) + '</pre>';
    } else {
      el.innerHTML = '<em>No users found</em>';
    }
  } else {
    el.innerHTML = '<strong>Error:</strong> ' + (result.error || JSON.stringify(result.data));
  }
}

document.getElementById('create-user-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);
  const user = {
    name: formData.get('name'),
    email: formData.get('email')
  };
  const result = await apiCall('/api/users', 'POST', user);
  const el = document.getElementById('create-user-response');
  if (result.success) {
    el.innerHTML = '<strong>Success:</strong><pre>' + JSON.stringify(result.data, null, 2) + '</pre>';
    e.target.reset();
    loadUsers(); // Reload the list
  } else {
    el.innerHTML = '<strong>Error:</strong> ' + (result.error || JSON.stringify(result.data));
  }
});

document.getElementById('get-user-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);
  const id = formData.get('id');
  const result = await apiCall(`/api/users/${id}`);
  const el = document.getElementById('get-user-response');
  if (result.success) {
    el.innerHTML = '<strong>User #' + id + ':</strong><pre>' + JSON.stringify(result.data, null, 2) + '</pre>';
  } else {
    el.innerHTML = '<strong>Error:</strong> ' + (result.error || JSON.stringify(result.data));
  }
});

// Products
async function loadProducts() {
  const result = await apiCall('/api/products');
  const el = document.getElementById('products-list');
  if (result.success) {
    if (result.data && result.data.length > 0) {
      el.innerHTML = '<pre>' + JSON.stringify(result.data, null, 2) + '</pre>';
    } else {
      el.innerHTML = '<em>No products found</em>';
    }
  } else {
    el.innerHTML = '<strong>Error:</strong> ' + (result.error || JSON.stringify(result.data));
  }
}

document.getElementById('create-product-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);
  const product = {
    name: formData.get('name'),
    description: formData.get('description'),
    price: parseFloat(formData.get('price')),
    stock: parseInt(formData.get('stock'))
  };
  const result = await apiCall('/api/products', 'POST', product);
  const el = document.getElementById('create-product-response');
  if (result.success) {
    el.innerHTML = '<strong>Success:</strong><pre>' + JSON.stringify(result.data, null, 2) + '</pre>';
    e.target.reset();
    loadProducts(); // Reload the list
  } else {
    el.innerHTML = '<strong>Error:</strong> ' + (result.error || JSON.stringify(result.data));
  }
});

document.getElementById('get-product-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);
  const id = formData.get('id');
  const result = await apiCall(`/api/products/${id}`);
  const el = document.getElementById('get-product-response');
  if (result.success) {
    el.innerHTML = '<strong>Product #' + id + ':</strong><pre>' + JSON.stringify(result.data, null, 2) + '</pre>';
  } else {
    el.innerHTML = '<strong>Error:</strong> ' + (result.error || JSON.stringify(result.data));
  }
});
