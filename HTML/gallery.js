document.addEventListener('DOMContentLoaded', function () {

    const grid = document.getElementById('galleryGrid');
    const input = document.getElementById('photoInput');
    const lightbox = document.getElementById('galleryLightbox');
    const wrapper = document.getElementById('lightboxWrapper');
    const closeButton = document.getElementById('lightboxClose');
    const deleteButton = document.getElementById('lightboxDelete');

    if (!grid || !lightbox || !wrapper) {
        return;
    }

    // ====================================================
    // СОСТОЯНИЕ ВЫБРАННЫХ ФОТО
    // ====================================================

    let selectedFiles = [];


    // ====================================================
    // SWIPER
    // ====================================================

    const lightboxSwiper = new Swiper(
        '.lightbox-swiper',
        {
            loop: false,
            speed: 600,
            slidesPerView: 1,
            spaceBetween: 20,
            centeredSlides: true,
            grabCursor: true,

            navigation: {
                nextEl: '.swiper-button-next',
                prevEl: '.swiper-button-prev'
            }
        }
    );

    if (deleteButton) {

        deleteButton.addEventListener(
            'click',
            async function () {

                const activeSlide =
                    lightboxSwiper.slides[
                        lightboxSwiper.activeIndex
                        ];


                if (!activeSlide) {
                    return;
                }


                const filename =
                    activeSlide.dataset.filename;


                if (!filename) {
                    return;
                }


                const confirmed =
                    confirm(
                        'Вы уверены, что хотите удалить это фото?'
                    );


                if (!confirmed) {
                    return;
                }


                deleteButton.disabled = true;

                deleteButton.textContent =
                    'Удаление...';


                try {

                    const response =
                        await fetch(
                            '/upload/' +
                            encodeURIComponent(filename),
                            {
                                method: 'DELETE'
                            }
                        );


                    if (!response.ok) {

                        const errorText =
                            await response.text();

                        throw new Error(
                            errorText ||
                            'Ошибка удаления фотографии'
                        );

                    }


                    // Обновляем галерею
                    await loadPhotos();


                    const slidesCount =
                        lightboxSwiper.slides.length;


                    if (slidesCount === 0) {

                        lightbox.classList.remove(
                            'active'
                        );

                    } else {

                        let newIndex =
                            lightboxSwiper.activeIndex;

                        if (
                            newIndex >=
                            slidesCount
                        ) {
                            newIndex =
                                slidesCount - 1;
                        }

                        lightboxSwiper.slideTo(
                            newIndex,
                            0
                        );

                    }


                    console.log(
                        'Фото удалено:',
                        filename
                    );


                } catch (err) {

                    console.error(
                        'Ошибка удаления фотографии:',
                        err
                    );


                    alert(
                        'Не удалось удалить фотографию.'
                    );

                } finally {

                    deleteButton.disabled = false;

                    deleteButton.textContent =
                        '🗑 Удалить фото';

                }

            }
        );

    }

    // ====================================================
    // ПРЕДПРОСМОТР
    // ====================================================

    function createPreview() {

        const oldPreview =
            document.getElementById('photoPreviewOverlay');

        if (oldPreview) {
            oldPreview.remove();
        }

        if (!selectedFiles.length) {
            return;
        }


        // ====================================================
        // OVERLAY
        // ====================================================

        const overlay =
            document.createElement('div');

        overlay.id = 'photoPreviewOverlay';


        // ====================================================
        // КОНТЕЙНЕР
        // ====================================================

        const container =
            document.createElement('div');

        container.className =
            'photo-preview-container';


        // ====================================================
        // ЗАГОЛОВОК
        // ====================================================

        const title =
            document.createElement('h3');

        title.textContent =
            'Предпросмотр фотографий';


        // ====================================================
        // СПИСОК ФОТО
        // ====================================================

        const previewGrid =
            document.createElement('div');

        previewGrid.className =
            'photo-preview-grid';


        selectedFiles.forEach(function (file) {

            if (!file.type.startsWith('image/')) {
                return;
            }


            const item =
                document.createElement('div');

            item.className =
                'photo-preview-item';


            const image =
                document.createElement('img');

            const objectUrl =
                URL.createObjectURL(file);

            image.src = objectUrl;

            image.alt = file.name;


            image.onload = function () {
                URL.revokeObjectURL(objectUrl);
            };


            const filename =
                document.createElement('div');

            filename.className =
                'photo-preview-name';

            filename.textContent =
                file.name;


            item.appendChild(image);
            item.appendChild(filename);

            previewGrid.appendChild(item);

        });


        // ====================================================
        // КНОПКИ
        // ====================================================

        const buttons =
            document.createElement('div');

        buttons.className =
            'photo-preview-buttons';


        // ----------------------------------------------------
        // ОТМЕНИТЬ
        // ----------------------------------------------------

        const cancelButton =
            document.createElement('button');

        cancelButton.type = 'button';

        cancelButton.className =
            'photo-preview-cancel';

        cancelButton.textContent =
            'Отменить';


        cancelButton.addEventListener(
            'click',
            function () {

                overlay.remove();

                selectedFiles = [];

                input.value = '';

            }
        );


        // ----------------------------------------------------
        // ПОДТВЕРДИТЬ
        // ----------------------------------------------------

        const confirmButton =
            document.createElement('button');

        confirmButton.type = 'button';

        confirmButton.className =
            'photo-preview-confirm';

        confirmButton.textContent =
            'Загрузить';


        confirmButton.addEventListener(
            'click',
            async function () {

                if (!selectedFiles.length) {
                    return;
                }


                confirmButton.disabled = true;

                cancelButton.disabled = true;

                confirmButton.textContent =
                    'Загрузка...';


                try {

                    for (
                        const file of selectedFiles
                        ) {

                        const formData =
                            new FormData();

                        formData.append(
                            'photo',
                            file
                        );


                        const response =
                            await fetch(
                                '/upload',
                                {
                                    method: 'POST',
                                    body: formData
                                }
                            );


                        if (!response.ok) {

                            const errorText =
                                await response.text();

                            throw new Error(
                                errorText ||
                                'Ошибка загрузки фотографии'
                            );

                        }

                    }


                    // ====================================================
                    // УСПЕШНАЯ ЗАГРУЗКА
                    // ====================================================

                    overlay.remove();

                    selectedFiles = [];

                    input.value = '';

                    await loadPhotos();


                } catch (err) {

                    console.error(
                        'Ошибка загрузки фотографий:',
                        err
                    );


                    alert(
                        'Не удалось загрузить фотографии.'
                    );


                    confirmButton.disabled = false;

                    cancelButton.disabled = false;

                    confirmButton.textContent =
                        'Загрузить';

                }

            }
        );

        // ====================================================
        // СОБИРАЕМ
        // ====================================================

        buttons.appendChild(
            cancelButton
        );

        buttons.appendChild(
            confirmButton
        );


        container.appendChild(title);

        container.appendChild(
            previewGrid
        );

        container.appendChild(
            buttons
        );

        overlay.appendChild(
            container
        );

        document.body.appendChild(
            overlay
        );
    }


    // ====================================================
    // СТИЛИ ПРЕДПРОСМОТРА
    // ====================================================

    function injectPreviewStyles() {

        if (
            document.getElementById(
                'photoPreviewStyles'
            )
        ) {
            return;
        }

        const style =
            document.createElement('style');

        style.id =
            'photoPreviewStyles';

        style.textContent = `
    /* ================================================
       ОКНО ПРЕДПРОСМОТРА
       ================================================ */

    #photoPreviewOverlay {
        position: fixed;
        inset: 0;

        z-index: 100001;

        display: flex;
        align-items: center;
        justify-content: center;

        padding: 20px;

        box-sizing: border-box;

        background: rgba(20, 16, 12, 0.94);

        overflow-y: auto;
    }


    /* ================================================
       КОНТЕЙНЕР
       ================================================ */

    .photo-preview-container {
        width: 94%;
        max-width: 1000px;

        max-height: 90vh;

        box-sizing: border-box;

        padding: 25px;

        background: #ffffff;

        border-radius: 16px;

        overflow-y: auto;
    }


    /* ================================================
       ЗАГОЛОВОК
       ================================================ */

    .photo-preview-container h3 {
        margin: 0 0 20px;

        text-align: center;

        font-family:
            'Montserrat',
            sans-serif;

        font-size: 24px;

        color: #6d614f;
    }


    /* ================================================
       СЕТКА ФОТО
       ================================================ */

    .photo-preview-grid {
        display: flex;

        flex-wrap: wrap;

        justify-content: center;

        align-items: flex-start;

        gap: 15px;

        width: 100%;

        box-sizing: border-box;
    }


    /* ================================================
       ОДНО ФОТО
       ================================================ */

    .photo-preview-item {
        display: flex;

        flex-direction: column;

        align-items: center;

        justify-content: flex-start;

        width: auto;

        max-width: 100%;

        box-sizing: border-box;
    }


    /* ================================================
       КАРТИНКА
       ================================================ */

    .photo-preview-item img {
        display: block;

        width: auto;

        height: auto;

        max-width: 100%;

        max-height: 350px;

        object-fit: contain;

        margin: 0 auto;
    }


    /* ================================================
       ИМЯ ФАЙЛА
       ================================================ */

    .photo-preview-name {
        width: 100%;

        max-width: 100%;

        box-sizing: border-box;

        padding: 8px;

        font-family:
            'Montserrat',
            sans-serif;

        font-size: 12px;

        color: #6d614f;

        text-align: center;

        white-space: nowrap;

        overflow: hidden;

        text-overflow: ellipsis;
    }


    /* ================================================
       КНОПКИ
       ================================================ */

    .photo-preview-buttons {
        display: flex;

        justify-content: center;

        gap: 15px;

        margin-top: 20px;
    }


    .photo-preview-buttons button {
        min-width: 150px;

        padding: 12px 24px;

        border: none;

        border-radius: 40px;

        font-family:
            'Montserrat',
            sans-serif;

        font-size: 16px;

        font-weight: 600;

        cursor: pointer;

        transition:
            transform 0.2s ease,
            opacity 0.2s ease;
    }


    .photo-preview-buttons button:hover {
        transform: translateY(-2px);
    }


    .photo-preview-buttons button:disabled {
        opacity: 0.6;

        cursor: wait;

        transform: none;
    }


    .photo-preview-cancel {
        background: #e8dfd2;

        color: #5c5145;

        box-shadow:
            0 4px 10px
            rgba(0, 0, 0, 0.15);
    }


    .photo-preview-confirm {
        background: linear-gradient(
            145deg,
            #f5ead5 0%,
            #d8c29a 45%,
            #b99a63 100%
        );

        color: #3f352b;

        box-shadow:
            inset 0 2px 3px
            rgba(255, 255, 255, 0.8),

            inset 0 -3px 5px
            rgba(90, 65, 30, 0.25),

            0 3px 0 #9a7b48,

            0 7px 12px
            rgba(70, 50, 25, 0.25);
    }


    /* ================================================
       ТЕЛЕФОН
       ================================================ */

    @media (max-width: 600px) {

        #photoPreviewOverlay {
            padding: 10px;
        }


        .photo-preview-container {
            width: 94%;

            max-width: 94%;

            max-height: 88vh;

            padding: 15px;

            border-radius: 16px;

            box-sizing: border-box;
        }


        .photo-preview-container h3 {
            font-size: 22px;

            margin-bottom: 15px;
        }


        .photo-preview-grid {
            display: flex;

            flex-direction: column;

            align-items: center;

            justify-content: flex-start;

            width: 100%;

            gap: 10px;
        }


        .photo-preview-item {
            width: 100%;

            max-width: 100%;

            display: flex;

            flex-direction: column;

            align-items: center;

            justify-content: center;

            box-sizing: border-box;
        }


        .photo-preview-item img {
            display: block;

            width: auto;

            height: auto;

            max-width: 100%;

            max-height: 250px;

            object-fit: contain;

            margin: 0 auto;
        }


        .photo-preview-name {
            max-width: 100%;

            text-align: center;

            font-size: 11px;
        }


        .photo-preview-buttons {
            flex-direction: column;

            align-items: center;

            gap: 10px;

            margin-top: 15px;
        }


        .photo-preview-buttons button {
            width: 100%;

            min-width: 0;
        }

    }
`;


        document.head.appendChild(style);
    }


    // ====================================================
    // СОЗДАЁМ СЛАЙДЫ LIGHTBOX
    // ====================================================

    function buildSlides() {

        wrapper.innerHTML = '';


        const images =
            grid.querySelectorAll('img');


        images.forEach(function (image) {

            const slide =
                document.createElement('div');

            slide.className =
                'swiper-slide';


            // Запоминаем имя файла
            slide.dataset.filename =
                image.dataset.filename;


            const img =
                document.createElement('img');

            img.src =
                image.src;

            img.alt =
                image.alt || '';


            slide.appendChild(img);

            wrapper.appendChild(slide);

        });


        lightboxSwiper.update();

    }


    // ====================================================
    // ЗАГРУЗКА ФОТО С СЕРВЕРА
    // ====================================================

    async function loadPhotos() {

        try {

            const response =
                await fetch('/upload');


            if (!response.ok) {
                throw new Error(
                    'Ошибка получения фотографий'
                );
            }


            const photos =
                await response.json();


            grid.innerHTML = '';


            const imagePromises = [];


            photos.forEach(function (photo) {

                const photoContainer =
                    document.createElement('div');

                photoContainer.className =
                    'gallery-photo';


                const img =
                    document.createElement('img');

                img.src =
                    '/uploads/' +
                    encodeURIComponent(photo);

                img.alt =
                    'Фото';


                // Запоминаем имя файла
                // для последующего удаления
                img.dataset.filename =
                    photo;


                photoContainer.appendChild(img);

                grid.appendChild(photoContainer);


                imagePromises.push(
                    new Promise(function (resolve) {

                        if (img.complete) {
                            resolve();
                            return;
                        }

                        img.onload =
                            resolve;

                        img.onerror =
                            resolve;

                    })
                );

            });


            await Promise.all(
                imagePromises
            );


            buildSlides();


        } catch (err) {

            console.error(
                'Не удалось загрузить фотографии:',
                err
            );

        }

    }

// ====================================================
// УДАЛЕНИЕ ФОТО
// ====================================================


    // ====================================================
    // ОТКРЫТЬ LIGHTBOX
    // ====================================================

    function openLightbox(index) {

        lightbox.classList.add('active');

        requestAnimationFrame(function () {

            lightboxSwiper.update();

            lightboxSwiper.slideTo(
                index,
                0
            );

        });

    }


    // ====================================================
    // КЛИК ПО ФОТО
    // ====================================================

    grid.addEventListener(
        'click',
        function (event) {

            const image =
                event.target.closest('img');

            if (!image) {
                return;
            }

            const images =
                Array.from(
                    grid.querySelectorAll('img')
                );

            const index =
                images.indexOf(image);

            if (index === -1) {
                return;
            }

            openLightbox(index);

        }
    );


    // ====================================================
    // КРЕСТИК
    // ====================================================

    if (closeButton) {

        closeButton.addEventListener(
            'click',
            function () {

                lightbox.classList.remove(
                    'active'
                );

            }
        );

    }


    // ====================================================
    // КЛИК ПО ФОНУ
    // ====================================================

    lightbox.addEventListener(
        'click',
        function (event) {

            if (event.target === lightbox) {

                lightbox.classList.remove(
                    'active'
                );

            }

        }
    );


    // ====================================================
    // ESC
    // ====================================================

    document.addEventListener(
        'keydown',
        function (event) {

            if (
                event.key === 'Escape' &&
                lightbox.classList.contains('active')
            ) {

                lightbox.classList.remove(
                    'active'
                );

            }

        }
    );


    // ====================================================
    // ВЫБОР ФОТО
    // ====================================================

    if (input) {

        input.addEventListener(
            'change',
            function () {

                selectedFiles =
                    Array.from(input.files)
                        .filter(function (file) {

                            return file.type
                                .startsWith('image/');

                        });


                if (!selectedFiles.length) {
                    return;
                }


                // Создаём стили один раз
                injectPreviewStyles();

                // Показываем preview
                createPreview();



            }
        );

    }


    // ====================================================
    // ОТКРЫТИЕ ФОТО ПО ?photo=2
    // ====================================================

    let startPhotoIndex = null;

    const urlParams =
        new URLSearchParams(
            window.location.search
        );

    const photoParam =
        urlParams.get('photo');


    if (photoParam !== null) {

        const index =
            Number(photoParam);

        if (
            Number.isInteger(index) &&
            index >= 0
        ) {
            startPhotoIndex = index;
        }

        window.history.replaceState(
            {},
            document.title,
            window.location.pathname
        );

    }


    // ====================================================
    // ЗАГРУЖАЕМ ФОТО ПРИ ОТКРЫТИИ
    // ====================================================

    loadPhotos().then(function () {

        if (startPhotoIndex !== null) {

            openLightbox(
                startPhotoIndex
            );

        }

    });

});

