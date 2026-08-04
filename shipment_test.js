#!/usr/bin/env node

const baseURL = process.env.BASE_URL || 'http://localhost:8081';

async function request(path, options = {}) {
  const url = `${baseURL}${path}`;
  const res = await fetch(url, options);
  const body = await res.text();
  let data;
  try {
    data = body ? JSON.parse(body) : null;
  } catch (err) {
    data = body;
  }
  return { status: res.status, data };
}

async function main() {
  console.log('Base URL:', baseURL);

  const health = await request('/');
  console.log('GET / =>', health.status, health.data);
  if (health.status !== 200) {
    throw new Error('Server health check failed');
  }

  const shipmentPayload = {
    sender: {
      senderName: 'Alice Example',
      senderPhone: '+1234567890',
      pickupLocation: '100 Main St',
      pickupCity: 'Sample City',
      pickupCountry: 'Exampleland'
    },
    receiver: {
      receiverName: 'Bob Receiver',
      receiverPhone: '+0987654321',
      deliveryLocation: '200 Market Ave',
      deliveryCity: 'Destination City',
      deliveryCountry: 'Dreamland'
    },
    package: {
      packageDescription: 'Test package',
      packageWeight: 1,
      packageCategories: 'Test'
    }
  };

  const createResp = await request('/shipments', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(shipmentPayload),
  });
  console.log('POST /shipments =>', createResp.status, createResp.data);
  if (createResp.status !== 201 || !createResp.data?.id) {
    throw new Error('Shipment creation failed');
  }

  const shipmentId = createResp.data.id;

  const getResp = await request(`/shipments/${shipmentId}`);
  console.log(`GET /shipments/${shipmentId} =>`, getResp.status, getResp.data);
  if (getResp.status !== 200) {
    throw new Error('Fetch shipment failed');
  }

  const adminHeader = process.env.ADMIN_TOKEN ? { 'X-Admin-Token': process.env.ADMIN_TOKEN } : { 'X-Admin-Token': 'admin-secret' };
  const updateResp = await request(`/shipments/${shipmentId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...adminHeader },
    body: JSON.stringify({ status: 'in-transit', timeline: 'picked up' }),
  });
  console.log(`PUT /shipments/${shipmentId} =>`, updateResp.status, updateResp.data);
  if (updateResp.status !== 200 || updateResp.data?.status !== 'in-transit') {
    throw new Error('Shipment update failed');
  }

  const listResp = await request('/shipments');
  console.log('GET /shipments =>', listResp.status, Array.isArray(listResp.data) ? listResp.data.length : listResp.data);
  if (listResp.status !== 200) {
    throw new Error('List shipments failed');
  }

  const deleteResp = await request(`/shipments/${shipmentId}`, {
    method: 'DELETE',
  });
  console.log(`DELETE /shipments/${shipmentId} =>`, deleteResp.status);
  if (deleteResp.status !== 204) {
    throw new Error('Delete shipment failed');
  }

  const verifyDelete = await request(`/shipments/${shipmentId}`);
  console.log(`GET /shipments/${shipmentId} after delete =>`, verifyDelete.status, verifyDelete.data);
  if (verifyDelete.status !== 404) {
    throw new Error('Expected not found after delete');
  }

  console.log('All shipment endpoint checks passed.');
}

main().catch((error) => {
  console.error('Test failed:', error.message || error);
  process.exit(1);
});
