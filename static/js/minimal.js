document.addEventListener('DOMContentLoaded', function () {
    const typeSelect = document.getElementById('type_id');
    const categorySelect = document.getElementById('category_id');
    const categoryOptions = Array.from(categorySelect.options);

    function toggleCategoryOptions() {
        const selectedType = typeSelect.value;
        let hasVisibleOptions = false;

        categoryOptions.forEach(option => {
            if (option.value === '') {
                option.style.display = 'block';
                return;
            }

            const optionType = option.className.match(/category-(\d+)/)[1];

            if (selectedType === optionType) {
                option.style.display = 'block';
                hasVisibleOptions = true;
            } else {
                option.style.display = 'none';
            }
        });

        if (!hasVisibleOptions) {
            categorySelect.value = '';
        }
    }

    typeSelect.addEventListener('change', toggleCategoryOptions);
    toggleCategoryOptions();
});