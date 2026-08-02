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
// 4. ЗАГРУЗКА СПИСКА ГОСТЕЙ ДЛЯ МОДАЛКИ
// ============================================================
async function loadModalGuests() {
    try {
        const res = await fetch('/user/get');
        const guests = await res.json();
        const list = document.getElementById('modalGuestList');
        list.innerHTML = '';
        if (guests.length === 0) {
            list.innerHTML = '<p style="text-align:center;color:#999;">Пока никто не подтвердил 😔</p>';
            return;
        }
        guests.forEach((g, index) => {
            const div = document.createElement('div');
            div.className = 'guest-item';
            div.innerHTML = `<div class="guest-name">${index + 1}. ${g.firstname} ${g.lastname}</div>`;
            list.appendChild(div);
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

if (guestCountBlock) {
    guestCountBlock.addEventListener('click', async () => {
        await loadModalGuests();
        modal.style.display = 'flex';
    });
}
if (modalClose) {
    modalClose.addEventListener('click', () => {
        modal.style.display = 'none';
    });
}
window.addEventListener('click', (e) => {
    if (e.target === modal) modal.style.display = 'none';
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
    }, 4000);
});