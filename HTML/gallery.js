document.addEventListener('DOMContentLoaded', function () {

    const grid = document.getElementById('galleryGrid');
    const input = document.getElementById('photoInput');
    const lightbox = document.getElementById('galleryLightbox');
    const wrapper = document.getElementById('lightboxWrapper');
    const closeButton = document.getElementById('lightboxClose');

    if (!grid || !lightbox || !wrapper) {
        return;
    }


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


    // ====================================================
    // СОЗДАЁМ СЛАЙДЫ
    // ====================================================

    function buildSlides() {

        wrapper.innerHTML = '';

        const images = grid.querySelectorAll('img');

        images.forEach(function (image) {

            const slide = document.createElement('div');

            slide.className = 'swiper-slide';

            const img = document.createElement('img');

            img.src = image.src;
            img.alt = image.alt || '';

            slide.appendChild(img);

            wrapper.appendChild(slide);

        });

        lightboxSwiper.update();
    }


    // ====================================================
    // ЗАГРУЗКА ФОТО ИЗ СЕРВЕРА
    // ====================================================

    async function loadPhotos() {

        try {

            const response = await fetch('/upload');

            if (!response.ok) {
                throw new Error('Ошибка получения фотографий');
            }

            const photos = await response.json();

            grid.innerHTML = '';

            const imagePromises = [];

            photos.forEach(function (photo) {

                const img = document.createElement('img');

                img.src =
                    '/uploads/' +
                    encodeURIComponent(photo);

                img.alt = 'Фото';

                grid.appendChild(img);

                // Ждём фактической загрузки картинки
                imagePromises.push(
                    new Promise(function (resolve) {

                        if (img.complete) {
                            resolve();
                            return;
                        }

                        img.onload = resolve;
                        img.onerror = resolve;

                    })
                );

            });

            // Ждём загрузки всех фотографий
            await Promise.all(imagePromises);

            // Теперь строим Swiper
            buildSlides();

        } catch (err) {

            console.error(
                'Не удалось загрузить фотографии:',
                err
            );

        }
    }


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
    // ДОБАВЛЕНИЕ ФОТО
    // ====================================================

    if (input) {

        input.addEventListener(
            'change',
            async function () {

                const files =
                    Array.from(input.files);

                if (!files.length) {
                    return;
                }

                for (const file of files) {

                    if (
                        !file.type.startsWith('image/')
                    ) {
                        continue;
                    }

                    const formData =
                        new FormData();

                    formData.append(
                        'photo',
                        file
                    );

                    try {

                        const response =
                            await fetch(
                                '/upload',
                                {
                                    method: 'POST',
                                    body: formData
                                }
                            );

                        if (!response.ok) {
                            throw new Error(
                                'Ошибка загрузки'
                            );
                        }

                        console.log(
                            'Фото загружено:',
                            file.name
                        );

                    } catch (err) {

                        console.error(
                            'Ошибка загрузки:',
                            err
                        );

                    }
                }

                input.value = '';

                // После загрузки заново
                // получаем фотографии с сервера

                await loadPhotos();

            }
        );

    }


    // ====================================================
    // ОТКРЫТИЕ ФОТО ПО ?photo=2
    // ====================================================

    let startPhotoIndex = null;

    const urlParams =
        new URLSearchParams(window.location.search);

    const photoParam =
        urlParams.get('photo');

    if (photoParam !== null) {

        const index = Number(photoParam);

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

            openLightbox(startPhotoIndex);

        }

        });
});