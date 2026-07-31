const form = document.getElementById("guestForm");

const firstName = document.getElementById("firstName");
const lastName = document.getElementById("lastName");

const message = document.getElementById("message");

const days = document.getElementById("days");
const hours = document.getElementById("hours");
const minutes = document.getElementById("minutes");
const seconds = document.getElementById("seconds");


// -------------------------
// Обратный отсчет
// -------------------------

const weddingDate = new Date("2026-10-10T00:00:00").getTime();

function updateCountdown() {

    const now = new Date().getTime();

    const distance = weddingDate - now;

    if (distance <= 0) {

        days.textContent = "0";
        hours.textContent = "0";
        minutes.textContent = "0";
        seconds.textContent = "0";

        return;
    }

    const d = Math.floor(distance / (1000 * 60 * 60 * 24));

    const h = Math.floor(
        (distance % (1000 * 60 * 60 * 24))
        / (1000 * 60 * 60)
    );

    const m = Math.floor(
        (distance % (1000 * 60 * 60))
        / (1000 * 60)
    );

    const s = Math.floor(
        (distance % (1000 * 60))
        / 1000
    );

    days.textContent = d;
    hours.textContent = h;
    minutes.textContent = m;
    seconds.textContent = s;

}

updateCountdown();

setInterval(updateCountdown, 1000);


// -------------------------
// Отправка формы
// -------------------------

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

        const response = await fetch("http://localhost:8080/user", {

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