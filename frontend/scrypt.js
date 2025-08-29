const API_BASE_URL = 'http://localhost:4141/api/v1';

function searchOrder() {
    const orderId = document.getElementById('orderId').value.trim();

    if (!orderId) {
        showError('Пожалуйста, введите ID заказа');
        return;
    }

    showLoading(true);
    hideError();
    hideResult();

    fetch(`${API_BASE_URL}/order/${orderId}`)
        .then(response => {
            if (!response.ok) {
                if (response.status === 404) {
                    throw new Error('Заказ не найден');
                }
                throw new Error(`Ошибка сервера: ${response.status}`);
            }
            return response.json();
        })
        .then(order => {
            displayOrder(order);
            showLoading(false);
        })
        .catch(error => {
            showLoading(false);
            showError(error.message);
            console.error('Error:', error);
        });
}

function displayOrder(order) {
    // Основная информация
    document.getElementById('order-uid').textContent = order.order_uid || 'N/A';
    document.getElementById('track-number').textContent = order.track_number || 'N/A';
    document.getElementById('entry').textContent = order.entry || 'N/A';
    document.getElementById('locale').textContent = order.locale || 'N/A';
    document.getElementById('customer-id').textContent = order.customer_id || 'N/A';
    document.getElementById('delivery-service').textContent = order.delivery_service || 'N/A';
    document.getElementById('date-created').textContent = formatDate(order.date_created) || 'N/A';

    // Доставка
    if (order.delivery) {
        document.getElementById('delivery-name').textContent = order.delivery.name || 'N/A';
        document.getElementById('delivery-phone').textContent = order.delivery.phone || 'N/A';
        document.getElementById('delivery-email').textContent = order.delivery.email || 'N/A';
        document.getElementById('delivery-address').textContent = order.delivery.address || 'N/A';
        document.getElementById('delivery-city').textContent = order.delivery.city || 'N/A';
        document.getElementById('delivery-region').textContent = order.delivery.region || 'N/A';
        document.getElementById('delivery-zip').textContent = order.delivery.zip || 'N/A';
    }

    // Оплата
    if (order.payment) {
        document.getElementById('payment-transaction').textContent = order.payment.transaction || 'N/A';
        document.getElementById('payment-currency').textContent = order.payment.currency || 'N/A';
        document.getElementById('payment-amount').textContent = order.payment.amount ? `${order.payment.amount} USD` : 'N/A';
        document.getElementById('payment-provider').textContent = order.payment.provider || 'N/A';
        document.getElementById('payment-bank').textContent = order.payment.bank || 'N/A';
        document.getElementById('payment-delivery-cost').textContent = order.payment.delivery_cost ? `${order.payment.delivery_cost} USD` : 'N/A';
        document.getElementById('payment-dt').textContent = formatTimestamp(order.payment.payment_dt) || 'N/A';
    }

    // Товары
    const itemsList = document.getElementById('items-list');
    itemsList.innerHTML = '';

    if (order.items && order.items.length > 0) {
        order.items.forEach((item, index) => {
            const itemDiv = document.createElement('div');
            itemDiv.className = 'item';
            itemDiv.innerHTML = `
                <p><strong>Товар ${index + 1}:</strong> ${item.name || 'N/A'}</p>
                <p><strong>Бренд:</strong> ${item.brand || 'N/A'}</p>
                <p><strong>Цена:</strong> ${item.price ? `${item.price} USD` : 'N/A'}</p>
                <p><strong>Цена со скидкой:</strong> ${item.total_price ? `${item.total_price} USD` : 'N/A'}</p>
                <p><strong>Размер:</strong> ${item.size || 'N/A'}</p>
                <p><strong>Статус:</strong> ${item.status || 'N/A'}</p>
                <p><strong>Артикул:</strong> ${item.chrt_id || 'N/A'}</p>
            `;
            itemsList.appendChild(itemDiv);
        });
    } else {
        itemsList.innerHTML = '<p>Товары не найдены</p>';
    }

    showResult();
}

function formatDate(dateString) {
    if (!dateString) return 'N/A';
    try {
        const date = new Date(dateString);
        return date.toLocaleString('ru-RU');
    } catch (e) {
        return dateString;
    }
}

function formatTimestamp(timestamp) {
    if (!timestamp) return 'N/A';
    try {
        const date = new Date(timestamp * 1000);
        return date.toLocaleString('ru-RU');
    } catch (e) {
        return timestamp;
    }
}

function showLoading(show) {
    const loading = document.getElementById('loading');
    loading.classList.toggle('hidden', !show);
}

function showError(message) {
    const errorDiv = document.getElementById('error');
    errorDiv.textContent = message;
    errorDiv.classList.remove('hidden');
}

function hideError() {
    const errorDiv = document.getElementById('error');
    errorDiv.classList.add('hidden');
}

function showResult() {
    document.getElementById('result').classList.remove('hidden');
}

function hideResult() {
    document.getElementById('result').classList.add('hidden');
}

// Обработка нажатия Enter в поле ввода
document.getElementById('orderId').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        searchOrder();
    }
});