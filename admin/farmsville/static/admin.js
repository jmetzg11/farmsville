document.addEventListener('DOMContentLoaded', function() {
    if (typeof django === 'undefined' || typeof django.jQuery === 'undefined') {
        console.error('django.jQuery is not available');
        return;
    }

    const $ = django.jQuery;

    function updateFieldVisibility(inline) {
        const blockTypeSelect = inline.find('select[name$="-block_type"]');
        const textField = inline.find('textarea[name$="-text_content"]').closest('.form-row, div[class*="field"]');
        const photoField = inline.find('select[name$="-photo"]').closest('.form-row, div[class*="field"]');
        const youtubeField = inline.find('input[name$="-youtube_url"]').closest('.form-row, div[class*="field"]');

        const blockType = blockTypeSelect.val();

        if (blockType === 'text') {
            textField.show();
            photoField.hide();
            youtubeField.hide();
        } else if (blockType === 'photo') {
            textField.hide();
            photoField.show();
            youtubeField.hide();
        } else if (blockType === 'youtube') {
            textField.hide();
            photoField.hide();
            youtubeField.show();
        } else {
            textField.show();
            photoField.show();
            youtubeField.show();
        }
    }

    function setupInline(inline) {
        if (inline && inline.length) {
            updateFieldVisibility(inline);
        }
    }

    // Setup existing inlines
    $('.inline-related').each(function() {
        setupInline($(this));
    });

    // Setup new inlines when they're added
    $(document).on('formset:added', function(event, $row, formsetName) {
        if ($row) {
            setupInline($row);
        } else {
            // Try to find the newly added inline
            const newInline = $('.inline-related').last();
            setupInline(newInline);
        }
    });

    // Use event delegation for change events on block_type selects
    $(document).on('change', 'select[name$="-block_type"]', function() {
        const inline = $(this).closest('.inline-related');
        updateFieldVisibility(inline);
    });
});
