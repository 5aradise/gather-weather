const apiBase = '/api';

const statusEl = document.getElementById('status');
const subsOutput = document.getElementById('subscriptions-output');
const weatherOutput = document.getElementById('weather-output');

const updateStatus = (message, isError = false) => {
  statusEl.textContent = message;
  statusEl.className = isError ? 'status error' : 'status';
};

const requestForm = async (path, form) => {
  const values = new URLSearchParams(new FormData(form));
  const res = await fetch(`${apiBase}${path}`, {
    method: 'POST',
    body: values,
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `${res.status} ${res.statusText}`);
  }
  return res;
};

const handleCreate = async () => {
  const form = document.getElementById('subscribe-form');
  try {
    await requestForm('/subscriptions', form);
    updateStatus('Subscription created immediately.');
  } catch (err) {
    updateStatus(`Subscribe failed: ${err.message}`, true);
  }
};

const handleList = async () => {
  try {
    const res = await fetch(`${apiBase}/subscriptions`);
    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || `${res.status} ${res.statusText}`);
    }
    const data = await res.json();
    subsOutput.textContent = JSON.stringify(data, null, 2);
    updateStatus('Subscriptions loaded.');
  } catch (err) {
    updateStatus(`Load subscriptions failed: ${err.message}`, true);
  }
};

const handleWeather = async () => {
  const city = document.getElementById('weather-city').value.trim();
  if (!city) return updateStatus('Please enter a city.', true);
  try {
    const res = await fetch(
      `${apiBase}/weather?city=${encodeURIComponent(city)}`,
    );
    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || `${res.status} ${res.statusText}`);
    }
    const data = await res.json();
    weatherOutput.textContent = JSON.stringify(data, null, 2);
    updateStatus('Weather data loaded.');
  } catch (err) {
    updateStatus(`Weather fetch failed: ${err.message}`, true);
  }
};

document.getElementById('create-btn').addEventListener('click', handleCreate);
document.getElementById('list-btn').addEventListener('click', handleList);
document.getElementById('weather-btn').addEventListener('click', handleWeather);
