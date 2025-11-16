function openClaimModal(productId, productName, remaining, qty) {
    document.getElementById('modal-product-name').textContent = productName;
    const unit = productName === 'Eggs' ? ' (dozen)' : '';
    document.getElementById('modal-remaining').textContent = remaining;
    document.getElementById('modal-qty').textContent = qty;
    document.getElementById('modal-available-unit').textContent = unit;
    document.getElementById('modal-qty-unit').textContent = unit;

    // Set hidden form fields
    document.getElementById('product-id').value = productId;

    // Set max for quantity input
    const qtyInput = document.getElementById('claim-qty');
    qtyInput.max = remaining;
    qtyInput.value = '';

    // Clear name input
    document.getElementById('claim-name').value = '';

    // Clear any previous errors
    const errorDiv = document.getElementById('claim-error');
    errorDiv.style.display = 'none';
    errorDiv.textContent = '';

    document.getElementById('claim-modal').style.display = 'flex';
}

function closeClaimModal() {
    document.getElementById('claim-modal').style.display = 'none';
    document.getElementById('claim-form').reset();

    // Clear errors
    const errorDiv = document.getElementById('claim-error');
    errorDiv.style.display = 'none';
    errorDiv.textContent = '';
}

function handleClaimResponse(event) {
    const xhr = event.detail.xhr;
    const errorDiv = document.getElementById('claim-error');

    if (xhr.status === 200) {
        // Success - close modal and refresh page
        closeClaimModal();
        location.reload();
    } else {
        // Error - show error message
        errorDiv.textContent = xhr.responseText;
        errorDiv.style.display = 'block';

        // Scroll error into view
        errorDiv.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
}
