const getId = (id) => document.getElementById(id);

const orderInfo = getId("orderInfo");
const orderNotFound = getId("orderNotFound");
const orderContainer = getId("orderContainer");
const terminalOutput = getId("terminalOutput");
const inputField = getId("orderInput");

inputField.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
        const orderUID = inputField.value;

        loadOrder(orderUID);
    }
});

function loadOrder(orderUID) {
    if (orderUID != "") {
        fetch(`/orders/${orderUID}`)
            .then((responce) => {
                if (responce.status === 404) {
                    return null;
                }
                return responce.json();
            })
            .then((order) => {
                if (order) {
                    renderOrderInfo(order);
                } else {
                    renderOrderNotFound(orderUID);
                }
            })
            .catch((error) => {
                console.error("Error fetching order:", error);
            });
    } else {
        console.log("Empty order_uid is ignored");
    }
}

function renderOrderNotFound(orderUID) {
    orderInfo.style.display = "none";
    orderContainer.style.display = "flex";
    terminalOutput.style.display = "block";
    orderNotFound.style.display = "block";
    orderNotFound.textContent = "Заказ не найден";

    terminalOutput.innerHTML = "";

    const text = [
        `~$ orderctl get "${orderUID}"`,
        "[INFO] Connecting to database: orders_db",
        "[INFO] Querying order...",
        "[INFO] Validating order data...",
        "[ERROR] Can't retrieve order: not found in database",
        "[INFO] Triggering not-found handler...",
    ];

    text.forEach((line) => {
        const div = document.createElement("div");

        if (line.includes("[INFO]")) {
            div.innerHTML = line.replace(
                "INFO",
                '<span style="color:#004003">INFO</span>',
            );
        } else if (line.includes("[ERROR]")) {
            div.innerHTML = line.replace(
                "ERROR",
                '<span style="color:#590000">ERROR</span>',
            );
        } else {
            div.textContent = line;
        }

        terminalOutput.appendChild(div);

        const emptyDiv = document.createElement("div");
        emptyDiv.innerHTML = "&nbsp;";
        terminalOutput.appendChild(emptyDiv);
    });
}

function renderOrderInfo(order) {
    terminalOutput.style.display = "none";
    orderNotFound.style.display = "none";
    orderContainer.style.display = "block";
    orderInfo.style.display = "block";

    const dateCreated = new Date(order.date_created);
    const dateFormatted = dateCreated.toLocaleString("ru-RU");

    getId("orderUID").textContent = order.order_uid;
    getId("trackNumber").textContent = order.track_number;
    getId("creationDate").textContent = dateFormatted;
    getId("locale").textContent = order.locale;
    getId("deliveryService").textContent = order.delivery_service;
    getId("internalSignature").textContent = order.internal_signature;

    getId("receiverName").textContent = order.delivery.name;
    getId("phoneNumber").textContent = order.delivery.phone;
    getId("city").textContent = order.delivery.city;
    getId("address").textContent = order.delivery.address;
    getId("email").textContent = order.delivery.email;

    getId("transaction").textContent = order.payment.transaction;
    getId("currency").textContent = order.payment.currency;
    getId("amount").textContent = order.payment.amount;
    getId("deliveryCost").textContent = order.payment.delivery_cost;
    getId("itemsCost").textContent = order.payment.goods_total;
    getId("fee").textContent = order.payment.custom_fee;
    getId("bank").textContent = order.payment.bank;

    const itemsAmount = getId("itemsAmount");
    const itemsBody = getId("itemsBody");

    itemsAmount.textContent = order.items.length;
    itemsBody.innerHTML = "";

    order.items.forEach((item) => {
        const row = document.createElement("tr");

        row.innerHTML = `
        <td>${item.name}</td>
        <td>${item.brand}</td>
        <td>${item.total_price}</td>
        <td>${item.size}</td>
        <td>${item.status}</td>
        `;

        itemsBody.appendChild(row);
    });
}
