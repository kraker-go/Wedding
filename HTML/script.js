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
// 7. ПОДСКАЗКА ПРИ ФОКУСЕ НА ПОЛЕ
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

// ============================================================
// 8. ВСЕ АНИМАЦИИ И ИНИЦИАЛИЗАЦИЯ ПРИ ЗАГРУЗКЕ СТРАНИЦЫ
// ============================================================
document.addEventListener('DOMContentLoaded', function () {

    // ---- 8.1. Загрузка счётчика гостей ----
    loadGuestCount();

    // ---- 8.2. Анимация блока «Ждём вас!» ----
    const waitingBlock = document.getElementById('waitingBlock');
    if (waitingBlock) {
        const observerWaiting = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    waitingBlock.classList.add('visible');
                    observerWaiting.unobserve(waitingBlock);
                }
            });
        }, { threshold: 0.3 });
        observerWaiting.observe(waitingBlock);
    }

    // ---- 8.3. Анимация появления таймера ----
    const timerCards = document.querySelectorAll('.timer-card');
    const timerSection = document.querySelector('.timer');
    if (timerSection && timerCards.length) {
        const observerTimer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    timerCards.forEach((card, index) => {
                        setTimeout(() => {
                            card.classList.add('visible');
                        }, index * 200);
                    });
                    observerTimer.unobserve(timerSection);
                }
            });
        }, { threshold: 0.3 });
        observerTimer.observe(timerSection);
    }

    // ---- 8.4. Анимация появления имён и амперсанда ----
    const hero = document.querySelector('.hero');
    const nameLeft = document.querySelector('.name-left');
    const nameRight = document.querySelector('.name-right');
    const ampersand = document.querySelector('.ampersand-glow');
    if (hero && nameLeft && nameRight && ampersand) {
        const observerHero = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    setTimeout(() => nameLeft.classList.add('visible'), 200);
                    setTimeout(() => nameRight.classList.add('visible'), 1600);
                    setTimeout(() => ampersand.classList.add('visible'), 2600);
                    observerHero.unobserve(hero);
                }
            });
        }, { threshold: 0.3 });
        observerHero.observe(hero);
    }

    // ---- 8.5. Анимация кругов с фото и подписей ----
    const coupleSection = document.getElementById('couplePhotos');
    const labels = document.getElementById('coupleLabels');
    if (coupleSection) {
        const observerCouple = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    coupleSection.classList.add('visible');
                    if (labels) {
                        setTimeout(() => {
                            labels.classList.add('visible');
                        }, 2200);
                    }
                    observerCouple.unobserve(coupleSection);
                }
            });
        }, { threshold: 0.3 });
        observerCouple.observe(coupleSection);
    }

    // ---- 8.6. Инициализация Swiper ----
    if (typeof Swiper !== 'undefined') {
        new Swiper('.photo-swiper', {
            loop: true,
            autoplay: { delay: 5000, disableOnInteraction: false },
            pagination: { el: '.swiper-pagination', clickable: true },
            navigation: { nextEl: '.swiper-button-next', prevEl: '.swiper-button-prev' },
            effect: 'fade',
            fadeEffect: { crossFade: true },
            speed: 1200,
            grabCursor: true,
            breakpoints: { 640: { autoplay: false } }
        });
    }

    // ---- 8.7. Показ сердец (hero-bg) ----
    setTimeout(() => {
        if (hero) hero.classList.add('show-hearts');
    }, 2200);

});