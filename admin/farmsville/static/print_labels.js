document.addEventListener('DOMContentLoaded', function() {
    const labels = document.querySelectorAll('.label');
    const selectAllBtn = document.getElementById('select-all');
    const selectNoneBtn = document.getElementById('select-none');
    const printSelectedBtn = document.getElementById('print-selected');

    // Toggle selection on label click
    labels.forEach(label => {
        label.addEventListener('click', function() {
            this.classList.toggle('selected');
            updatePrintButton();
        });
    });

    // Select all
    selectAllBtn.addEventListener('click', function() {
        labels.forEach(label => label.classList.add('selected'));
        updatePrintButton();
    });

    // Select none
    selectNoneBtn.addEventListener('click', function() {
        labels.forEach(label => label.classList.remove('selected'));
        updatePrintButton();
    });

    // Print selected - sends to local Python bridge
    printSelectedBtn.addEventListener('click', async function() {
        const selectedLabels = document.querySelectorAll('.label.selected');

        if (selectedLabels.length === 0) {
            alert('Please select at least one label to print');
            return;
        }

        printSelectedBtn.disabled = true;
        const originalText = printSelectedBtn.textContent;

        for (const label of selectedLabels) {
            const product = label.dataset.product || '';
            const user = label.dataset.user || '';

            printSelectedBtn.textContent = `Printing: ${product}...`;

            try {
                // Use 127.0.0.1 instead of localhost for maximum compatibility
                const response = await fetch('http://127.0.0.1:5555/print', {
                    method: 'POST',
                    mode: 'cors',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        product: product,
                        user: user
                    })
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    alert(`Error printing ${product}: ${errorData.message}`);
                }
            } catch (err) {
                console.error("Connection Error:", err);
                alert("The Printer Bridge is not responding. Check the PowerShell window!");
                break;
            }
        }

        printSelectedBtn.disabled = false;
        printSelectedBtn.textContent = originalText;
    });

    function updatePrintButton() {
        const selectedCount = document.querySelectorAll('.label.selected').length;
        printSelectedBtn.textContent = selectedCount > 0
            ? `Print Selected (${selectedCount})`
            : 'Print Selected';
    }

    // Start with all selected
    labels.forEach(label => label.classList.add('selected'));
    updatePrintButton();
});
