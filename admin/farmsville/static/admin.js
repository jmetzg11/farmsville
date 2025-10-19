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
        const blockTypeSelect = inline.find('select[name$="-block_type"]');
        if (blockTypeSelect.length) {
            blockTypeSelect.on('change', function() {
                updateFieldVisibility(inline);
            });
            updateFieldVisibility(inline);
        }
    }

    // Setup existing inlines
    $('.inline-related').each(function() {
        setupInline($(this));
    });

    // Setup new inlines when they're added
    $(document).on('formset:added', function(event, $row, formsetName) {
        setupInline($row);
    });
});
