document.addEventListener('DOMContentLoaded', () => {
    const deviceList = document.getElementById('deviceList');
    const addDeviceForm = document.getElementById('addDeviceForm');

    // Fetch and display devices
    fetchDevices();

    // Handle form submission
    addDeviceForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const name = document.getElementById('name').value;
        const mac = document.getElementById('mac').value;

        try {
            const response = await fetch('/devices', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name, mac })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to add device');
            }

            // Reset form
            addDeviceForm.reset();

            // Refresh list
            fetchDevices();
            showToast('Device added successfully!');

        } catch (err) {
            showToast(`Error: ${err.message}`, true);
        }
    });

    async function fetchDevices() {
        try {
            const response = await fetch('/devices');
            if (!response.ok) throw new Error('Failed to fetch devices');

            const devices = await response.json();
            renderDevices(devices);
        } catch (err) {
            deviceList.innerHTML = `<div class="loading">Error loading devices: ${err.message}</div>`;
        }
    }

    function renderDevices(devices) {
        deviceList.innerHTML = '';

        if (!devices || devices.length === 0) {
            deviceList.innerHTML = `
                <div class="empty-state">
                    <p>No devices found. Add one below to get started.</p>
                </div>
            `;
            return;
        }

        devices.forEach(device => {
            const card = document.createElement('div');
            card.className = 'device-card';

            card.innerHTML = `
                <div class="device-info">
                    <div class="device-name">${escapeHtml(device.name)}</div>
                    <div class="device-mac">${escapeHtml(device.mac)}</div>
                </div>
                <div class="device-actions">
                    <button class="btn btn-wake" data-id="${device.id}" data-name="${escapeHtml(device.name)}">
                        ⚡ Wake Up
                    </button>
                    <button class="btn btn-delete" data-id="${device.id}">
                        ✕ Delete
                    </button>
                </div>
            `;

            deviceList.appendChild(card);
        });

        // Add event listeners for buttons
        document.querySelectorAll('.btn-wake').forEach(btn => {
            btn.addEventListener('click', () => wakeDevice(btn.dataset.id));
        });

        document.querySelectorAll('.btn-delete').forEach(btn => {
            btn.addEventListener('click', () => deleteDevice(btn.dataset.id));
        });
    }

    async function wakeDevice(id) {
        try {
            const response = await fetch(`/wake/${id}`, { method: 'POST' });
            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to wake device');
            }
            showToast('Packet sent! ⚡');
        } catch (err) {
            showToast(`Error: ${err.message}`, true);
        }
    }

    async function deleteDevice(id) {
        if (!confirm('Are you sure you want to delete this device?')) return;

        try {
            const response = await fetch(`/devices/${id}`, { method: 'DELETE' });
            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to delete device');
            }
            fetchDevices();
            showToast('Device deleted');
        } catch (err) {
            showToast(`Error: ${err.message}`, true);
        }
    }

    function showToast(message, isError = false) {
        const toast = document.getElementById('toast');
        toast.textContent = message;
        toast.classList.remove('hidden');

        if (isError) {
            toast.style.border = '1px solid #ef4444';
        } else {
            toast.style.border = '1px solid #10b981';
        }

        setTimeout(() => {
            toast.classList.add('hidden');
        }, 3000);
    }

    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
});
