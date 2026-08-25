// ============================================================
// 1. ПОЛУЧЕНИЕ ЭЛЕМЕНТОВ
// ============================================================
const form = document.getElementById("guestForm");
const firstName = document.getElementById("firstName");
const lastName = document.getElementById("lastName");
const message = document.getElementById("message");

const daysEl = document.getElementById("days");
const hoursEl = document.getElementById("hours");
const minutesEl = document.getElementById("minutes");
const secondsEl = document.getElementById("seconds");

const weddingDate = new Date("2026-10-10T00:00:00").getTime();

// ============================================================
// 2. ТАЙМЕР (БЕЗ АНИМАЦИИ)
// ============================================================
function updateCountdown() {
    const now = Date.now();
    const distance = weddingDate - now;

    if (distance <= 0) {
        daysEl.textContent = "0";
        hoursEl.textContent = "0";
        minutesEl.textContent = "0";
        secondsEl.textContent = "0";
        return;
    }

    const d = Math.floor(distance / (1000 * 60 * 60 * 24));
    const h = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const m = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    const s = Math.floor((distance % (1000 * 60)) / 1000);

    daysEl.textContent = d;
    hoursEl.textContent = h;
    minutesEl.textContent = m;
    secondsEl.textContent = s;
}

// Запускаем сразу и обновляем каждую секунду
updateCountdown();
setInterval(updateCountdown, 1000);

// ============================================================
// 3. ЗАГРУЗКА КОЛИЧЕСТВА ГОСТЕЙ
// ============================================================
async function loadGuestCount() {
    try {

        const res = await fetch('/user/count');

        if (!res.ok) {
            throw new Error('Ошибка HTTP: ' + res.status);
        }

        const data = await res.text();

        const el = document.getElementById('guestCount');

        if (!el) {
            console.error('Элемент #guestCount не найден');
            return;
        }

        el.textContent = data;

        console.log('Количество гостей:', data);

    } catch (e) {

        console.error(
            'Ошибка загрузки количества:',
            e
        );

    }
}

// ============================================================
// ===== ЗАГРУЗКА СПИСКА ГОСТЕЙ С INLINE-РЕДАКТИРОВАНИЕМ =====
async function loadModalGuests() {
    try {
        const res = await fetch('/user/get');
        const guests = await res.json();
        const list = document.getElementById('modalGuestList');
        const msg = document.getElementById('guestListMessage');
        msg.textContent = '';
        msg.className = '';
        list.innerHTML = '';

        if (guests.length === 0) {
            list.innerHTML = '<p style="text-align:center;color:#999;">Пока никто не подтвердил 😔</p>';
            return;
        }

        guests.forEach((g, index) => {
            const div = document.createElement('div');
            div.className = 'guest-item';
            div.dataset.id = g.id;

            // Контейнер для отображения имени
            const nameSpan = document.createElement('span');
            nameSpan.className = 'guest-name';
            nameSpan.textContent = `${index + 1}. ${g.firstname} ${g.lastname}`;

            // Контейнер для редактирования (скрыт по умолчанию)
            const editContainer = document.createElement('span');
            editContainer.className = 'edit-container';
            editContainer.style.display = 'none';

            const inputFirst = document.createElement('input');
            inputFirst.type = 'text';
            inputFirst.className = 'edit-input';
            inputFirst.value = g.firstname;

            const inputLast = document.createElement('input');
            inputLast.type = 'text';
            inputLast.className = 'edit-input';
            inputLast.value = g.lastname;

            const saveBtn = document.createElement('button');
            saveBtn.textContent = '💾';
            saveBtn.className = 'action-icon edit';
            saveBtn.title = 'Сохранить';

            const cancelBtn = document.createElement('button');
            cancelBtn.textContent = '✖';
            cancelBtn.className = 'action-icon delete';
            cancelBtn.title = 'Отмена';

            editContainer.appendChild(inputFirst);
            editContainer.appendChild(inputLast);
            editContainer.appendChild(saveBtn);
            editContainer.appendChild(cancelBtn);

            // Кнопка редактирования (карандаш)
            const editBtn = document.createElement('button');
            editBtn.className = 'action-icon edit';
            editBtn.textContent = '✏️';
            editBtn.title = 'Редактировать';
            editBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                // Скрываем имя, показываем поля
                nameSpan.style.display = 'none';
                editContainer.style.display = 'inline-block';
                editBtn.style.display = 'none';
            });

            // Кнопка удаления
            const deleteBtn = document.createElement('button');
            deleteBtn.className = 'action-icon delete';
            deleteBtn.textContent = '🗑️';
            deleteBtn.title = 'Удалить';
            deleteBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                deleteGuest(g.id);
            });

            const actions = document.createElement('span');
            actions.appendChild(editBtn);
            actions.appendChild(deleteBtn);

            div.appendChild(nameSpan);
            div.appendChild(editContainer);
            div.appendChild(actions);
            list.appendChild(div);

            // Обработчик сохранения
            saveBtn.addEventListener('click', async () => {
                const newFirst = inputFirst.value.trim();
                const newLast = inputLast.value.trim();
                if (!newFirst || !newLast) {
                    msg.textContent = 'Заполните все поля.';
                    msg.className = 'error';
                    return;
                }
                try {
                    const res = await fetch(`/user/update/${g.id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ id: g.id, firstname: newFirst, lastname: newLast })
                    });
                    if (res.ok) {
                        const data = await res.json();
                        msg.textContent = data.message || '✏️ Запрос на обновление данных ... ';
                        msg.className = 'success';
                        nameSpan.textContent = `${index + 1}. ${newFirst} ${newLast}`;
                        nameSpan.style.display = 'inline';
                        editContainer.style.display = 'none';
                        editBtn.style.display = 'inline-block';
                        g.firstname = newFirst;
                        g.lastname = newLast;
                    } else {
                        const text = await res.text();
                        msg.textContent = '❌ Ошибка: ' + text;
                        msg.className = 'error';
                    }
                } catch (err) {
                    msg.textContent = '❌ Сервер недоступен.';
                    msg.className = 'error';
                }
            });

            // Обработчик отмены
            cancelBtn.addEventListener('click', () => {
                nameSpan.style.display = 'inline';
                editContainer.style.display = 'none';
                editBtn.style.display = 'inline-block';
            });
        });
    } catch (e) {
        console.error('Ошибка загрузки списка:', e);
    }
}
// ============================================================
// 5. ОТПРАВКА ФОРМЫ
// ============================================================
if (form) {
    form.addEventListener("submit", async function (e) {
        e.preventDefault();
        message.innerHTML = "";
        message.className = "";

        const data = {
            firstname: firstName.value.trim(),
            lastname: lastName.value.trim()
        };

        if (!data.firstname || !data.lastname) {
            message.className = "error";
            message.innerHTML = "Заполните все поля.";
            return;
        }

        try {
            const response = await fetch("/user", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                message.className = "success";
                message.innerHTML =
                    "✔ Спасибо! Ваше присутствие подтверждено.";

                form.reset();

                await loadGuestCount();

            } else {

                const text = await response.text();

                message.className = "error";
                message.innerHTML =
                    "Ошибка: " + text;
            }

        } catch (err) {

            message.className = "error";
            message.innerHTML =
                "Сервер недоступен.";

            console.error(err);
        }
    });
}

// ============================================================
// 6. МОДАЛЬНОЕ ОКНО СПИСКА ГОСТЕЙ
// ============================================================

const modal = document.getElementById("guestModal");
const modalClose = document.getElementById("modalClose");
const guestCountBlock = document.getElementById("guestCountBlock");
const overlay = document.getElementById("overlay");

if (guestCountBlock && modal) {

    guestCountBlock.addEventListener("click", async function () {

        await loadModalGuests();

        modal.style.display = "block";

        if (overlay) {
            overlay.style.display = "block";
        }
    });
}

if (modalClose && modal) {

    modalClose.addEventListener("click", function () {

        modal.style.display = "none";

        if (overlay) {
            overlay.style.display = "none";
        }
    });
}

if (overlay && modal) {

    overlay.addEventListener("click", function () {

        modal.style.display = "none";
        overlay.style.display = "none";
    });
}
// ============================================================
// 7. АНИМАЦИЯ БЛОКА «ЖДЁМ ВАС!» ПРИ ПРОКРУТКЕ (оставляем)
// ============================================================
document.addEventListener('DOMContentLoaded', function () {
    const waitingBlock = document.getElementById('waitingBlock');
    if (waitingBlock) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    waitingBlock.classList.add('visible');
                    observer.unobserve(waitingBlock);
                }
            });
        }, { threshold: 0.3 });
        observer.observe(waitingBlock);
    }
});

// ============================================================
// 8. ЗАГРУЗКА ФОТОГРАФИЙ ИЗ /uploads/
// ============================================================

document.addEventListener('DOMContentLoaded', async function () {

    const wrapper = document.getElementById('mainPhotoWrapper');

    if (!wrapper) {
        return;
    }

    try {

        const response = await fetch('/upload');

        if (!response.ok) {
            throw new Error('Ошибка получения фотографий');
        }

        const photos = await response.json();

        wrapper.innerHTML = '';

        photos.forEach(function (photo) {

            const slide = document.createElement('div');
            slide.className = 'swiper-slide';

            const img = document.createElement('img');

            img.src = '/uploads/' + encodeURIComponent(photo);
            img.alt = 'Фото';

            slide.appendChild(img);
            wrapper.appendChild(slide);

        });

        // ====================================================
        // SWIPER
        // ====================================================

        if (typeof Swiper !== 'undefined') {

            new Swiper('.photo-swiper', {

                loop: photos.length > 4,

                autoplay: {
                    delay: 5000,
                    disableOnInteraction: false
                },

                slidesPerView: 4,
                spaceBetween: 10,

                grabCursor: true,

                pagination: {
                    el: '.swiper-pagination',
                    clickable: true
                },

                breakpoints: {

                    0: {
                        slidesPerView: 4,
                        spaceBetween: 6
                    },

                    500: {
                        slidesPerView: 4,
                        spaceBetween: 8
                    },

                    900: {
                        slidesPerView: 4,
                        spaceBetween: 10
                    }

                }

            });

        }

    } catch (error) {

        console.error(
            'Ошибка загрузки фотографий:',
            error
        );

    }

});

// ============================================================
// КЛИК ПО ФОТО → ГАЛЕРЕЯ
// ============================================================

document
    .querySelector('.photo-swiper')
    ?.addEventListener('click', function (event) {

        const slide =
            event.target.closest('.swiper-slide');

        if (!slide) {
            return;
        }

        const slides =
            Array.from(
                document.querySelectorAll(
                    '.photo-swiper .swiper-slide'
                )
            );

        const index =
            slides.indexOf(slide);

        if (index === -1) {
            return;
        }

        window.location.href =
            'gallery.html?photo=' + index;

    });
// ============================================================
// АНИМАЦИЯ ПОЯВЛЕНИЯ ТАЙМЕРА (как у календаря)
// ============================================================
document.addEventListener('DOMContentLoaded', function () {
    const timerCards = document.querySelectorAll('.timer-card');
    const timerSection = document.querySelector('.timer'); // Родительский блок

    if (timerSection && timerCards.length) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    // Добавляем класс 'visible' к каждой карточке с задержкой
                    timerCards.forEach((card, index) => {
                        setTimeout(() => {
                            card.classList.add('visible');
                        }, index * 200); // Задержка 0.2с между карточками
                    });
                    observer.unobserve(timerSection); // Отключаем observer после запуска
                }
            });
        }, { threshold: 0.3 });

        observer.observe(timerSection);
    }
});

// ============================================================
// АНИМАЦИЯ ПОЯВЛЕНИЯ ИМЁН И АМПЕРСАНДА
// ============================================================
document.addEventListener('DOMContentLoaded', function () {
    const hero = document.querySelector('.hero');
    const nameLeft = document.querySelector('.name-left');
    const nameRight = document.querySelector('.name-right');
    const ampersand = document.querySelector('.ampersand-glow');

    if (hero && nameLeft && nameRight && ampersand) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    // 1. Александр прилетает первым
                    setTimeout(() => {
                        nameLeft.classList.add('visible');
                    }, 200);

                    // 2. Через 2 секунды прилетает Милана
                    setTimeout(() => {
                        nameRight.classList.add('visible');
                    }, 1600);

                    // 3. Амперсанд проявляется через 2.5 секунды (между ними)
                    setTimeout(() => {
                        ampersand.classList.add('visible');
                    }, 2600);

                    observer.unobserve(hero);
                }
            });
        }, { threshold: 0.3 });

        observer.observe(hero);
    }
});

window.addEventListener("load", () => {
    setTimeout(() => {
        document.querySelector(".hero").classList.add("show-hearts");
    }, 2200);
});

document.addEventListener('DOMContentLoaded', function () {
    const coupleSection = document.getElementById('couplePhotos');

    if (coupleSection) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    coupleSection.classList.add('visible');
                    observer.unobserve(coupleSection);
                }
            });
        }, { threshold: 0.3 });

        observer.observe(coupleSection);
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const coupleSection = document.getElementById('couplePhotos');
    const labels = document.getElementById('coupleLabels');

    if (coupleSection && labels) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    // Добавляем класс visible кругам (это уже есть в твоём коде, но оставим для ясности)
                    // Здесь мы просто ждём 2 секунды после появления блока и показываем подписи
                    setTimeout(() => {
                        labels.classList.add('visible');
                    }, 2200); // 2.2 сек (чтобы круги успели встать)
                    observer.unobserve(coupleSection);
                }
            });
        }, { threshold: 0.3 });

        observer.observe(coupleSection);
    }
});

// ============================================================
// ПОДСКАЗКА ПРИ ФОКУСЕ НА ПОЛЕ ИМЕНИ
// ============================================================
const hint = document.getElementById('familyHint');

function showHint() {
    if (hint) hint.classList.add('visible');
}

function hideHint() {
    if (hint && document.activeElement !== firstName) {
        hint.classList.remove('visible');
    }
}

if (firstName && hint) {
    firstName.addEventListener('focus', showHint);
    firstName.addEventListener('blur', hideHint);
}

// ===== ПОКАЗ СООБЩЕНИЙ В СПИСКЕ ГОСТЕЙ =====
function showGuestListMessage(text, isSuccess = true) {
    const msg = document.getElementById('guestListMessage');
    msg.textContent = text;
    msg.className = isSuccess ? 'success' : 'error';
    // Автоматически скрыть через 3 секунды
    clearTimeout(msg._hideTimer);
    msg._hideTimer = setTimeout(() => {
        msg.textContent = '';
        msg.className = '';
    }, 3000);
}

// ===== ЗАПРОС НА УДАЛЕНИЕ ГОСТЯ =====
async function deleteGuest(id) {
    if (!confirm('Отправить запрос на удаление гостя?')) {
        return;
    }

    try {
        const res = await fetch(`/user/${id}`, {
            method: 'DELETE'
        });

        if (res.ok) {
            const data = await res.json();

            showGuestListMessage(
                data.message || '🗑️ Запрос на удаление отправлен ...',
                true
            );

        } else {
            const text = await res.text();

            showGuestListMessage(
                '❌ Ошибка: ' + text,
                false
            );
        }

    } catch (err) {

        showGuestListMessage(
            '❌ Сервер недоступен.',
            false
        );

        console.error(err);
    }
}

// ===== РЕДАКТИРОВАНИЕ ГОСТЯ =====
function openEditModal(id, firstname, lastname) {
    console.log('openEditModal called with id:', id);
    document.getElementById('editId').value = id;
    document.getElementById('editFirstName').value = firstname;
    document.getElementById('editLastName').value = lastname;
    document.getElementById('editModal').style.display = 'flex';
    document.getElementById('editMessage').innerHTML = '';
}

// Закрыть редактирование
function closeEditModal() {
    document.getElementById('editModal').style.display = 'none';
}

// Обработчик формы редактирования (должен быть добавлен один раз)
document.getElementById('editForm')?.addEventListener('submit', async function(e) {
    e.preventDefault();
    const id = document.getElementById('editId').value;
    const firstname = document.getElementById('editFirstName').value.trim();
    const lastname = document.getElementById('editLastName').value.trim();

    if (!firstname || !lastname) {
        document.getElementById('editMessage').innerHTML = 'Заполните все поля.';
        return;
    }

    try {
        const res = await fetch(`/user/update/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id, firstname, lastname })
        });
        if (res.ok) {
            const data = await res.json();
            document.getElementById('editMessage').innerHTML = data.message || '✅ Данные обновлены!';
            // Скрыть модалку и обновить список через секунду
            setTimeout(() => {
                closeEditModal();
                loadModalGuests();
                loadGuestCount();
            }, 1000);
        } else {
            const text = await res.text();
            document.getElementById('editMessage').innerHTML = '❌ Ошибка: ' + text;
        }
    } catch (err) {
        document.getElementById('editMessage').innerHTML = '❌ Сервер недоступен.';
    }
});

// Закрытие редактирования по клику вне окна
document.getElementById('editModalClose')?.addEventListener('click', closeEditModal);
window.addEventListener('click', function(e) {
    const modal = document.getElementById('editModal');
    if (e.target === modal) closeEditModal();
});

document.addEventListener('DOMContentLoaded', function () {
    
    fetch('/visit').catch(err => {
        console.error('Ошибка регистрации посещения:', err);
    });

    const video = document.getElementById('weddingVideo');
    const musicButton = document.getElementById('musicToggle');
    const venueText = document.getElementById('venueText');

    if (!video) return;


    // ====================================================
    // НАЧАЛЬНОЕ СОСТОЯНИЕ ВИДЕО
    // ====================================================

    video.muted = true;
    video.volume = 0;

    video.setAttribute('muted', '');
    video.setAttribute('playsinline', '');

    video.play().catch(function () {});


    // ====================================================
// ПОЯВЛЕНИЕ ТЕКСТА НА ВИДЕО
// ====================================================

    if (video && venueText) {

        const venueObserver = new IntersectionObserver(
            function (entries) {

                const entry = entries[0];

                if (!entry.isIntersecting) {
                    return;
                }

                // Запускаем всю CSS-анимацию
                venueText.classList.add('animate');

                // Больше не отслеживаем
                venueObserver.unobserve(video);

            },
            {
                threshold: 0.25
            }
        );

        venueObserver.observe(video);
    }


    // ====================================================
    // КЛИК ПО ВИДЕО
    // ====================================================

    video.addEventListener('click', function () {

        // Всегда начинаем видео с начала
        video.currentTime = 0;

        // Включаем звук
        video.muted = false;
        video.volume = 1;

        video.removeAttribute('muted');

        // Запускаем видео
        video.play().catch(function () {});


        // Android / Chrome / Firefox / Edge
        if (video.requestFullscreen) {

            video.requestFullscreen().catch(function () {});

            return;
        }


        // iPhone / iPad Safari
        if (typeof video.webkitEnterFullscreen === 'function') {

            try {
                video.webkitEnterFullscreen();
            } catch (error) {
                console.log('Fullscreen недоступен');
            }
        }

    });


    // ====================================================
    // КНОПКА ЗВУКА
    // ====================================================

    if (musicButton) {

        musicButton.textContent = '🔊';

        musicButton.addEventListener('click', function () {

            if (video.muted) {

                // Включить звук
                video.muted = false;
                video.volume = 1;

                musicButton.textContent = '🔊';
                musicButton.setAttribute(
                    'aria-label',
                    'Выключить звук'
                );

                video.play().catch(function () {});

            } else {

                // Выключить звук
                video.muted = true;
                video.volume = 0;

                musicButton.textContent = '🔇';
                musicButton.setAttribute(
                    'aria-label',
                    'Включить звук'
                );
            }

        });

    }

});

document.addEventListener('DOMContentLoaded', function () {
    loadGuestCount();
});