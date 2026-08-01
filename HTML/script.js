const form = document.getElementById("guestForm");
const firstName = document.getElementById("firstName");
const lastName = document.getElementById("lastName");
const message = document.getElementById("message");

// ===== ТАЙМЕР =====
const daysEl = document.getElementById("days");
const hoursEl = document.getElementById("hours");
const minutesEl = document.getElementById("minutes");
const secondsEl = document.getElementById("seconds");

const weddingDate = new Date("2026-10-10T00:00:00").getTime();

function updateCountdown() {
    const now = new Date().getTime();
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

// ===== ЗАГРУЗКА КОЛИЧЕСТВА ГОСТЕЙ =====
async function loadGuestCount() {
    try {
        const res = await fetch('/user/count');
        const data = await res.json();
        document.getElementById('guestCount').textContent = data;
    } catch (e) {
        console.error('Ошибка загрузки количества:', e);
    }
}

// ===== ЗАГРУЗКА ПОЛНОГО СПИСКА ДЛЯ МОДАЛКИ =====
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
            div.innerHTML = `
                <span>${index + 1}. ${g.firstname} ${g.lastname}</span>
                <span class="date">${new Date(g.created_at).toLocaleDateString('ru-RU')}</span>
            `;
            list.appendChild(div);
        });
    } catch (e) {
        console.error('Ошибка загрузки списка:', e);
    }
}

// ===== МОДАЛЬНОЕ ОКНО =====
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
    if (e.target === modal) {
        modal.style.display = 'none';
    }
});

// ===== ОТПРАВКА ФОРМЫ =====
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
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(data)
        });

        if (response.ok) {
            message.className = "success";
            message.innerHTML = "✔ Спасибо! Ваше присутствие подтверждено.";
            form.reset();
            await loadGuestCount();  // ← обновляем счётчик
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

// ===== ЗАГРУЗКА ПРИ СТАРТЕ =====
document.addEventListener('DOMContentLoaded', () => {
    loadGuestCount();
});