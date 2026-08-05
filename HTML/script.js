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
        const data = await res.json();
        const el = document.getElementById('guestCount');
        if (el) el.textContent = data;
    } catch (e) {
        console.error('Ошибка загрузки количества:', e);
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
                        msg.textContent = data.message || '✅ Обновлено';
                        msg.className = 'success';
                        // Возвращаем отображение имени
                        nameSpan.textContent = `${index + 1}. ${newFirst} ${newLast}`;
                        nameSpan.style.display = 'inline';
                        editContainer.style.display = 'none';
                        editBtn.style.display = 'inline-block';
                        // Обновляем данные в объекте (для повторного редактирования)
                        g.firstname = newFirst;
                        g.lastname = newLast;
                        // Можно перезагрузить список, чтобы обновить всё
                        // loadModalGuests(); // раскомментируй, если хочешь полный перезапуск
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
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data)
        });
        if (response.ok) {
            message.className = "success";
            message.innerHTML = "✔ Спасибо! Ваше присутствие подтверждено.";
            form.reset();
            await loadGuestCount();
        } else {
            const text = await response.text();
            message.className = "error";
            message.innerHTML = "Ошибка: " + text;
        }
    } catch (err) {
        message.className = "error";
        message.innerHTML = "Сервер недоступен.";
        console.error(err);
    }
});

// ============================================================
// 6. МОДАЛЬНОЕ ОКНО
// ============================================================
const modal = document.getElementById('guestModal');
const modalClose = document.getElementById('modalClose');
const guestCountBlock = document.getElementById('guestCountBlock');
const overlay = document.getElementById('overlay');

if (guestCountBlock) {
    guestCountBlock.addEventListener('click', async () => {
        await loadModalGuests();
        modal.style.display = 'block';
        if (overlay) overlay.style.display = 'block';
    });
}
if (modalClose) {
    modalClose.addEventListener('click', () => {
        modal.style.display = 'none';
        if (overlay) overlay.style.display = 'none';
    });
}
window.addEventListener('click', (e) => {
    if (e.target === modal) {
        modal.style.display = 'none';
        if (overlay) overlay.style.display = 'none';
    }
});

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
// 8. ИНИЦИАЛИЗАЦИЯ ПРИ ЗАГРУЗКЕ СТРАНИЦЫ
// ============================================================
document.addEventListener('DOMContentLoaded', function () {
    // Загружаем счётчик гостей
    loadGuestCount();

    // Инициализация Swiper (карусель)
    if (typeof Swiper !== 'undefined') {
        const swiper = new Swiper('.photo-swiper', {
            loop: true,
            autoplay: {
                delay: 5000,
                disableOnInteraction: false,
            },
            pagination: {
                el: '.swiper-pagination',
                clickable: true,
            },
            navigation: {
                nextEl: '.swiper-button-next',
                prevEl: '.swiper-button-prev',
            },
            effect: 'fade',
            fadeEffect: { crossFade: true },
            speed: 1200,
            grabCursor: true,
            breakpoints: {
                640: { autoplay: false }
            }
        });
    }
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

// ===== УДАЛЕНИЕ ГОСТЯ =====
async function deleteGuest(id) {
    if (!confirm('Удалить гостя из списка?')) return;

    try {
        const res = await fetch(`/user/${id}`, { method: 'DELETE' });
        if (res.ok) {
            const data = await res.json();
            showGuestListMessage(data.message || '✅ Гость удалён', true);
            loadModalGuests();    // обновляем список
            loadGuestCount();     // обновляем счётчик
        } else {
            const text = await res.text();
            showGuestListMessage('❌ Ошибка: ' + text, false);
        }
    } catch (err) {
        showGuestListMessage('❌ Сервер недоступен', false);
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